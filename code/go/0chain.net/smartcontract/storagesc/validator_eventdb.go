package storagesc

import (
	"0chain.net/smartcontract/dbs"
	"0chain.net/smartcontract/dto"
	"0chain.net/smartcontract/provider"
	"0chain.net/smartcontract/stakepool"
	"0chain.net/smartcontract/stakepool/spenum"
	"github.com/0chain/common/core/logging"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/smartcontract/dbs/event"
)

func validatorTableToValidationNode(v event.Validator) *ValidationNode {
	return &ValidationNode{
		Provider: provider.Provider{
			ID:           v.ID,
			ProviderType: spenum.Validator,
		},
		BaseURL:   v.BaseUrl,
		PublicKey: v.PublicKey,
		StakePoolSettings: stakepool.Settings{
			DelegateWallet:     v.DelegateWallet,
			MaxNumDelegates:    v.NumDelegates,
			ServiceChargeRatio: v.ServiceCharge,
		},
	}
}

func getValidators(validatorIDs []string, edb *event.EventDb) ([]*ValidationNode, error) {
	validators, err := edb.GetValidatorsByIDs(validatorIDs)
	if err != nil {
		return nil, err
	}
	vNodes := make([]*ValidationNode, len(validators))
	for i := range validators {
		vNodes[i] = validatorTableToValidationNode(validators[i])
	}

	return vNodes, nil
}

func emitUpdateValidationNode(validationNode *dto.ValidationDtoNode,
	existingStakePool *stakePool, balances cstate.StateContextI) error {
	staked, err := existingStakePool.stake()
	if err != nil {
		return err
	}

	logging.Logger.Info("emitting validator update event")

	data := &event.Validator{
		Provider: event.Provider{
			ID:         validationNode.ID,
			TotalStake: staked,
		},
	}

	if validationNode.BaseURL != nil {
		data.BaseUrl = *validationNode.BaseURL
	}

	if validationNode.StakePoolSettings != nil {
		if validationNode.StakePoolSettings.DelegateWallet != nil {
			data.DelegateWallet = *validationNode.StakePoolSettings.DelegateWallet
		}
		if validationNode.StakePoolSettings.ServiceChargeRatio != nil {
			data.ServiceCharge = *validationNode.StakePoolSettings.ServiceChargeRatio
		}
		if validationNode.StakePoolSettings.MaxNumDelegates != nil {
			data.NumDelegates = *validationNode.StakePoolSettings.MaxNumDelegates
		}
	}

	if validationNode.LastHealthCheck != nil {
		data.LastHealthCheck = *validationNode.LastHealthCheck
	}

	balances.EmitEvent(event.TypeStats, event.TagUpdateValidator, validationNode.ID, data)
	return nil
}

func (vn *ValidationNode) emitAddOrOverwrite(sp *stakePool, balances cstate.StateContextI) error {
	staked, err := sp.stake()
	if err != nil {
		return err
	}

	logging.Logger.Info("emitting validator add or overwrite event")
	data := &event.Validator{
		BaseUrl: vn.BaseURL,
		Provider: event.Provider{
			ID:              vn.ID,
			TotalStake:      staked,
			DelegateWallet:  vn.StakePoolSettings.DelegateWallet,
			NumDelegates:    vn.StakePoolSettings.MaxNumDelegates,
			ServiceCharge:   vn.StakePoolSettings.ServiceChargeRatio,
			Rewards:         event.ProviderRewards{ProviderID: vn.ID},
			LastHealthCheck: vn.LastHealthCheck,
		},
	}

	balances.EmitEvent(event.TypeStats, event.TagAddOrOverwiteValidator, vn.ID, data)
	return nil
}

func emitValidatorHealthCheck(vn *ValidationNode, downtime uint64, balances cstate.StateContextI) {
	data := dbs.DbHealthCheck{
		ID:              vn.ID,
		LastHealthCheck: vn.LastHealthCheck,
		Downtime:        downtime,
	}

	balances.EmitEvent(event.TypeStats, event.TagValidatorHealthCheck, vn.ID, data)
}
