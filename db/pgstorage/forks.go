package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

const (
	// FORKID_BLUEBERRY is the fork id 4
	FORKID_BLUEBERRY = 4
	// FORKID_DRAGONFRUIT is the fork id 5
	FORKID_DRAGONFRUIT = 5
	// FORKID_INCABERRY is the fork id 6
	FORKID_INCABERRY = 6
	// FORKID_ETROG is the fork id 7
	FORKID_ETROG = 7
	// FORKID_ELDERBERRY is the fork id 8
	FORKID_ELDERBERRY = 8
	// FORKID_ZERO is the fork id 0
	FORKID_ZERO = 0
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}

// UpdateForkIDIntervalsInMemory updates the forkID intervals in memory
func (s *PostgresStorage) UpdateForkIDIntervalsInMemory(intervals []ForkIDInterval) {
}

// AddForkIDInterval updates the forkID intervals
func (s *PostgresStorage) AddForkIDInterval(ctx context.Context, newForkID ForkIDInterval, dbTx pgx.Tx) error {
	return nil
}

// GetForkIDByBatchNumber returns the fork id for a given batch number
func (s *PostgresStorage) GetForkIDByBatchNumber(batchNumber uint64) uint64 {
	return FORKID_ZERO
}

// GetForkIDByBlockNumber returns the fork id for a given block number
func (s *PostgresStorage) GetForkIDByBlockNumber(blockNumber uint64) uint64 {
	return FORKID_ZERO
}
