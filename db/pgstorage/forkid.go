package pgstorage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddForkID(ctx context.Context, forkID ForkIDInterval, dbTx pgx.Tx) error {
	const addForkIDSQL = "INSERT INTO sync.fork_id (from_batch_num, to_batch_num, fork_id, version, block_num) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (fork_id) DO UPDATE SET block_num = $5 WHERE sync.fork_id.fork_id = $3;"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addForkIDSQL, forkID.FromBatchNumber, forkID.ToBatchNumber, forkID.ForkId, forkID.Version, forkID.BlockNumber)
	return err
}

// GetForkIDs get all the forkIDs stored
func (p *PostgresStorage) GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]ForkIDInterval, error) {
	const getForkIDsSQL = "SELECT from_batch_num, to_batch_num, fork_id, version, block_num FROM sync.fork_id ORDER BY from_batch_num ASC"
	q := p.getExecQuerier(dbTx)

	rows, err := q.Query(ctx, getForkIDsSQL)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	forkIDs := make([]ForkIDInterval, 0, len(rows.RawValues()))

	for rows.Next() {
		var forkID ForkIDInterval
		if err := rows.Scan(
			&forkID.FromBatchNumber,
			&forkID.ToBatchNumber,
			&forkID.ForkId,
			&forkID.Version,
			&forkID.BlockNumber,
		); err != nil {
			return forkIDs, err
		}
		forkIDs = append(forkIDs, forkID)
	}
	return forkIDs, err
}

// UpdateForkID updates the forkID stored in db
func (p *PostgresStorage) UpdateForkID(ctx context.Context, forkID ForkIDInterval, dbTx pgx.Tx) error {
	const updateForkIDSQL = "UPDATE sync.fork_id SET to_batch_num = $1 WHERE fork_id = $2"
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, updateForkIDSQL, forkID.ToBatchNumber, forkID.ForkId); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) uint64 {
	return 0
}
