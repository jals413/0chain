package minersc

import (
	"fmt"
	"reflect"
	"runtime"

	"0chain.net/chaincore/block"
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/node"
	"0chain.net/chaincore/threshold/bls"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/smartcontract/stakepool/spenum"
	"github.com/0chain/common/core/logging"
	"github.com/0chain/common/core/util"

	"go.uber.org/zap"
)

var moveFunctions = make(map[Phase]movePhaseFunctions)

/*
- Start      : moveToContribute
- Contribute : moveToShareOrPublish
- Share      : moveToShareOrPublish
- Publish    : moveToWait
- Wait       : moveToStart
*/

func (msc *MinerSmartContract) moveToContribute(balances cstate.StateContextI,
	pn *PhaseNode, gn *GlobalNode) error {

	var (
		allMinersList *MinerNodes
		dkgMinersList *DKGMinerNodes

		allShardersList *MinerNodes

		err error
	)

	if allMinersList, err = msc.getMinersList(balances); err != nil {
		return common.NewError("move_to_contribute_failed", err.Error())
	}

	if dkgMinersList, err = getDKGMinersList(balances); err != nil {
		return common.NewError("move_to_contribute_failed", err.Error())
	}
	allShardersList, err = getAllShardersList(balances)
	if err != nil {
		return common.NewError("move_to_contribute_failed", err.Error())
	}

	if len(allShardersList.Nodes) < gn.MinS {
		return common.NewErrorf("move_to_contribute_failed",
			"not enough sharders in all sharders list to move phase, all: %d, min_s: %d",
			len(allShardersList.Nodes), gn.MinS)
	}

	if !gn.hasPrevShader(allShardersList, balances) {
		return common.NewErrorf("move_to_contribute_failed",
			"invalid state: all sharders list hasn't a sharder from previous VC set, "+
				"all: %d, min_s: %d", len(allShardersList.Nodes), gn.MinS)
	}

	if !gn.hasPrevMiner(allMinersList, balances) {
		return common.NewErrorf("move_to_contribute_failed",
			"invalid state: all miners list hasn't a miner from previous VC set, "+
				"all: %d, min_n: %d", len(allMinersList.Nodes), gn.MinN)
	}

	if allMinersList == nil {
		return common.NewError("move_to_contribute_failed", "allMinersList is nil")
	}

	if len(allMinersList.Nodes) < dkgMinersList.K {
		return common.NewErrorf("move_to_contribute_failed",
			"len(allMinersList.Nodes) < dkgMinersList.K, l_miners: %d, K: %d",
			len(allMinersList.Nodes), dkgMinersList.K)
	}

	if len(allShardersList.Nodes) < gn.MinS {
		return common.NewErrorf("move_to_contribute_failed",
			"len(allShardersList.Nodes) < gn.MinS, l_shards: %d, min_s: %d",
			len(allShardersList.Nodes), gn.MinS)
	}

	logging.Logger.Debug("miner sc: move phase to contribute",
		zap.Int("miners", len(allMinersList.Nodes)),
		zap.Int("K", dkgMinersList.K),
		zap.Int("sharders", len(allShardersList.Nodes)),
		zap.Int("min_s", gn.MinS))
	return nil
}

// moveToShareOrPublish move function for contribute phase
func (msc *MinerSmartContract) moveToShareOrPublish(
	balances cstate.StateContextI, pn *PhaseNode, gn *GlobalNode) error {

	// TODO: do sharders in auto VC
	// shardersKeep, err := getShardersKeepList(balances)
	// if err != nil {
	// 	return common.NewErrorf("move_to_share_or_publish_failed",
	// 		"failed to get sharders keep list: %v", err)
	// }

	// if len(shardersKeep.Nodes) < gn.MinS {
	// 	return common.NewErrorf("move_to_share_or_publish_failed",
	// 		"not enough sharders in keep list to move phase keep: %d, min_s: %d", len(shardersKeep.Nodes), gn.MinS)
	// }

	// if !gn.hasPrevShader(shardersKeep, balances) {
	// 	return common.NewErrorf("move_to_share_or_publish_failed",
	// 		"missing at least one sharder from previous set in "+
	// 			"sharders keep list to move phase, keep: %d, min_s: %d",
	// 		len(shardersKeep.Nodes), gn.MinS)
	// }

	dkgMinersList, err := getDKGMinersList(balances)
	if err != nil {
		return common.NewErrorf("move_to_share_or_publish_failed",
			"failed to get miners DKG, phase: %v, err: %v", pn.Phase, err)
	}

	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()

	mpks, err := getMinersMPKs(balances)
	if err != nil {
		return common.NewErrorf("move_to_share_or_publish_failed",
			"phase: %v, err: %v", pn.Phase, err)
	}

	if len(mpks.Mpks) == 0 {
		return common.NewError("move_to_share_or_publish_failed", "empty miners mpks keys")
	}

	// should have at least one miner from previous VC set
	if !gn.hasPrevMinerInMPKs(mpks, balances) {
		return common.NewErrorf("move_to_share_or_publish_failed",
			"no miner from previous VC set in MPKS, l_mpks: %d, DB: %d, DB version: %d",
			len(mpks.Mpks),
			int(balances.GetState().GetVersion()),
			int(balances.GetState().GetVersion()))
	}

	if len(mpks.Mpks) < dkgMinersList.K {
		return common.NewErrorf("move_to_share_or_publish_failed",
			"len(mpks.Mpks) < dkgMinersList.K, l_mpks: %d, K: %d",
			len(mpks.Mpks), dkgMinersList.K)
	}

	logging.Logger.Debug("[mvc] miner sc: move phase to share or publish",
		zap.Int("mpks", len(mpks.Mpks)),
		zap.Int("K", dkgMinersList.K),
		zap.Int64("DB version", int64(balances.GetState().GetVersion())))

	return nil
}

func (msc *MinerSmartContract) moveToWait(balances cstate.StateContextI,
	pn *PhaseNode, gn *GlobalNode) error {

	dkgMinersList, err := getDKGMinersList(balances)
	if err != nil {
		return common.NewErrorf("move_to_wait_failed",
			"failed to get miners DKG, phase: %v, err: %v", pn.Phase, err)
	}

	gsos, err := getGroupShareOrSigns(balances)
	switch err {
	case nil:
	case util.ErrValueNotPresent:
		return common.NewError("move_to_wait_failed", "empty sharder or sign keys")
	default:
		return common.NewErrorf("move_to_wait_failed", "phase: %v, err: %v", pn.Phase, err)
	}

	if !gn.hasPrevMinerInGSoS(gsos, balances) {
		return common.NewErrorf("move_to_wait_failed",
			"no miner from previous VC set in GSoS, "+
				"l_gsos: %d, K: %d, DB version: %d, gsos_shares: %v",
			len(gsos.Shares), dkgMinersList.K, int(balances.GetState().GetVersion()), gsos.Shares)
	}

	if len(gsos.Shares) < dkgMinersList.K {
		return common.NewErrorf("move_to_wait_failed",
			"len(gsos.Shares) < dkgMinersList.K, l_gsos: %d, K: %d",
			len(gsos.Shares), dkgMinersList.K)
	}

	logging.Logger.Debug("miner sc: move phase to wait",
		zap.Int("shares", len(gsos.Shares)),
		zap.Int("K", dkgMinersList.K))

	return nil
}

func (msc *MinerSmartContract) moveToStart(balances cstate.StateContextI,
	pn *PhaseNode, gn *GlobalNode) error {
	return nil
}

// GetPhaseNode gets phase nodes from client state
func GetPhaseNode(statectx cstate.CommonStateContextI) (*PhaseNode, error) {
	pn := &PhaseNode{}
	err := statectx.GetTrieNode(pn.GetKey(), pn)
	if err != nil {
		if err != util.ErrValueNotPresent {
			return nil, err
		}

		pn.Phase = Start
		pn.CurrentRound = statectx.GetBlock().Round
		pn.StartRound = statectx.GetBlock().Round
		return pn, nil
	}

	pn.CurrentRound = statectx.GetBlock().Round
	return pn, nil
}

func (msc *MinerSmartContract) setPhaseNode(balances cstate.StateContextI,
	pn *PhaseNode,
	gn *GlobalNode,
	t *transaction.Transaction) (err error) {
	defer func() {
		if err == nil {
			// TODO: optimize this to check if there are any changes, avoid unnecessary MPT saving
			_, err = balances.InsertTrieNode(pn.GetKey(), pn)
			return
		}
	}()

	// move phase condition
	var movePhase = pn.CurrentRound-pn.StartRound >= PhaseRounds[pn.Phase]
	if !movePhase {
		return nil
	}

	logging.Logger.Debug("[mvc] setPhaseNode",
		zap.String("phase", pn.Phase.String()),
		zap.Int64("phase current round", pn.CurrentRound),
		zap.Int64("phase start round", pn.StartRound))

	currentMoveFunc := moveFunctions[pn.Phase]
	err = currentMoveFunc(balances, pn, gn)
	if err != nil {
		if cstate.ErrInvalidState(err) {
			logging.Logger.Error("[mvc] setPhaseNode failed",
				zap.Error(err),
				zap.String("phase", pn.Phase.String()),
				zap.Int64("phase current round", pn.CurrentRound),
				zap.Int64("phase start round", pn.StartRound))
			return err
		}

		funcName := getFunctionName(currentMoveFunc)
		logging.Logger.Error("[mvc] failed to move phase, restartDKG",
			zap.String("phase", pn.Phase.String()),
			zap.Int64("phase start round", pn.StartRound),
			zap.Int64("phase current round", pn.CurrentRound),
			zap.String("move_func", funcName),
			zap.Error(err))
		err = msc.RestartDKG(pn, balances)
		if err != nil {
			logging.Logger.Error("[mvc] setPhaseNode restart DKG failed",
				zap.Error(err),
				zap.Int64("phase start round", pn.StartRound),
				zap.Int64("phase current round", pn.CurrentRound),
				zap.String("move_func", funcName))
		}
		return err
	}

	logging.Logger.Debug("[mvc] setPhaseNode move phase success", zap.String("phase", pn.Phase.String()))
	phaseFunc, ok := phaseFuncs[pn.Phase]
	if ok {
		err = phaseFunc(balances, gn)
		if err != nil {
			if cstate.ErrInvalidState(err) {
				logging.Logger.Debug("[mvc] setPhaseNode move phase failed",
					zap.Error(err),
					zap.String("phase", pn.Phase.String()))
				return err
			}

			logging.Logger.Error("[mvc] setPhaseNode failed to set phase node - restarting DKG",
				zap.Error(err),
				zap.String("phase", pn.Phase.String()))

			if err = msc.RestartDKG(pn, balances); err != nil {
				logging.Logger.Debug("setPhaseNode move phase failed",
					zap.Error(err),
					zap.String("phase", pn.Phase.String()))
				return err
			}
		}
	}

	if pn.Phase >= Phase(len(PhaseRounds)-1) {
		pn.Phase = 0
		pn.Restarts = 0
	} else {
		pn.Phase++
	}
	pn.StartRound = pn.CurrentRound
	logging.Logger.Debug("[mvc] setPhaseNode", zap.String("next_phase", pn.Phase.String()))
	return nil
}

func (msc *MinerSmartContract) createDKGMinersForContribute(
	balances cstate.StateContextI, gn *GlobalNode) error {
	logging.Logger.Debug("[mvc] createDKGMinersForContribute start")
	allMinersList, err := msc.getMinersList(balances)
	if err != nil {
		logging.Logger.Error("createDKGMinersForContribute -- failed to get miner list",
			zap.Error(err))
		return err
	}

	// get deleted miners list
	dKey, ok := deleteNodeKeyMap[spenum.Miner]
	if !ok {
		return fmt.Errorf("get delete node key failed, invalid node type: %s", spenum.Miner)
	}

	deleteMinersIDs, err := getDeleteNodeIDs(balances, dKey)
	if err != nil && err != util.ErrValueNotPresent {
		logging.Logger.Error("[mvc] createDKGMinersForContribute failed to get delete miners IDs", zap.Error(err))
		return err
	}

	logging.Logger.Debug("[mvc] createDKGMinersForContribute remove deleted miners",
		zap.Strings("miners", deleteMinersIDs),
		zap.Int("all miners num", len(allMinersList.Nodes)))

	allMinersList.RemoveNodes(deleteMinersIDs)

	if len(allMinersList.Nodes) < gn.MinN {
		logging.Logger.Error("[mvc] createDKGMinersForContribute too few miners for dkg",
			zap.Int("all miners num", len(allMinersList.Nodes)),
			zap.Int("min miners num", gn.MinN))
		return common.NewErrorf("failed to create dkg miners", "too few miners for dkg, l_all_miners: %d, N: %d", len(allMinersList.Nodes), gn.MinN)
	}

	var (
		allMinersNum = len(allMinersList.Nodes)
		lmb          = balances.GetChainCurrentMagicBlock()
		mbMinersNum  int
	)

	if lmb != nil {
		mbMinersNum = lmb.Miners.Size()
	}

	dkgMiners := NewDKGMinerNodes()
	if mbMinersNum < allMinersNum && mbMinersNum >= gn.MinN {
		logging.Logger.Debug("[mvc] createDKGMinersForContribute calculate TKN from lmb",
			zap.Int64("starting round", lmb.StartingRound),
			zap.Int("miners num", mbMinersNum))
		dkgMiners.calculateTKN(gn, mbMinersNum)
	} else {
		logging.Logger.Debug("[mvc] createDKGMinersForContribute calculate TKN from all miner list",
			zap.Int("all count", allMinersNum),
			zap.Int64("gn.LastRound", gn.LastRound))
		dkgMiners.calculateTKN(gn, allMinersNum)
	}

	for _, nd := range allMinersList.Nodes {
		dkgMiners.SimpleNodes[nd.ID] = nd.SimpleNode
	}

	dkgMiners.StartRound = gn.LastRound
	logging.Logger.Debug("[mvc] createDKGMinersForContribute, new dkg miners",
		zap.Int("T", dkgMiners.T),
		zap.Int("K", dkgMiners.K),
		zap.Int("N", dkgMiners.N),
		zap.Int("all count", len(dkgMiners.SimpleNodes)),
		zap.Int64("start round", dkgMiners.StartRound))
	if err := updateDKGMinersList(balances, dkgMiners); err != nil {
		logging.Logger.Error("[mvc] createDKGMinersForContribute, failed to update dkg miners list")
		return err
	}

	// TODO: do the below in auto VC
	// sharders
	// allSharderKeepList := new(MinerNodes)
	// return updateShardersKeepList(balances, allSharderKeepList)
	return nil
}

func (msc *MinerSmartContract) widdleDKGMinersForShare(
	balances cstate.StateContextI, gn *GlobalNode) error {

	dkgMiners, err := getDKGMinersList(balances)
	if err != nil {
		logging.Logger.Error("widdle dkg miners -- failed to get dkgMiners",
			zap.Error(err))
		return err
	}

	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()

	mpks, err := getMinersMPKs(balances)
	if err != nil {
		logging.Logger.Error("widdle dkg miners -- failed to get miners mpks",
			zap.Error(err))
		return err
	}

	for k := range dkgMiners.SimpleNodes {
		if _, ok := mpks.Mpks[k]; !ok {
			delete(dkgMiners.SimpleNodes, k)
		}
	}

	if err = dkgMiners.reduceNodes(false, gn, balances); err != nil {
		logging.Logger.Error("widdle dkg miners", zap.Error(err))
		return err
	}

	if err := updateDKGMinersList(balances, dkgMiners); err != nil {
		logging.Logger.Error("widdle dkg miners -- failed to insert dkg miners",
			zap.Error(err))
		return err
	}
	return nil
}

func (msc *MinerSmartContract) reduceShardersList(
	keep,
	all *MinerNodes,
	gn *GlobalNode,
	balances cstate.StateContextI) (nodes []*MinerNode, err error) {

	simpleNodes := NewSimpleNodes()

	tmpMinerNodes := make([]*MinerNode, 0, len(keep.Nodes))

	for _, keepNode := range keep.Nodes {
		var found = all.FindNodeById(keepNode.ID)
		if found == nil {
			return nil, common.NewErrorf("invalid state", "a sharder exists in"+
				" keep list doesn't exists in all sharders list: %s", keepNode.ID)
		}
		tmpMinerNodes = append(tmpMinerNodes, found)
		simpleNodes[found.ID] = found.SimpleNode
	}

	if len(simpleNodes) < gn.MinS {
		return nil, fmt.Errorf("to few sharders: %d, want at least: %d", len(simpleNodes), gn.MinS)
	}

	var pmbrss int64
	var pmbnp *node.Pool
	pmb := balances.GetLastestFinalizedMagicBlock()
	if pmb != nil {
		pmbrss = pmb.RoundRandomSeed
		if pmb.MagicBlock != nil {
			pmbnp = pmb.MagicBlock.Sharders
		}
	}
	logging.Logger.Debug("sharder keep before", zap.Int("num", len(simpleNodes)))
	simpleNodes.reduce(gn.MaxS, gn.XPercent, pmbrss, pmbnp)
	logging.Logger.Debug("sharder keep after", zap.Int("num", len(simpleNodes)))

	nodes = make([]*MinerNode, 0, len(simpleNodes))

	for _, mn := range tmpMinerNodes {
		if sn, ok := simpleNodes[mn.ID]; ok {
			mn.SimpleNode = sn
			nodes = append(nodes, mn)
		}
	}

	if !hasPrevSharderInList(pmb.MagicBlock, nodes) {
		var prev = rankedPrevSharders(pmb.MagicBlock, nodes)
		if len(prev) == 0 {
			panic("must not happen")
		}
		nodes = append(nodes, prev[0])
	}

	return
}

func (msc *MinerSmartContract) createMagicBlockForWait(
	balances cstate.StateContextI, gn *GlobalNode) error {

	pn, err := GetPhaseNode(balances)
	if err != nil {
		return err
	}
	dkgMinersList, err := getDKGMinersList(balances)
	if err != nil {
		return err
	}

	gsos, err := getGroupShareOrSigns(balances)
	switch err {
	case nil:
	case util.ErrValueNotPresent:
		gsos = block.NewGroupSharesOrSigns()
	default:
		return err
	}

	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()

	mpks, err := getMinersMPKs(balances)
	if err != nil {
		return common.NewError("create_magic_block_failed", err.Error())
	}

	for key := range mpks.Mpks {
		if _, ok := gsos.Shares[key]; !ok {
			logging.Logger.Debug("create magic block - delete miner because no share found", zap.String("key", key))
			delete(dkgMinersList.SimpleNodes, key)
			delete(gsos.Shares, key)
			delete(mpks.Mpks, key)
		}
	}
	// TODO: check the necessary of the commented code below
	//for key, sharesRevealed := range dkgMinersList.RevealedShares {
	//	if sharesRevealed >= dkgMinersList.T {
	//		Logger.Debug("create magic block - delete miner because share revealed >= T found",
	//			zap.String("key", key),
	//			zap.Int("shares revealed", sharesRevealed),
	//			zap.Int("T", dkgMinersList.T))
	//		delete(dkgMinersList.SimpleNodes, key)
	//		delete(gsos.Shares, key)
	//		delete(mpks.Mpks, key)
	//	}
	//}

	// sharders
	sharders, err := getShardersKeepList(balances)
	if err != nil {
		return err
	}
	allSharderList, err := getAllShardersList(balances)
	if err != nil {
		return err
	}

	// TODO: reduce sharders in auto VC
	// if sharders == nil || len(sharders.Nodes) == 0 {
	// 	sharders = allSharderList
	// } else {
	// 	sharders.Nodes, err = msc.reduceShardersList(sharders, allSharderList, gn, balances)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	sharders = allSharderList

	if err = dkgMinersList.reduceNodes(true, gn, balances); err != nil {
		logging.Logger.Error("create magic block for wait - reduce nodes failed", zap.Error(err))
		return err
	}

	for id := range gsos.Shares {
		if _, ok := dkgMinersList.SimpleNodes[id]; !ok {
			delete(gsos.Shares, id)
		}
	}

	for id := range mpks.Mpks {
		if _, ok := dkgMinersList.SimpleNodes[id]; !ok {
			delete(mpks.Mpks, id)
		}
	}

	if len(dkgMinersList.SimpleNodes) < dkgMinersList.K {
		return common.NewErrorf("create_magic_block_failed",
			"len(dkgMinersList.SimpleNodes) [%d] < dkgMinersList.K [%d]", len(dkgMinersList.SimpleNodes), dkgMinersList.K)
	}

	magicBlock, err := msc.createMagicBlock(balances, sharders, dkgMinersList, gsos, mpks, pn)
	if err != nil {
		return err
	}

	gn.ViewChange = magicBlock.StartingRound
	mpks = block.NewMpks()
	if err := updateMinersMPKs(balances, mpks); err != nil {
		return err
	}

	logging.Logger.Debug("create_mpks in createMagicBlockForWait", zap.Int64("DB version", int64(balances.GetState().GetVersion())))

	gsos = block.NewGroupSharesOrSigns()
	if err := updateGroupShareOrSigns(balances, gsos); err != nil {
		return err
	}

	err = updateMagicBlock(balances, magicBlock)
	if err != nil {
		logging.Logger.Error("failed to insert magic block", zap.Error(err))
		return err
	}
	// dkgMinersList = NewDKGMinerNodes()
	// _, err = balances.InsertTrieNode(DKGMinersKey, dkgMinersList)
	// if err != nil {
	// 	return err
	// }
	// allSharderKeepList := new(MinerNodes)
	return updateShardersKeepList(balances, NodeIDs{})
}

func (msc *MinerSmartContract) contributeMpk(t *transaction.Transaction,
	inputData []byte, gn *GlobalNode, balances cstate.StateContextI) (
	string, error) {
	logging.Logger.Debug("[mvc] miner smart contract, contribute Mpk")
	pn, err := GetPhaseNode(balances)
	if err != nil {
		logging.Logger.Error("[mvc] contribute mpk failed to get phase node", zap.Error(err))
		return "", common.NewErrorf("contribute_mpk_failed", "can't get phase node: %v", err)
	}

	if pn.Phase != Contribute {
		logging.Logger.Error("[mvc] contribute mpk, not the correct phase to contribute mpk",
			zap.String("phase", pn.Phase.String()))
		return "", common.NewErrorf("contribute_mpk_failed",
			"this is not the correct phase to contribute mpk: %v", pn.Phase.String())
	}

	dmn, err := getDKGMinersList(balances)
	if err != nil {
		logging.Logger.Error("[mvc] contribute mpk failed to get dkg miners", zap.Error(err))
		return "", err
	}

	if _, ok := dmn.SimpleNodes[t.ClientID]; !ok {
		logging.Logger.Error("[mvc] contribute mpk miner not part of dkg set", zap.String("miner", t.ClientID))
		return "", common.NewError("contribute_mpk_failed",
			"miner not part of dkg set")
	}

	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()

	mpk := &block.MPK{ID: t.ClientID}
	if err := mpk.Decode(inputData); err != nil {
		return "", common.NewErrorf("contribute_mpk_failed",
			"decoding request: %v", err)
	}

	if len(mpk.Mpk) != dmn.T {
		logging.Logger.Error("[mvc] contribute mpk, mpk sent (size: %v) is not correct size: %v",
			zap.Int("mpk size", len(mpk.Mpk)), zap.Int("dkg size", dmn.T))
		return "", common.NewErrorf("contribute_mpk_failed",
			"mpk sent (size: %v) is not correct size: %v", len(mpk.Mpk), dmn.T)
	}

	mpks, err := getMinersMPKs(balances)
	switch err {
	case util.ErrValueNotPresent:
		// the mpks could be empty when the first time to contribute mpks
		logging.Logger.Debug("[mvc] contribute create new mpks")
		mpks = block.NewMpks()
	case nil:
	default:
		return "", common.NewError("contribute_mpk_failed", err.Error())
	}

	if _, ok := mpks.Mpks[mpk.ID]; ok {
		return "", common.NewError("contribute_mpk_failed", "already have mpk for miner")
	}

	mpks.Mpks[mpk.ID] = mpk
	if err := updateMinersMPKs(balances, mpks); err != nil {
		return "", common.NewError("contribute_mpk_failed", err.Error())
	}

	logging.Logger.Debug("[mvc] contribute_mpk success",
		zap.Int64("DB version", int64(balances.GetState().GetVersion())),
		zap.String("mpk id", mpk.ID),
		zap.Int("len", len(mpks.Mpks)),
		zap.Int64("pn_start_round", pn.StartRound),
		zap.String("phase", pn.Phase.String()))

	return string(mpk.Encode()), nil
}

func (msc *MinerSmartContract) shareSignsOrShares(t *transaction.Transaction,
	inputData []byte, gn *GlobalNode, balances cstate.StateContextI) (
	resp string, err error) {

	var pn *PhaseNode
	if pn, err = GetPhaseNode(balances); err != nil {
		logging.Logger.Error("[mvc] shareSignsOrShares, can't get phase node",
			zap.Error(err))
		return "", common.NewErrorf("share_signs_or_shares",
			"can't get phase node: %v", err)
	}

	if pn.Phase != Publish {
		logging.Logger.Error("[mvc] shareSignsOrShares, not in publish phase",
			zap.String("phase", pn.Phase.String()))
		return "", common.NewErrorf("share_signs_or_shares",
			"this is not the correct phase to publish signs or shares, phase node: %v",
			string(pn.Encode()))
	}

	var gsos *block.GroupSharesOrSigns
	gsos, err = getGroupShareOrSigns(balances)
	switch err {
	case nil:
	case util.ErrValueNotPresent:
		gsos = block.NewGroupSharesOrSigns()
	default:
		logging.Logger.Error("[mvc] failed to get group share or signs",
			zap.Error(err))
		return "", common.NewError("share_signs_or_shares_failed", err.Error())
	}

	var ok bool
	if _, ok = gsos.Shares[t.ClientID]; ok {
		logging.Logger.Error("[mvc] shareSignsOrShares, already have share or signs for miner",
			zap.String("miner", t.ClientID))
		return "", common.NewErrorf("share_signs_or_shares",
			"already have share or signs for miner %v", t.ClientID)
	}

	var dmn *DKGMinerNodes
	if dmn, err = getDKGMinersList(balances); err != nil {
		logging.Logger.Error("[mvc] shareSignsOrShares, failed to get miners DKG list",
			zap.Error(err))
		return "", common.NewErrorf("share_signs_or_shares",
			"getting miners DKG list %v", err)
	}

	var sos = block.NewShareOrSigns()
	if err = sos.Decode(inputData); err != nil {
		logging.Logger.Error("[mvc] shareSignsOrShares, failed to decode sc input", zap.Error(err))
		return "", common.NewErrorf("share_signs_or_shares", "decoding input %v", err)
	}

	if len(sos.ShareOrSigns) < dmn.K-1 {
		logging.Logger.Debug("[mvc] shareSignsOrShares, not enough share or signs for this dkg",
			zap.Int("l_sos", len(sos.ShareOrSigns)),
			zap.Int("K", dmn.K-1))
		return "", common.NewErrorf("share_signs_or_shares",
			"not enough share or signs for this dkg, l_sos: %d, K - 1: %d",
			len(sos.ShareOrSigns), dmn.K-1)
	}

	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()

	var mpks *block.Mpks
	mpks, err = getMinersMPKs(balances)
	if err != nil {
		logging.Logger.Error("[mvc] shareSignsOrShares, getting miners MPKs",
			zap.Error(err))
		return "", common.NewError("share_signs_or_shares_failed", err.Error())
	}

	var publicKeys = make(map[string]string)
	for key, miner := range dmn.SimpleNodes {
		publicKeys[key] = miner.PublicKey
	}

	var shares []string
	shares, ok = sos.Validate(mpks, publicKeys, balances.GetSignatureScheme())
	if !ok {
		logging.Logger.Error("[mvc] shareSignsOrShares, validation failed")
		return "", common.NewError("share_signs_or_shares", "share or signs failed validation")
	}

	for _, share := range shares {
		dmn.RevealedShares[share]++
	}

	sos.ID = t.ClientID
	gsos.Shares[t.ClientID] = sos

	err = updateGroupShareOrSigns(balances, gsos)
	if err != nil {
		logging.Logger.Error("[mvc] miner sc: shareSignsOrShares, failed to save group share of signs",
			zap.Error(err))
		return "", common.NewErrorf("share_signs_or_shares",
			"saving group share of signs: %v", err)
	}

	if err := updateDKGMinersList(balances, dmn); err != nil {
		logging.Logger.Error("[mvc] miner sc: shareSignsOrShares, failed to save DKG miners",
			zap.Error(err))
		return "", common.NewErrorf("share_signs_or_shares",
			"saving DKG miners: %v", err)
	}

	logging.Logger.Debug("[mvc] miner sc: shareSignsOrShares, update gsos",
		zap.String("miner", t.ClientID),
		zap.Int("gsos shares len", len(gsos.Shares)),
		zap.Int64("gn.LastRound", gn.LastRound))

	return string(sos.Encode()), nil
}

// Wait used to notify SC that DKG summary and magic block data saved by
// a miner locally and VC can moves on. If >= T miners don't save their
// DKG/MB then we can't do a view change because some miners can save,
// but some can not. And this case we got BC stuck.
func (msc *MinerSmartContract) wait(t *transaction.Transaction,
	inputData []byte, gn *GlobalNode, balances cstate.StateContextI) (
	resp string, err error) {

	var pn *PhaseNode
	if pn, err = GetPhaseNode(balances); err != nil {
		return "", common.NewErrorf("msc - wait",
			"can't get phase node: %v", err)
	}

	if pn.Phase != Wait {
		return "", common.NewErrorf("msc - wait", "this is not the"+
			" correct phase to wait: %s", pn.Phase)
	}

	var dmn *DKGMinerNodes
	if dmn, err = getDKGMinersList(balances); err != nil {
		return "", common.NewErrorf("msc - wait", "can't get DKG miners: %v", err)
	}

	if already, ok := dmn.Waited[t.ClientID]; ok && already {
		return "", common.NewError("msc - wait", "already checked in")
	}

	dmn.Waited[t.ClientID] = true
	logging.Logger.Debug("[mvc] wait", zap.Int("dkg miners num", len(dmn.SimpleNodes)))

	if err := updateDKGMinersList(balances, dmn); err != nil {
		return "", common.NewErrorf("msc - wait", "saving DKG miners: %v", err)
	}

	return
}

func (msc *MinerSmartContract) createMagicBlock(
	balances cstate.StateContextI,
	sharders *MinerNodes,
	dkgMinersList *DKGMinerNodes,
	gsos *block.GroupSharesOrSigns,
	mpks *block.Mpks,
	pn *PhaseNode,
) (*block.MagicBlock, error) {

	pmb := balances.GetLastestFinalizedMagicBlock()

	magicBlock := block.NewMagicBlock()

	magicBlock.Miners = node.NewPool(node.NodeTypeMiner)
	magicBlock.Sharders = node.NewPool(node.NodeTypeSharder)
	magicBlock.SetShareOrSigns(gsos)
	magicBlock.Mpks = mpks
	magicBlock.T = dkgMinersList.T
	magicBlock.K = dkgMinersList.K
	magicBlock.N = dkgMinersList.N
	magicBlock.MagicBlockNumber = pmb.MagicBlock.MagicBlockNumber + 1
	magicBlock.PreviousMagicBlockHash = pmb.MagicBlock.Hash
	magicBlock.StartingRound = pn.CurrentRound + PhaseRounds[Wait]

	logging.Logger.Debug("create magic block",
		zap.Int64("view change", magicBlock.StartingRound),
		zap.Int("dkg miners num", len(dkgMinersList.SimpleNodes)),
		zap.String("mb miners pool type", magicBlock.Miners.Type.String()),
		zap.String("mb sharders pool type", magicBlock.Sharders.Type.String()))

	for _, v := range dkgMinersList.SimpleNodes {
		n := node.Provider()
		n.ID = v.ID
		n.N2NHost = v.N2NHost
		n.Host = v.Host
		n.Port = v.Port
		n.Path = v.Path
		n.PublicKey = v.PublicKey
		n.Description = v.ShortName
		n.Type = node.NodeTypeMiner
		n.Info.BuildTag = v.BuildTag
		n.Status = node.NodeStatusActive
		n.InPrevMB = pmb.Miners.HasNode(v.ID)
		if err := magicBlock.Miners.AddNode(n); err != nil {
			return nil, err
		}

		mn := NewMinerNode()
		mn.SimpleNode = v
		mn.Status = n.Status
		// TODO: should not emit add miner for existing miner
		// emitAddMiner(mn, balances)
	}

	for _, v := range sharders.Nodes {
		n := node.Provider()
		n.ID = v.ID
		n.N2NHost = v.N2NHost
		n.Host = v.Host
		n.Port = v.Port
		n.Path = v.Path
		n.PublicKey = v.PublicKey
		n.Description = v.ShortName
		n.Type = node.NodeTypeSharder
		n.Info.BuildTag = v.BuildTag
		n.Status = node.NodeStatusActive
		n.InPrevMB = pmb.Sharders.HasNode(v.ID)
		if err := magicBlock.Sharders.AddNode(n); err != nil {
			return nil, err
		}

		sn := v
		sn.Status = n.Status
		// emitAddSharder(sn, balances)
	}

	magicBlock.Hash = magicBlock.GetHash()
	return magicBlock, nil
}

func (msc *MinerSmartContract) RestartDKG(pn *PhaseNode,
	balances cstate.StateContextI) error {
	logging.Logger.Debug("RestartDKG", zap.Int64("DB version", int64(balances.GetState().GetVersion())))
	msc.mutexMinerMPK.Lock()
	defer msc.mutexMinerMPK.Unlock()
	mpks := block.NewMpks()
	if err := updateMinersMPKs(balances, mpks); err != nil {
		logging.Logger.Error("failed to restart dkg", zap.Error(err))
		return err
	}

	logging.Logger.Debug("create_mpks in RestartDKG", zap.Int64("DB version", int64(balances.GetState().GetVersion())))

	gsos := block.NewGroupSharesOrSigns()
	if err := updateGroupShareOrSigns(balances, gsos); err != nil {
		logging.Logger.Error("failed to restart dkg", zap.Error(err))
		return err
	}
	dkgMinersList := NewDKGMinerNodes()
	dkgMinersList.StartRound = pn.CurrentRound
	if err := updateDKGMinersList(balances, dkgMinersList); err != nil {
		logging.Logger.Error("failed to restart dkg", zap.Error(err))
		return err
	}

	// TODO: do the below in auto VC
	// if err := updateShardersKeepList(balances, NodeIDs{}); err != nil {
	// 	logging.Logger.Error("failed to restart dkg", zap.Error(err))
	// 	return err
	// }
	pn.Phase = Start
	pn.Restarts++
	pn.StartRound = pn.CurrentRound
	return nil
}

func (msc *MinerSmartContract) SetMagicBlock(gn *GlobalNode,
	balances cstate.StateContextI) error {
	magicBlock, err := getMagicBlock(balances)
	if err != nil {
		logging.Logger.Error("could not get magic block from MPT", zap.Error(err))
		return err
	}

	if magicBlock.StartingRound == 0 && magicBlock.MagicBlockNumber == 0 {
		logging.Logger.Error("SetMagicBlock smart contract, starting round is 0, magic block number is 0")
	}

	// keep the magic block to track previous nodes list next view change
	// (deny VC leaving for at least 1 miner and 1 sharder of previous set)
	gn.PrevMagicBlock = magicBlock

	logging.Logger.Debug("SetMagicBlock",
		zap.String("hash", magicBlock.Hash),
		zap.Int64("starting round", magicBlock.StartingRound),
		zap.String("miners pool type", magicBlock.Miners.Type.String()),
		zap.String("sharders pool type", magicBlock.Sharders.Type.String()),
		zap.Int("miners num", magicBlock.Miners.Size()),
		zap.Int("sharders num", magicBlock.Sharders.Size()))
	balances.SetMagicBlock(magicBlock)
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func newDKGWithMagicBlock(mb *block.MagicBlock, summary *bls.DKGSummary) (*bls.DKG, error) {
	selfNodeKey := node.Self.Underlying().GetKey()

	if summary.SecretShares == nil {
		return nil, common.NewError("failed to set dkg from store", "no saved shares for dkg")
	}

	var newDKG = bls.MakeDKG(mb.T, mb.N, selfNodeKey)
	newDKG.MagicBlockNumber = mb.MagicBlockNumber
	newDKG.StartingRound = mb.StartingRound

	if mb.Miners == nil {
		return nil, common.NewError("failed to set dkg from store", "miners pool is not initialized in magic block")
	}

	for k := range mb.Miners.CopyNodesMap() {
		if savedShare, ok := summary.SecretShares[computeBlsID(k)]; ok {
			if err := newDKG.AddSecretShare(bls.ComputeIDdkg(k), savedShare, false); err != nil {
				return nil, err
			}
		} else if v, ok := mb.GetShareOrSigns().Get(k); ok {
			if share, ok := v.ShareOrSigns[node.Self.Underlying().GetKey()]; ok && share.Share != "" {
				if err := newDKG.AddSecretShare(bls.ComputeIDdkg(k), share.Share, false); err != nil {
					return nil, err
				}
			}
		}
	}

	if !newDKG.HasAllSecretShares() {
		return nil, common.NewError("failed to set dkg from store",
			"not enough secret shares for dkg")
	}

	newDKG.AggregateSecretKeyShares()
	newDKG.Pi = newDKG.Si.GetPublicKey()
	mpks, err := mb.Mpks.GetMpkMap()
	if err != nil {
		return nil, err
	}

	if err := newDKG.AggregatePublicKeyShares(mpks); err != nil {
		return nil, err
	}

	return newDKG, nil
}
