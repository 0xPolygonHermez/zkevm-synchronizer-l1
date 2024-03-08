package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
	"github.com/jackc/pgx/v4"
)

const (
	// FORKID_ZERO is the fork id 0 (no forkid)
	FORKID_ZERO = uint64(0)
	// FORKID_BLUEBERRY is the fork id 4
	FORKID_BLUEBERRY = uint64(4)
	// FORKID_DRAGONFRUIT is the fork id 5
	FORKID_DRAGONFRUIT = uint64(5)
	// FORKID_INCABERRY is the fork id 6
	FORKID_INCABERRY = uint64(6)
	// FORKID_ETROG is the fork id 7
	FORKID_ETROG = uint64(7)
	// FORKID_ELDERBERRY is the fork id 8
	FORKID_ELDERBERRY = uint64(8)
)

var (
	// ErrNewForkIdIsNotNext is returned when the new fork id is not the next one
	ErrNewForkIdIsNotNext = errors.New("ForkID must be the greatest that last one stored")
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}

type StorageForkIdInterface interface {
	AddForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx pgx.Tx) error
	GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx pgx.Tx) error
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
func (s *ForkIdState) AddForkID(ctx context.Context, newForkID ForkIDInterval, dbTx pgx.Tx) error {
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
	// TODO: When add batch support check that
	//s.checkIfForkIDAffectsPreviousBatches(ctx, newForkID, dbTx)
	return err
}
func (p *ForkIdState) GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]ForkIDInterval, error) {
	currentForksId, err := p.storage.GetForkIDs(ctx, dbTx)
	if errors.Is(err, pgstorage.ErrNotFound) {
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
func (s *ForkIdState) GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) uint64 {
	return FORKID_ZERO
}

// GetForkIDByBlockNumber returns the fork id for a given block number
func (s *ForkIdState) GetForkIDByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) uint64 {
	return FORKID_ZERO
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

// func (s *ForkIdState) checkIfForkIDAffectsPreviousBatches(ctx context.Context, newForkID ForkIDInterval, dbTx pgx.Tx) {
// 	latestBatchNumber, err := s.state.GetLastBatchNumber(ctx, dbTx)
// 	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
// 		log.Error("error getting last batch number. Error: ", err)
// 		rollbackErr := dbTx.Rollback(ctx)
// 		if rollbackErr != nil {
// 			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
// 			return rollbackErr
// 		}
// 		return err
// 	}
// }
