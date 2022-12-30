package event

import (
	"errors"
	"fmt"

	"github.com/0chain/common/core/currency"
	"gorm.io/gorm/clause"
)

type Authorizer struct {
	Provider

	URL string `json:"url"`

	// Configuration
	Fee currency.Coin `json:"fee"`

	// Geolocation
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

}

func (edb *EventDb) AddAuthorizer(a *Authorizer) error {
	exists, err := a.exists(edb)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("authorizer already exists")
	}

	result := edb.Store.Get().Create(a)

	return result.Error
}

func (edb *EventDb) GetAuthorizer(id string) (*Authorizer, error) {
	var auth Authorizer

	result := edb.Store.Get().
		Model(&Authorizer{}).
		Where(&Authorizer{Provider: Provider{ID: id}}).
		First(&auth)

	if result.Error != nil {
		return nil, fmt.Errorf(
			"error retrieving authorizer %v, error %v",
			id, result.Error,
		)
	}

	return &auth, nil
}

func (edb *EventDb) GetAuthorizers() ([]Authorizer, error) {
	var authorizers []Authorizer
	result := edb.Store.Get().
		Model(&Authorizer{}).
		Find(&authorizers)
	return authorizers, result.Error
}

func (edb *EventDb) DeleteAuthorizer(id string) error {
	result := edb.Store.Get().
		Where("authorizer_id = ?", id).
		Delete(&Authorizer{})
	return result.Error
}

func (a *Authorizer) exists(edb *EventDb) (bool, error) {
	var count int64

	result := edb.Get().
		Model(&Authorizer{}).
		Where(&Authorizer{Provider: Provider{ID: a.ID}}).
		Count(&count)

	if result.Error != nil {
		return false,
			fmt.Errorf(
				"error searching for authorizer %v, error %v",
				a.ID, result.Error,
			)
	}
	return count > 0, nil
}

func NewUpdateAuthorizerTotalStakeEvent(ID string, totalStake currency.Coin) (tag EventTag, data interface{}) {
	return TagUpdateAuthorizerTotalStake, Authorizer{
		Provider: Provider{
			ID:         ID,
			TotalStake: totalStake,
		},
	}
}

func NewUpdateAuthorizerTotalUnStakeEvent(ID string, totalUnstake currency.Coin) (tag EventTag, data interface{}) {
	return TagUpdateAuthorizerTotalStake, Authorizer{
		Provider: Provider{
			ID:         ID,
			TotalStake: totalUnstake,
		},
	}
}

func (edb *EventDb) updateAuthorizersTotalStakes(authorizer []Authorizer) error {
	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"total_stake"}),
	}).Create(&authorizer).Error
}

func mergeUpdateAuthorizerTotalStakesEvents() *eventsMergerImpl[Authorizer] {
	return newEventsMerger[Authorizer](TagUpdateAuthorizerTotalStake, withUniqueEventOverwrite())
}
func (edb *EventDb) updateAuthorizersTotalUnStakes(authorizer []Authorizer) error {
	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"unstake_total"}),
	}).Create(&authorizer).Error
}

func mergeUpdateAuthorizerTotalUnStakesEvents() *eventsMergerImpl[Authorizer] {
	return newEventsMerger[Authorizer](TagUpdateAuthorizerTotalUnStake, withUniqueEventOverwrite())
}
