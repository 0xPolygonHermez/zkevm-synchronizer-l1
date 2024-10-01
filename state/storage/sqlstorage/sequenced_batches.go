package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	sequencedBatchesTable     = "sequenced_batches"
	sqlSelectSequencedBatches = "SELECT from_batch_num, to_batch_num, fork_id, timestamp, block_num, l1_info_root, received_at, source "
)

type SequencedBatches = entities.SequencedBatches
type SequencedBatchesSlice = entities.SequencesBatchesSlice

// AddForkID adds a new forkID to the storage
func (p *SqlStorage) AddSequencedBatches(ctx context.Context, sequence *SequencedBatches, dbTx dbTxType) error {
	sql := "INSERT INTO " + p.BuildTableName(sequencedBatchesTable) + " (from_batch_num, to_batch_num, fork_id,timestamp,block_num, l1_info_root, received_at, source) VALUES ($1, $2, $3, $4,$5, $6,$7,$8);"
	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, sql, sequence.FromBatchNumber, sequence.ToBatchNumber, sequence.ForkID, sequence.Timestamp.UTC(),
		sequence.L1BlockNumber, sequence.L1InfoRoot.String(), sequence.ReceivedAt.UTC(), sequence.Source)
	return translateSqlError(err, fmt.Sprintf("AddSequencedBatches %d", sequence.Key()))
}

func (p *SqlStorage) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatches, error) {
	sql := sqlSelectSequencedBatches +
		"FROM " + p.BuildTableName(sequencedBatchesTable) + " " +
		"WHERE  $1 >= from_batch_num  AND $1 <= to_batch_num " +
		"ORDER BY from_batch_num DESC LIMIT 1;"
	sequences, err := p.querySequences(ctx, fmt.Sprintf("GetSequenceByBatchNumber %d", batchNumber), sql, getSqlTx(dbTx), batchNumber)
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

func (p *SqlStorage) GetLatestSequence(ctx context.Context, dbTx dbTxType) (*SequencedBatches, error) {
	sql := sqlSelectSequencedBatches +
		"FROM " + p.BuildTableName(sequencedBatchesTable) + "  " +
		"ORDER BY from_batch_num DESC LIMIT 1;"
	sequences, err := p.querySequences(ctx, "GetLatestSequence", sql, getSqlTx(dbTx))
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

func (p *SqlStorage) GetSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatchesSlice, error) {
	sql := sqlSelectSequencedBatches +
		"FROM " + p.BuildTableName(sequencedBatchesTable) + "  " +
		"WHERE    from_batch_num >= $1 OR to_batch_num >= $1 " +
		"ORDER BY from_batch_num;"
	seq, err := p.querySequences(ctx, fmt.Sprintf("GetSequencesGreatestOrEqualBatchNumber %d", batchNumber), sql, getSqlTx(dbTx), batchNumber)
	if err != nil {
		return nil, err
	}
	if seq == nil {
		return nil, nil
	}
	seqSlice := SequencedBatchesSlice(*seq)
	return &seqSlice, nil
}

func (p *SqlStorage) DeleteSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) error {
	sql := "DELETE " +
		"FROM " + p.BuildTableName(sequencedBatchesTable) + "  " +
		"WHERE    from_batch_num >= $1 OR to_batch_num >= $1;"
	e := p.getExecQuerier(getSqlTx(dbTx))
	if _, err := e.ExecContext(ctx, sql, batchNumber); err != nil {
		return err
	}
	return nil
}

func (p *SqlStorage) querySequences(ctx context.Context, desc string, sql string, dbTx *sql.Tx, args ...interface{}) (*[]SequencedBatches, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sequences []SequencedBatches
	for rows.Next() {
		sequence, err := scanSequence(rows)
		if err != nil {
			err = translateSqlError(err, desc)
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
