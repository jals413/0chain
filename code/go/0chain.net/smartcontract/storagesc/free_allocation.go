package storagesc

import (
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/0chain/common/core/currency"
	"github.com/0chain/common/core/logging"
	"github.com/0chain/common/core/util"
	"go.uber.org/zap"
)

const (
	floatToBalance = 10 * 1000 * 1000 * 1000
)

//msgp:ignore freeStorageAllocationInput newFreeStorageAssignerInfo
//go:generate msgp -io=false -tests=false -unexported=true -v

type freeStorageMarker struct {
	Assigner   string   `json:"assigner"`
	Recipient  string   `json:"recipient"`
	FreeTokens float64  `json:"free_tokens"`
	Nonce      int64    `json:"nonce"`
	Signature  string   `json:"signature"`
	Blobbers   []string `json:"blobbers"`
}

func (frm *freeStorageMarker) decode(b []byte) error {
	return json.Unmarshal(b, frm)
}

// swagger:model freeStorageAllocationInput
type freeStorageAllocationInput struct {
	RecipientPublicKey string   `json:"recipient_public_key"`
	Marker             string   `json:"marker"`
	Blobbers           []string `json:"blobbers"`
}

func (frm *freeStorageAllocationInput) decode(b []byte) error {
	return json.Unmarshal(b, frm)
}

type freeStorageUpgradeInput struct {
	AllocationId string `json:"allocation_id"`
	Marker       string `json:"marker"`
}

type newFreeStorageAssignerInfo struct {
	Name            string  `json:"name"`
	PublicKey       string  `json:"public_key"`
	IndividualLimit float64 `json:"individual_limit"`
	TotalLimit      float64 `json:"total_limit"`
}

func (frm *newFreeStorageAssignerInfo) decode(b []byte) error {
	return json.Unmarshal(b, frm)
}

func freeStorageAssignerKey(sscKey, clientId string) datastore.Key {
	return sscKey + ":freestorageredeemed:" + clientId
}

type freeStorageAssigner struct {
	ClientId        string        `json:"client_id"`
	PublicKey       string        `json:"public_key"`
	IndividualLimit currency.Coin `json:"individual_limit"`
	TotalLimit      currency.Coin `json:"total_limit"`
	CurrentRedeemed currency.Coin `json:"current_redeemed"`
	RedeemedNonces  []int64       `json:"redeemed_nonces"`
}

func (fsa *freeStorageAssigner) Encode() []byte {
	var b, err = json.Marshal(fsa)
	if err != nil {
		panic(err)
	}
	return b
}

func (fsa *freeStorageAssigner) Decode(p []byte) error {
	return json.Unmarshal(p, fsa)
}

func (fsa *freeStorageAssigner) save(sscKey string, balances cstate.StateContextI) error {
	_, err := balances.InsertTrieNode(freeStorageAssignerKey(sscKey, fsa.ClientId), fsa)
	return err
}

// TODO test that we really send some value here
func (fsa *freeStorageAssigner) validate(
	marker freeStorageMarker,
	now common.Timestamp,
	value currency.Coin,
	balances cstate.StateContextI,
) error {
	verified, e := verifyFreeAllocationRequestNew(marker, fsa.PublicKey, balances)
	if e != nil {
		return e
	}
	if !verified {
		return fmt.Errorf("failed to verify signature")
	}

	newTotal, err := currency.AddCoin(fsa.CurrentRedeemed, value)
	if err != nil {
		return err
	}

	if newTotal > fsa.TotalLimit {
		return fmt.Errorf("%d exceeded total permitted free storage limit %d", newTotal, fsa.TotalLimit)
	}

	if value > fsa.IndividualLimit {
		return fmt.Errorf("%d exceeded permitted free storage  %d", value, fsa.IndividualLimit)
	}

	for _, nonce := range fsa.RedeemedNonces {
		if marker.Nonce == nonce {
			return fmt.Errorf("marker already redeemed, nonce: %v", marker.Nonce)
		}
	}

	return nil
}

func (ssc *StorageSmartContract) addFreeStorageAssigner(
	t *transaction.Transaction,
	input []byte,
	balances cstate.StateContextI,
) (string, error) {
	var conf *Config
	var err error
	if conf, err = ssc.getConfig(balances, true); err != nil {
		return "", common.NewErrorf("add_free_storage_assigner",
			"can't get config: %v", err)
	}

	if err := smartcontractinterface.AuthorizeWithOwner("add_free_storage_assigner", func() bool {
		return conf.OwnerId == t.ClientID
	}); err != nil {
		return "", err
	}

	var assignerInfo newFreeStorageAssignerInfo
	if err := assignerInfo.decode(input); err != nil {
		return "", common.NewErrorf("add_free_storage_assigner",
			"can't unmarshal input: %v", err)
	}

	newTotalLimit, err := currency.Float64ToCoin(assignerInfo.TotalLimit * floatToBalance)
	if err != nil {
		return "", common.NewErrorf("add_free_storage_assigner", "can't convert total limit to coin: %v", err)
	}

	if newTotalLimit > conf.MaxTotalFreeAllocation {
		return "", common.NewErrorf("add_free_storage_assigner",
			"total tokens limit %d exceeds maximum permitted: %d", newTotalLimit, conf.MaxTotalFreeAllocation)
	}

	newIndividualLimit, err := currency.Float64ToCoin(assignerInfo.IndividualLimit * floatToBalance)
	if err != nil {
		return "", common.NewErrorf("add_free_storage_assigner", "can't convert individual limit to coin: %v", err)
	}

	if newIndividualLimit > conf.MaxIndividualFreeAllocation {
		return "", common.NewErrorf("add_free_storage_assigner",
			"individual allocation token limit %d exceeds maximum permitted: %d", newIndividualLimit, conf.MaxIndividualFreeAllocation)
	}

	assigner, err := ssc.getFreeStorageAssigner(assignerInfo.Name, balances)
	if err != nil && err != util.ErrValueNotPresent {
		return "", common.NewError("add_free_storage_assigner", err.Error())
	}
	if err == util.ErrValueNotPresent || assigner == nil {
		assigner = &freeStorageAssigner{
			ClientId: assignerInfo.Name,
		}
	}
	assigner.PublicKey = assignerInfo.PublicKey
	assigner.TotalLimit = newTotalLimit
	assigner.IndividualLimit = newIndividualLimit
	err = assigner.save(ssc.ID, balances)
	if err != nil {
		return "", common.NewErrorf("add_free_storage_assigner", "error saving new assigner: %v", err)
	}

	return "", nil
}

func verifyFreeAllocationRequestNew(
	frm freeStorageMarker,
	publicKey string,
	balances cstate.StateContextI,
) (bool, error) {
	var ids string
	for _, b := range frm.Blobbers {
		ids += b
	}
	marker := fmt.Sprintf("%s:%f:%d:%s", frm.Recipient, frm.FreeTokens, frm.Nonce, ids)
	logging.Logger.Debug("free_storage_marker verify", zap.String("marker", marker), zap.String("pub_key", publicKey))
	signatureScheme := balances.GetSignatureScheme()
	if err := signatureScheme.SetPublicKey(publicKey); err != nil {
		return false, err
	}
	return signatureScheme.Verify(frm.Signature, hex.EncodeToString([]byte(marker)))
}

func (ssc *StorageSmartContract) freeAllocationRequest(
	txn *transaction.Transaction,
	input []byte,
	balances cstate.StateContextI,
) (string, error) {
	var err error
	var inputObj freeStorageAllocationInput
	if err := inputObj.decode(input); err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"unmarshal input: %v", err)
	}

	var marker freeStorageMarker
	if err := marker.decode([]byte(inputObj.Marker)); err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"unmarshal request: %v", err)
	}

	var conf *Config
	if conf, err = ssc.getConfig(balances, true); err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"can't get config: %v", err)
	}

	if txn.ClientID != marker.Recipient {
		return "", common.NewErrorf("free_allocation_failed",
			"marker can be used only by its recipient")
	}

	assigner, err := ssc.getFreeStorageAssigner(marker.Assigner, balances)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"error getting assigner details: %v", err)
	}

	coin, err := currency.ParseZCN(marker.FreeTokens)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"marker verification failed: %v", err)
	}
	//todo query sharder on 0box to get the price of allocation
	if err := assigner.validate(marker, txn.CreationDate, coin, balances); err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"marker verification failed: %v", err)
	}

	var blobberAuthTickets []string

	request := newAllocationRequest{
		DataShards:           conf.FreeAllocationSettings.DataShards,
		ParityShards:         conf.FreeAllocationSettings.ParityShards,
		Size:                 conf.FreeAllocationSettings.Size,
		Owner:                marker.Recipient,
		OwnerPublicKey:       inputObj.RecipientPublicKey,
		ReadPriceRange:       conf.FreeAllocationSettings.ReadPriceRange,
		WritePriceRange:      conf.FreeAllocationSettings.WritePriceRange,
		Blobbers:             marker.Blobbers,
		ThirdPartyExtendable: true,
		StorageVersion:       1, // Use storageV2 for free allocation
	}

	for range marker.Blobbers {
		blobberAuthTickets = append(blobberAuthTickets, "")
	}
	request.BlobberAuthTickets = blobberAuthTickets

	arBytes, err := request.encode()
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"marshal request: %v", err)
	}

	free, err := currency.ParseZCN(marker.FreeTokens)
	if err != nil {
		return "", err
	}

	newRedeemed, err := currency.AddCoin(assigner.CurrentRedeemed, free)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed", "add coins: %v", err)
	}

	assigner.CurrentRedeemed = newRedeemed

	if err != nil {
		return "", err
	}
	readPoolTokens, err := currency.Float64ToCoin(float64(free) * conf.FreeAllocationSettings.ReadPoolFraction)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed", "converting read pool tokens to Coin: %v", err)
	}
	writePoolTokens, err := currency.MinusCoin(free, readPoolTokens)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed",
			"subtracting read pool token from transaction value: %v", err)
	}

	resp, err := ssc.newAllocationRequestInternal(txn, arBytes, conf, NewTokenTransfer(writePoolTokens, conf.OwnerId, txn.ToClientID, true), balances, nil)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed", "creating new allocation: %v", err)
	}

	var sa StorageAllocation
	if err := sa.Decode([]byte(resp)); err != nil {
		return "", common.NewErrorf("free_allocation_failed", "unmarshalling allocation: %v", err)
	}

	assigner.RedeemedNonces = append(assigner.RedeemedNonces, marker.Nonce)
	if err := assigner.save(ssc.ID, balances); err != nil {
		return "", common.NewErrorf("free_allocation_failed", "assigner Save failed: %v", err)
	}

	txn.Value = readPoolTokens
	_, err = ssc.readPoolLockInternal(txn, NewTokenTransfer(readPoolTokens, conf.OwnerId, txn.ToClientID, true), marker.Recipient, balances)
	if err != nil {
		return "", common.NewErrorf("free_allocation_failed", "locking tokens in read pool: %v", err)
	}

	return resp, err
}

func (ssc *StorageSmartContract) getFreeStorageAssigner(
	clientID datastore.Key,
	balances cstate.StateContextI,
) (*freeStorageAssigner, error) {
	fsa := new(freeStorageAssigner)
	err := balances.GetTrieNode(freeStorageAssignerKey(ssc.ID, clientID), fsa)
	if err != nil {
		return nil, err
	}

	return fsa, nil
}
