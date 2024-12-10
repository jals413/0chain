package round

import (
	"errors"
	"github.com/0chain/common/core/logging"
	"go.uber.org/zap"
	"sync"
)

var (
	ErrRoundEntityNotFound = errors.New("round entity not found")
)

type RoundStorageEntity = interface{}

type RoundStorage interface {
	Get(round int64) RoundStorageEntity
	GetLatest() RoundStorageEntity
	Put(entity RoundStorageEntity, round int64) error
	Prune(round int64) error
	Count() int
	GetRound(i int) int64
	FindRoundIndex(round int64) int
	GetRounds() []int64
}

type roundStartingStorage struct {
	mu     *sync.RWMutex
	max    int64
	items  map[int64]RoundStorageEntity
	rounds []int64
}

var _ RoundStorage = (*roundStartingStorage)(nil)

func NewRoundStartingStorage() *roundStartingStorage {
	store := &roundStartingStorage{
		items:  make(map[int64]RoundStorageEntity),
		rounds: make([]int64, 0),
		mu:     &sync.RWMutex{},
	}
	return store
}

func (s *roundStartingStorage) Get(round int64) RoundStorageEntity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	found := s.calcNearestRound(round)
	if found == -1 {
		logging.Logger.Info("roundStartingStorage", zap.Any("round", round), zap.Any("entity", nil))
		return nil
	}
	entity, ok := s.items[found]
	if !ok {
		logging.Logger.Info("roundStartingStorage", zap.Any("round", round), zap.Any("entity", nil))
		return nil
	}

	logging.Logger.Info("roundStartingStorage", zap.Any("round", round), zap.Any("entity", entity))
	return entity
}

func (s *roundStartingStorage) FindRoundIndex(round int64) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if round > s.max && s.max > 0 {
		return len(s.rounds) - 1
	}
	found := -1
	for i := 0; i < len(s.rounds); i++ {
		if round >= s.rounds[i] {
			found = i
		} else {
			break
		}
	}
	return found
}

func (s *roundStartingStorage) calcNearestRound(round int64) int64 {
	if round > s.max && s.max > 0 {
		return s.max
	}
	found := int64(-1)
	for i := 0; i < len(s.rounds); i++ {
		if round >= s.rounds[i] {
			found = s.rounds[i]
		} else {
			break
		}
	}
	return found
}

func (s *roundStartingStorage) GetLatest() RoundStorageEntity {

	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.items) == 0 {
		logging.Logger.Info("roundStartingStorage", zap.Any("round", s.max), zap.Any("entity", nil))
		return nil
	}

	logging.Logger.Info("roundStartingStorage", zap.Any("round", s.max), zap.Any("entity", s.items[s.max]))
	return s.items[s.max]
}

func (s *roundStartingStorage) Put(entity RoundStorageEntity, round int64) error {
	logging.Logger.Info("roundStartingStorage", zap.Any("round", round), zap.Any("entity", entity))
	s.mu.Lock()
	defer s.mu.Unlock()
	if round > s.max {
		s.max = round
	}

	_, found := s.items[round]
	s.items[round] = entity
	if !found {
		s.putToSlice(round)
	}

	return nil
}

func (s *roundStartingStorage) putToSlice(round int64) {
	logging.Logger.Info("roundStartingStorage", zap.Any("round", round))
	index := -1
	for i := len(s.rounds) - 1; i >= 0; i-- {
		sRound := s.rounds[i]
		if sRound < round {
			index = i
			break
		}
	}
	if index == -1 {
		s.rounds = append([]int64{round}, s.rounds...)
	} else {
		s.rounds = append(s.rounds[:index+1], append([]int64{round}, s.rounds[index+1:]...)...)
	}
}

func (s *roundStartingStorage) check(round int64) error { //nolint
	return nil
}

func (s *roundStartingStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	logging.Logger.Info("roundStartingStorage", zap.Any("count", len(s.items)))
	return len(s.items)
}

func (s *roundStartingStorage) GetRounds() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]int64, len(s.rounds))
	copy(result, s.rounds)

	logging.Logger.Info("roundStartingStorage", zap.Any("rounds", result))
	return result
}

func (s *roundStartingStorage) GetRound(i int) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.rounds[i]
}

func (s *roundStartingStorage) Prune(round int64) error {
	logging.Logger.Info("roundStartingStorage", zap.Any("round", round))
	s.mu.Lock()
	defer s.mu.Unlock()
	_, found := s.items[round]
	if !found {
		return ErrRoundEntityNotFound
	}
	pruneIndex := -1
	pruneRounds := make([]int64, 0)
	for i := 0; i < len(s.rounds); i++ {
		pruneRounds = append(pruneRounds, s.rounds[i])
		if round == s.rounds[i] {
			pruneIndex = i
			break
		}
	}
	if pruneIndex == -1 {
		return ErrRoundEntityNotFound
	}

	for _, roundRemove := range pruneRounds {
		delete(s.items, roundRemove)
	}
	s.rounds = s.rounds[pruneIndex+1:]
	return nil
}
