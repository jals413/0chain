package storagesc

import (
	"encoding/json"
	"fmt"

	"0chain.net/smartcontract/stakepool"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/smartcontract/dbs/event"
)

func writeMarkerToValidationNode(vn *ValidationNode) *event.Validator {
	return &event.Validator{
		ValidatorID: vn.ID,
		BaseUrl:     vn.BaseURL,
		PublicKey:   vn.PublicKey,
		// TO-DO: Update stake in eventDB
		Stake: 0,

		DelegateWallet: vn.StakePoolSettings.DelegateWallet,
		MinStake:       vn.StakePoolSettings.MinStake,
		MaxStake:       vn.StakePoolSettings.MaxStake,
		NumDelegates:   vn.StakePoolSettings.MaxNumDelegates,
		ServiceCharge:  vn.StakePoolSettings.ServiceCharge,
	}
}

func validatorTableToValidationNode(v event.Validator) *ValidationNode {
	return &ValidationNode{
		ID:        v.ValidatorID,
		BaseURL:   v.BaseUrl,
		PublicKey: v.PublicKey,
		StakePoolSettings: stakepool.StakePoolSettings{
			DelegateWallet:  v.DelegateWallet,
			MinStake:        v.MinStake,
			MaxStake:        v.MaxStake,
			MaxNumDelegates: v.NumDelegates,
			ServiceCharge:   v.ServiceCharge,
		},
	}
}

func emitAddOrOverwriteValidatorTable(vn *ValidationNode, balances cstate.StateContextI, t *transaction.Transaction) error {

	data, err := json.Marshal(writeMarkerToValidationNode(vn))
	if err != nil {
		return fmt.Errorf("failed to marshal writemarker: %v", err)
	}

	balances.EmitEvent(event.TypeStats, event.TagAddOrOverwriteValidator, t.Hash, string(data))

	return nil
}

func getValidators(validatorIDs []string, edb *event.EventDb) ([]*ValidationNode, error) {
	validators, err := edb.GetValidatorsByID(validatorIDs)
	if err != nil {
		return nil, err
	}
	vNodes := make([]*ValidationNode, len(validators))
	for i := range validators {
		vNodes[i] = validatorTableToValidationNode(validators[i])
	}

	return vNodes, nil
}
