package miner

import (
	"context"
	"math/rand"
	"time"

	"0chain.net/block"
	"0chain.net/chain"
	"0chain.net/common"
	"0chain.net/datastore"
	. "0chain.net/logging"
	"0chain.net/memorystore"
	"0chain.net/node"
	"0chain.net/round"
	"0chain.net/transaction"
	"go.uber.org/zap"
)

var BLOCK_TIME = 3 * chain.DELTA

func SetNetworkRelayTime(delta time.Duration) {
	chain.SetNetworkRelayTime(delta)
	BLOCK_TIME = 3 * delta
}

func (mc *Chain) startNewRound(ctx context.Context, mr *Round) {
	if mr.Number < mc.CurrentRound {
		Logger.Debug("start new round (current round higher)", zap.Int64("round", mr.Number), zap.Int64("current_round", mc.CurrentRound))
		return
	}
	if !mc.AddRound(mr) {
		Logger.Debug("start new round (round already exists)", zap.Int64("round", mr.Number))
		return
	}
	pr := mc.GetRound(mr.Number - 1)
	//TODO: If for some reason the server is lagging behind (like network outage) we need to fetch the previous round info
	// before proceeding
	if pr == nil {
		Logger.Debug("start new round (previous round not found)", zap.Int64("round", mr.Number))
		return
	}
	self := node.GetSelfNode(ctx)
	rank := mr.GetRank(self.SetIndex)
	Logger.Info("*** starting round ***", zap.Any("round", mr.Number), zap.Any("index", self.SetIndex), zap.Any("rank", rank))
	if !mc.CanGenerateRound(&mr.Round, self.Node) {
		return
	}
	//NOTE: If there are not enough txns, this will not advance further even though rest of the network is. That's why this is a goroutine
	go mc.GenerateRoundBlock(ctx, mr)
}

/*GetBlockToExtend - Get the block to extend from the given round */
func (mc *Chain) GetBlockToExtend(r *Round) *block.Block {
	for true { // Need to do this for timing issues where a start round might come before a notarization and there is no notarized block to extend from
		rnb := r.GetNotarizedBlocks()
		if len(rnb) > 0 {
			if len(rnb) == 1 {
				return rnb[0]
			}
			//TODO: pick the best possible block
			return rnb[0]
		}
		if r.Number+1 != mc.CurrentRound {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	Logger.Debug("no block to extend", zap.Int64("round", r.Number), zap.Int64("current_round", mc.CurrentRound), zap.Int("nb_count", len(r.GetNotarizedBlocks())))
	return nil
}

/*GenerateRoundBlock - given a round number generates a block*/
func (mc *Chain) GenerateRoundBlock(ctx context.Context, r *Round) (*block.Block, error) {
	txnEntityMetadata := datastore.GetEntityMetadata("txn")
	ctx = memorystore.WithEntityConnection(ctx, txnEntityMetadata)
	defer memorystore.Close(ctx)
	pround := mc.GetRound(r.Number - 1)
	if pround == nil {
		Logger.Error("generate block (prior round not found)", zap.Any("round", r.Number-1))
		return nil, common.NewError("invalid_round,", "Round not available")
	}
	pb := mc.GetBlockToExtend(pround)
	if pb == nil {
		Logger.Error("generate block (prior block not found)", zap.Any("round", r.Number))
		return nil, common.NewError("block_gen_no_block_to_extend", "Do not have the block to extend this round")
	}
	b := datastore.GetEntityMetadata("block").Instance().(*block.Block)
	b.ChainID = mc.ID
	b.MagicBlockHash = mc.CurrentMagicBlock.Hash
	b.RoundRandomSeed = r.RandomSeed
	b.SetPreviousBlock(pb)
	for true {
		if mc.CurrentRound > b.Round {
			Logger.Error("generate block (round mismatch)", zap.Any("round", r.Number), zap.Any("current_round", mc.CurrentRound))
			return nil, common.NewError("round_mismatch", "Current round and block round do not match")
		}
		txnCount := transaction.TransactionCount
		err := mc.GenerateBlock(ctx, b, mc)
		if err != nil {
			Logger.Error("generate block", zap.Error(err))
			cerr, ok := err.(*common.Error)
			if ok && cerr.Code == InsufficientTxns {
				delay := 128 * time.Millisecond
				for true {
					if txnCount != transaction.TransactionCount {
						break
					}
					if mc.CurrentRound > b.Round {
						Logger.Error("generate block (round mismatch)", zap.Any("round", r.Number), zap.Any("current_round", mc.CurrentRound))
						return nil, common.NewError("round_mismatch", "Current round and block round do not match")
					}
					time.Sleep(delay)
					Logger.Debug("generate block", zap.Any("round", r.Number), zap.Any("delay", delay), zap.Any("txn_count", txnCount), zap.Any("t.txn_count", transaction.TransactionCount))
					delay = 2 * delay
					if delay > time.Second {
						delay = time.Second
					}
				}
				continue
			}
			Logger.Error("generate block", zap.Error(err))
			return nil, err
		}
		break
	}
	if mc.CurrentRound > b.Round {
		Logger.Error("generate block (round mismatch)", zap.Any("round", r.Number), zap.Any("current_round", mc.CurrentRound))
		return nil, common.NewError("round_mismatch", "Current round and block round do not match")
	}
	mc.AddBlock(b)
	mc.AddToRoundVerification(ctx, r, b)
	mc.SendBlock(ctx, b)
	return b, nil
}

/*AddToRoundVerification - Add a block to verify : WARNING: does not support concurrent access for a given round */
func (mc *Chain) AddToRoundVerification(ctx context.Context, mr *Round, b *block.Block) {
	if mr.IsFinalizing() || mr.IsFinalized() {
		Logger.Debug("add to verification", zap.Any("round", b.Round), zap.Any("block", b.Hash), zap.Any("finalizing", mr.IsFinalizing()), zap.Any("finalized", mr.IsFinalized()))
		return
	}
	if !mc.ValidateMagicBlock(ctx, b) {
		Logger.Error("invalid magic block", zap.Any("round", b.Round), zap.Any("block", b.Hash), zap.Any("magic_block", b.MagicBlockHash))
		return
	}
	mc.AddBlock(b)
	vctx := mr.StartVerificationBlockCollection(ctx)
	if vctx != nil {
		go mc.CollectBlocksForVerification(vctx, mr)
	}
	mr.AddBlockToVerify(b)
}

/*CollectBlocksForVerification - keep collecting the blocks till timeout and then start verifying */
func (mc *Chain) CollectBlocksForVerification(ctx context.Context, r *Round) {
	var blockTimeTimer = time.NewTimer(chain.DELTA)
	var sendVerification = false
	verifyAndSend := func(ctx context.Context, r *Round, b *block.Block) bool {
		bvt, err := mc.VerifyRoundBlock(ctx, r, b)
		if err != nil {
			if err == ErrRoundMismatch {
				Logger.Debug("verify round block", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Any("current_round", mc.CurrentRound))
			} else {
				Logger.Error("verify round block", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Error(err))
			}
			return false
		}
		r.Block = b

		//TODO: Dfinity suggests broadcasting the prior block so it saturates the network
		//While saturation is good, it's going to be expensive, hence TODO for now. Also, if we are proceeding verification based on partial block info,
		// we can't broadcast that block
		if !mc.IsBlockNotarized(ctx, b) {
			if b.MinerID != node.Self.GetKey() {
				mc.SendVerificationTicket(ctx, b, bvt)
			}
			mc.ProcessVerifiedTicket(ctx, r, b, &bvt.VerificationTicket)
		}
		return true
	}
	var blocks = make([]*block.Block, 0, 10)

	initiateVerification := func() {
		sendVerification = true
		// Sort the accumulated blocks by the rank and process them
		blocks = r.GetBlocksByRank(blocks)
		var verified bool
		// Keep verifying all the blocks collected so far in the best rank order till the first successul verification
		for _, b := range blocks {
			if verifyAndSend(ctx, r, b) {
				verified = true
				break
			}
		}
		if !verified {
			mc.startRound(&r.Round)
		}
	}
	for true {
		select {
		case <-ctx.Done():
			if !sendVerification {
				initiateVerification()
			}
			return
		case <-blockTimeTimer.C:
			initiateVerification()
		case b := <-r.GetBlocksToVerifyChannel():
			if sendVerification {
				// Is this better than the current best block
				if r.Block == nil || b.RoundRank < r.Block.RoundRank {
					verifyAndSend(ctx, r, b)
				}
			} else { // Accumulate all the blocks into this array till the BlockTime timeout
				blocks = append(blocks, b)
			}
		}
	}
}

/*VerifyRoundBlock - given a block is verified for a round*/
func (mc *Chain) VerifyRoundBlock(ctx context.Context, r *Round, b *block.Block) (*block.BlockVerificationTicket, error) {
	if mc.CurrentRound != r.Number {
		return nil, ErrRoundMismatch
	}
	if b.MinerID == node.Self.GetKey() {
		return mc.SignBlock(ctx, b)
	}
	if b.PrevBlock == nil {
		pb, err := mc.GetBlock(ctx, b.PrevHash)
		if err != nil {
			Logger.Error("verify round", zap.Any("round", r.Number), zap.Any("block", b.Hash), zap.Any("prev_block", b.PrevHash), zap.Error(err))
			//NOTE: when we don't have the prior block, we construct partial block to try to proceed
			//TODO: Need to figure out how to convert this partial block into a full block
			// key missing info for pb: txns and prevblock data.
			pb = datastore.GetEntityMetadata("block").Instance().(*block.Block)
			pb.ChainID = mc.ID
			pb.Round = b.Round - 1
			pb.MagicBlockHash = mc.CurrentMagicBlock.Hash
			pb.RoundRandomSeed = r.RandomSeed
			pb.Hash = b.PrevHash
			pb.VerificationTickets = b.PrevBlockVerficationTickets
			mc.AddBlock(pb)
		}
		b.PrevBlock = pb
	}
	/* Note: We are verifying the notarization of the previous block we have with
	   the prev verification tickets of the current block. This is right as all the
	   necessary verification tickets & notarization message may not have arrived to us */
	if err := mc.VerifyNotarization(ctx, b.PrevBlock, b.PrevBlockVerficationTickets); err != nil {
		return nil, err
	}

	bvt, err := mc.VerifyBlock(ctx, b)
	if err != nil {
		return nil, err
	}
	return bvt, nil
}

/*ProcessVerifiedTicket - once a verified ticket is receiveid, do further processing with it */
func (mc *Chain) ProcessVerifiedTicket(ctx context.Context, r *Round, b *block.Block, vt *block.VerificationTicket) {
	notarized := mc.IsBlockNotarized(ctx, b)
	//NOTE: We keep collecting verification tickets even if a block is notarized.
	// This is useful since Dfinity suggest broadcasting the previous notarized block when verifying the current block
	// If we know how many verifications already exists for a block, we only need to broadcast to the rest. Hence collecting any prior block verifications is OK.
	if !mc.AddVerificationTicket(ctx, b, vt) {
		return
	}
	if notarized {
		return
	}
	if mc.IsBlockNotarized(ctx, b) {
		r.Block = b
		mc.CancelRoundVerification(ctx, r)
		notarization := datastore.GetEntityMetadata("block_notarization").Instance().(*Notarization)
		notarization.BlockID = b.Hash
		notarization.Round = b.Round
		notarization.VerificationTickets = b.VerificationTickets
		mc.SendNotarization(ctx, notarization)
		mc.AddNotarizedBlock(ctx, &r.Round, b)
	}
}

/*AddNotarizedBlock - add a notarized block for a given round */
func (mc *Chain) AddNotarizedBlock(ctx context.Context, r *round.Round, b *block.Block) {
	r.AddNotarizedBlock(b)
	mc.startRound(r)

	pr := mc.GetRound(r.Number - 1)
	if pr != nil {
		pr.CancelVerification()
		mc.FinalizeRound(ctx, &pr.Round, mc)
	}
}

func (mc *Chain) startRound(r *round.Round) {
	if mc.GetRound(r.Number+1) == nil {
		nr := datastore.GetEntityMetadata("round").Instance().(*round.Round)
		nr.Number = r.Number + 1
		//TODO: We need to do VRF
		nr.RandomSeed = rand.New(rand.NewSource(r.RandomSeed)).Int63()
		nmr := mc.CreateRound(nr)
		// Even if the context is cancelled, we want to proceed with the next round, hence start with a root context
		Logger.Debug("starting a new round", zap.Int64("round", nr.Number))
		go mc.startNewRound(common.GetRootContext(), nmr)
		mc.Miners.SendAll(RoundStartSender(nr))
	}
}

/*CancelRoundVerification - cancel verifications happening within a round */
func (mc *Chain) CancelRoundVerification(ctx context.Context, r *Round) {
	r.CancelVerification() // No need for further verification of any blocks
}
