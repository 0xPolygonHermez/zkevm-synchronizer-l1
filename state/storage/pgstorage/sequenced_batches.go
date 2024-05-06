package pgstorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type SequencedBatches = entities.SequencedBatches

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddSequencedBatches(ctx context.Context, sequence SequencedBatches, dbTx dbTxType) error {
	const sql = "INSERT INTO sync.sequenced_batches (from_batch_num, to_batch_num, timestamp,block_num, l1_info_root) VALUES ($1, $2, $3, $4,$5);"
	e := p.getExecQuerier(getPgTx(dbTx))
	_, err := e.Exec(ctx, sql, sequence.FromBatchNumber, sequence.ToBatchNumber, sequence.Timestamp, sequence.L1BlockNumber, sequence.L1InfoRoot.String())
	return err
}

func (p *PostgresStorage) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatches, error) {
	const sql = `SELECT from_batch_num, to_batch_num, timestamp,block_num, l1_info_root FROM sync.sequenced_batches 
		WHERE  $1 >= from_batch_num  AND $1 <= to_batch_num 
		ORDER BY block_num DESC LIMIT 1;`
	e := p.getExecQuerier(getPgTx(dbTx))
	row := e.QueryRow(ctx, sql, batchNumber)
	sequence := &SequencedBatches{}
	var l1InfoRootStr string
	err := row.Scan(&sequence.FromBatchNumber, &sequence.ToBatchNumber, &sequence.Timestamp, &sequence.L1BlockNumber, &l1InfoRootStr)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	sequence.L1InfoRoot = common.HexToHash(l1InfoRootStr)
	return sequence, nil
}
