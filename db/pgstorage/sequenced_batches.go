package pgstorage

import (
	"context"
	"errors"
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

func (p *PostgresStorage) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*SequencedBatches, error) {
	const sql = `SELECT from_batch_num, to_batch_num, timestamp,block_num FROM sync.sequenced_batches 
		WHERE  from_batch_num >= $1 AND to_batch_num <= $1
		ORDER BY block_num DESC LIMIT 1;`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, sql, batchNumber)
	sequence := &SequencedBatches{}
	err := row.Scan(&sequence.FromBatchNumber, &sequence.ToBatchNumber, &sequence.Timestamp, &sequence.L1BlockNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return sequence, nil
}
