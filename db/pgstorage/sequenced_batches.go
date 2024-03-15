package pgstorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

type SequencedBatches struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	L1BlockNumber   uint64
	Timestamp       time.Time
}

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddSequencedBatches(ctx context.Context, sequence SequencedBatches, dbTx pgx.Tx) error {
	const sql = "INSERT INTO sync.sequenced_batches (from_batch_num, to_batch_num, timestamp,block_num) VALUES ($1, $2, $3, $4);"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, sql, sequence.FromBatchNumber, sequence.ToBatchNumber, sequence.Timestamp, sequence.L1BlockNumber)
	return err
}
