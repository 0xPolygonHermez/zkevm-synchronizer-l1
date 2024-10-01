package pgstorage

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type SequencedBatches = entities.SequencedBatches
type SequencedBatchesSlice = entities.SequencesBatchesSlice

// AddForkID adds a new forkID to the storage
func (p *PostgresStorage) AddSequencedBatches(ctx context.Context, sequence *SequencedBatches, dbTx dbTxType) error {
	const sql = "INSERT INTO sync.sequenced_batches (from_batch_num, to_batch_num, fork_id,timestamp,block_num, l1_info_root, received_at, source) VALUES ($1, $2, $3, $4,$5, $6,$7,$8);"
	e := p.getExecQuerier(getPgTx(dbTx))
	_, err := e.Exec(ctx, sql, sequence.FromBatchNumber, sequence.ToBatchNumber, sequence.ForkID, sequence.Timestamp,
		sequence.L1BlockNumber, sequence.L1InfoRoot.String(), sequence.ReceivedAt, sequence.Source)
	return translatePgxError(err, fmt.Sprintf("AddSequencedBatches %d", sequence.Key()))
}

func (p *PostgresStorage) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatches, error) {
	const sql = `SELECT from_batch_num, to_batch_num,fork_id, timestamp,block_num, l1_info_root,received_at,source FROM sync.sequenced_batches 
		WHERE  $1 >= from_batch_num  AND $1 <= to_batch_num 
		ORDER BY from_batch_num DESC LIMIT 1;`
	sequences, err := p.querySequences(ctx, fmt.Sprintf("GetSequenceByBatchNumber %d", batchNumber), sql, getPgTx(dbTx), batchNumber)
	if err != nil {
		return nil, err
	}
	if sequences == nil || len(*sequences) == 0 {
		return nil, nil
	}
	if len(*sequences) > 1 {
		return nil, fmt.Errorf("more than one sequence found for batch number %d", batchNumber)
	}
	return &(*sequences)[0], nil
}

func (p *PostgresStorage) GetLatestSequence(ctx context.Context, dbTx dbTxType) (*SequencedBatches, error) {
	const sql = `SELECT from_batch_num, to_batch_num,fork_id, timestamp,block_num, l1_info_root,received_at,source FROM sync.sequenced_batches 
		ORDER BY block_num DESC LIMIT 1;`
	sequences, err := p.querySequences(ctx, "GetLatestSequence", sql, getPgTx(dbTx))
	if err != nil {
		return nil, err
	}
	if sequences == nil || len(*sequences) == 0 {
		return nil, nil
	}
	if len(*sequences) > 1 {
		return nil, fmt.Errorf("more than one row for GetLatestSequence!??!?")
	}
	return &(*sequences)[0], nil
}

func (p *PostgresStorage) GetSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatchesSlice, error) {
	const sql = `SELECT from_batch_num, to_batch_num,fork_id, timestamp,block_num, l1_info_root,received_at,source FROM sync.sequenced_batches 
		WHERE    from_batch_num >= $1 OR to_batch_num >= $1
		ORDER BY from_batch_num;`
	seq, err := p.querySequences(ctx, fmt.Sprintf("GetSequencesGreatestOrEqualBatchNumber %d", batchNumber), sql, getPgTx(dbTx), batchNumber)
	if err != nil {
		return nil, err
	}
	if seq == nil {
		return nil, nil
	}
	seqSlice := SequencedBatchesSlice(*seq)
	return &seqSlice, nil
}

func (p *PostgresStorage) DeleteSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) error {
	const sql = `DELETE FROM  sync.sequenced_batches 
	WHERE    from_batch_num >= $1 OR to_batch_num >= $1;`
	e := p.getExecQuerier(getPgTx(dbTx))
	if _, err := e.Exec(ctx, sql, batchNumber); err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) querySequences(ctx context.Context, desc string, sql string, dbTx pgx.Tx, args ...interface{}) (*[]SequencedBatches, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []SequencedBatches
	for rows.Next() {
		sequence, err := scanSequence(rows)
		if err != nil {
			err = translatePgxError(err, desc)
			return nil, err
		}
		sequences = append(sequences, sequence)
	}
	if len(sequences) == 0 {
		return nil, nil
	}
	return &sequences, nil
}

func scanSequence(row pgx.Row) (SequencedBatches, error) {
	var l1InfoRootStr string
	sequence := SequencedBatches{}
	if err := row.Scan(&sequence.FromBatchNumber, &sequence.ToBatchNumber, &sequence.ForkID, &sequence.Timestamp,
		&sequence.L1BlockNumber, &l1InfoRootStr, &sequence.ReceivedAt, &sequence.Source); err != nil {
		return sequence, err
	}
	sequence.L1InfoRoot = common.HexToHash(l1InfoRootStr)
	return sequence, nil
}
