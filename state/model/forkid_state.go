package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
)

var (
	// ErrNewForkIdIsNotNext is returned when the new fork id is not the next one
	ErrNewForkIdIsNotNext = errors.New("ForkID must be the greatest that last one stored")
)

type ForkIDInterval = entities.ForkIDInterval

const FORKID_ZERO = entities.FORKID_ZERO

type StorageForkIdInterface interface {
	AddForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx storageTxType) error
	GetForkIDs(ctx context.Context, dbTx storageTxType) ([]pgstorage.ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx storageTxType) error
}

type ForkIdState struct {
	storage StorageForkIdInterface
	// cacheForkId[forkid] = ForkIDInterval
	cacheForkId *utils.Cache[uint64, ForkIDInterval]
}

// NewForkIdState creates a new ForkIdState that manage forksIds
func NewForkIdState(storage StorageForkIdInterface) *ForkIdState {
	return &ForkIdState{
		storage:     storage,
		cacheForkId: utils.NewCache[uint64, ForkIDInterval](utils.DefaultTimeProvider{}, utils.InfiniteTimeOfLiveItems),
	}
}

// AddForkIDInterval updates the forkID intervals
func (s *ForkIdState) AddForkID(ctx context.Context, newForkID ForkIDInterval, dbTx stateTxType) error {
	currentForksIDs, err := s.GetForkIDs(ctx, dbTx)
	if err != nil {
		return err
	}
	// Check if the forkId is already stored
	if s.isForkIDAlreadyStored(newForkID.ForkId, currentForksIDs) {
		// It's already stored, so we don't need to store it again
		return nil
	}
	// Check that is nextForkId
	err = s.checkValidNewForkID(newForkID, currentForksIDs)
	if err != nil {
		return err
	}
	previousForkID := s.adjustPreviousForkIDToBatchNumber(newForkID, currentForksIDs)
	if previousForkID != nil {
		err = s.storage.UpdateForkID(ctx, pgstorage.ForkIDInterval(*previousForkID), dbTx)
		if err != nil {
			return err
		}
	}
	err = s.storage.AddForkID(ctx, pgstorage.ForkIDInterval(newForkID), dbTx)
	if err != nil {
		return err
	}
	return err
}
func (p *ForkIdState) GetForkIDs(ctx context.Context, dbTx stateTxType) ([]ForkIDInterval, error) {
	currentForksId, err := p.storage.GetForkIDs(ctx, dbTx)
	if errors.Is(err, entities.ErrNotFound) {
		// It's ok, it's the first forkId
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	res := make([]ForkIDInterval, len(currentForksId))
	for i, v := range currentForksId {
		res[i] = ForkIDInterval(v)
	}
	return res, nil
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func (s *ForkIdState) GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) uint64 {
	forks, err := s.GetForkIDs(ctx, dbTx)
	if err != nil {
		log.Warnf("error getting forkIDs. Error: %v", err)
		return FORKID_ZERO
	}
	maxForId := uint64(0)
	for _, v := range forks {
		if batchNumber >= v.FromBatchNumber && batchNumber <= v.ToBatchNumber {
			return v.ForkId
		}
		if v.ForkId > maxForId {
			maxForId = v.ForkId
		}
	}
	log.Warnf("error can't match batch: %d in current forkids. Returning last one: %d", batchNumber, maxForId)
	return maxForId
}

// GetForkIDByBlockNumber returns the fork id for a given block number
func (s *ForkIdState) GetForkIDByBlockNumber(ctx context.Context, blockNumber uint64, dbTx stateTxType) uint64 {
	forks, err := s.GetForkIDs(ctx, dbTx)
	if err != nil {
		log.Warnf("error getting forkIDs. Error: %v", err)
		return FORKID_ZERO
	}
	maxForId := uint64(0)
	for _, v := range forks {

		if v.BlockNumber > blockNumber {
			return maxForId
		} else {
			maxForId = v.ForkId
		}
	}
	return maxForId
}

func (s *ForkIdState) checkValidNewForkID(newForkID ForkIDInterval, currentForksIDs []ForkIDInterval) error {
	// Check that is nextForkId
	if len(currentForksIDs) == 0 {
		return nil
	}
	lastForkId := currentForksIDs[len(currentForksIDs)-1]
	if newForkID.ForkId < lastForkId.ForkId+1 {
		return fmt.Errorf("last ForkID stored: %d. New ForkID received: %d. err:%w", lastForkId.ForkId, newForkID.ForkId, ErrNewForkIdIsNotNext)
	}
	return nil
}

func (s *ForkIdState) isForkIDAlreadyStored(forkID uint64, currentForksIDs []ForkIDInterval) bool {
	for _, v := range currentForksIDs {
		if v.ForkId == forkID {
			return true
		}
	}
	return false
}

func (s *ForkIdState) adjustPreviousForkIDToBatchNumber(newForkID ForkIDInterval, currentForksIDs []ForkIDInterval) *ForkIDInterval {
	// If the new to adjust previous ForkID to the new FromBatchNumber
	if len(currentForksIDs) == 0 {
		return nil
	}
	previousForkID := currentForksIDs[len(currentForksIDs)-1]
	previousForkID.ToBatchNumber = newForkID.FromBatchNumber - 1
	return &previousForkID
}
