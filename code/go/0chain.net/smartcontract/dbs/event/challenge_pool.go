package event

import (
	"gorm.io/gorm/clause"
)

type ChallengePool struct {
	ID           string `gorm:"primarykey"`
	AllocationID string `gorm:"uniqueIndex"`
	Balance      int64  `json:"balance"`
	StartTime    int64  `json:"start_time"`
	Expiration   int64  `json:"expiration"`
	Finalized    bool   `json:"finalized"`
}

func (edb *EventDb) addOrUpdateChallengePools(cps []ChallengePool) error {
	updateFields := []string{"balance", "start_time", "expiration", "finalized"}

	return edb.Store.Get().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns(updateFields), // column needed to be updated
	}).Create(&cps).Error
}

func (edb *EventDb) GetChallengePool(allocationID string) (*ChallengePool, error) {
	var cp ChallengePool
	return &cp, edb.Store.Get().Model(&ChallengePool{}).
		Where("allocation_id = ?", allocationID).
		Take(&cp).Error
}

func mergeAddChallengePoolsEvents() *eventsMergerImpl[ChallengePool] {
	return newEventsMerger[ChallengePool](TagAddOrUpdateChallengePool, withUniqueEventOverwrite())
}

func mergeToChallengePoolsEvents() *eventsMergerImpl[ChallengePool] {
	return newEventsMerger[ChallengePool](TagToChallengePool, withMergeToChallengePool())
}

func mergeFromChallengePoolsEvents() *eventsMergerImpl[ChallengePool] {
	return newEventsMerger[ChallengePool](TagFromChallengePool, withMergeFromChallengePool())
}

func withMergeToChallengePool() eventMergeMiddleware {
	return withEventMerge(func(a, b *ChallengePool) (*ChallengePool, error) {
		a.Balance += b.Balance
		a.StartTime = b.StartTime
		a.Expiration = b.Expiration
		a.Finalized = b.Finalized
		return a, nil
	})
}

func withMergeFromChallengePool() eventMergeMiddleware {
	return withEventMerge(func(a, b *ChallengePool) (*ChallengePool, error) {
		a.Balance -= b.Balance
		a.StartTime = b.StartTime
		a.Expiration = b.Expiration
		a.Finalized = b.Finalized
		return a, nil
	})
}
