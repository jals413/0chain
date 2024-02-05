package provider

import (
	"fmt"

	"0chain.net/smartcontract/dbs"
	"0chain.net/smartcontract/dbs/event"

	"0chain.net/smartcontract/stakepool"

	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/smartcontractinterface"
)

var AlreadyShutdownError = fmt.Errorf("already killed or shutdown")

func ShutDown(
	input []byte,
	clientId, ownerId string,
	killSlash float64,
	providerSpecific func(ProviderRequest) (AbstractProvider, stakepool.AbstractStakePool, error),
	refreshProvider func(ProviderRequest) error,
	balances cstate.StateContextI,
) error {
	var req ProviderRequest
	if err := req.Decode(input); err != nil {
		return err
	}

	p, sp, err := providerSpecific(req)
	if err != nil {
		return err
	}

	actErr := cstate.WithActivation(balances, "apollo", func() (e error) {
		if p.IsShutDown() {
			e = fmt.Errorf("already shutdown")
		}
		if p.IsKilled() {
			e = fmt.Errorf("already killed")
		}
		return e
	}, func() (e error) {
		if p.IsKilled() || p.IsShutDown() {
			if refreshProvider != nil {
				e = refreshProvider(req)
				if e != nil {
					return e
				}
			}

			e = AlreadyShutdownError
		}
		return e
	})

	if actErr != nil {
		return actErr
	}

	p.ShutDown()

	actErr = cstate.WithActivation(balances, "apollo",
		func() (e error) { return nil },
		func() (e error) {
			if err = sp.Kill(killSlash, p.Id(), p.Type(), balances); err != nil {
				e = fmt.Errorf("can't kill the stake pool: %v", err)
			}
			return e
		})

	if actErr != nil {
		return actErr
	}

	if err = sp.Save(p.Type(), clientId, balances); err != nil {
		return err
	}

	var errCode = "shutdown_" + p.Type().String() + "_failed"
	if err := smartcontractinterface.AuthorizeWithOwner(errCode, func() bool {
		return ownerId == clientId || clientId == sp.GetSettings().DelegateWallet
	}); err != nil {
		return err
	}

	balances.EmitEvent(event.TypeStats, event.TagShutdownProvider, p.Id(), dbs.ProviderID{
		ID:   p.Id(),
		Type: p.Type(),
	})

	return nil
}
