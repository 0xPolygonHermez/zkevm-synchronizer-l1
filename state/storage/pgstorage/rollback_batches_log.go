package pgstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type RollbackBatchesLogEntry = entities.RollbackBatchesLogEntry

func (p *PostgresStorage) AddRollbackBatchesLogEntry(ctx context.Context, entry *RollbackBatchesLogEntry, dbTx dbTxType) error {
	const sql = `
	INSERT INTO sync.rollback_batches_log 
	(id, block_num, last_batch_number,last_batch_acc_input_hash,
	l1_event_at,received_at, undo_first_block_num,description, sequences_deleted, sync_version)    
	VALUES ($1, $2, $3, $4,$5, $6,$7,$8, $9, $10);
	`
	if entry == nil {
		return fmt.Errorf("AddRollbackBatchesLogEntry: entry is nil err:%w", entities.ErrBadParams)
	}
	id := entry.ID()
	seqJson, err := json.Marshal(entry.SequencesDeleted)
	if err != nil {
		return fmt.Errorf("AddRollbackBatchesLogEntry: error marshalling sequencesDeleted err:%w", err)
	}
	e := p.getExecQuerier(getPgTx(dbTx))
	_, err = e.Exec(ctx, sql,
		id.String(), entry.BlockNumber, entry.LastBatchNumber, entry.LastBatchAccInputHash.String(),
		entry.L1EventAt, entry.ReceivedAt, entry.UndoFirstBlockNumber, entry.Description, seqJson, entry.SyncVersion())
	return translatePgxError(err, fmt.Sprintf("AddRollbackBatchesLogEntry %s", id.String()))

}

func (p *PostgresStorage) GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx context.Context, l1BlockNumber uint64, dbTx dbTxType) ([]RollbackBatchesLogEntry, error) {
	const sql = `
	SELECT id, block_num, last_batch_number,last_batch_acc_input_hash,
	l1_event_at,received_at, undo_first_block_num,description, sequences_deleted, sync_version
	FROM sync.rollback_batches_log
	WHERE block_num >= $1
	ORDER BY block_num;
	`
	return p.queryRollbackBatchesLogEntries(ctx, fmt.Sprintf("GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber %d", l1BlockNumber), sql, getPgTx(dbTx), l1BlockNumber)
}

func (p *PostgresStorage) queryRollbackBatchesLogEntries(ctx context.Context, desc string, sql string, dbTx pgx.Tx, args ...interface{}) ([]RollbackBatchesLogEntry, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []RollbackBatchesLogEntry
	for rows.Next() {
		var entry RollbackBatchesLogEntry
		var sequencesDeletedJson string
		var id string
		var syncVersion string
		var lastBatchAccInputHash string
		err := rows.Scan(&id, &entry.BlockNumber, &entry.LastBatchNumber, &lastBatchAccInputHash,
			&entry.L1EventAt, &entry.ReceivedAt, &entry.UndoFirstBlockNumber, &entry.Description, &sequencesDeletedJson, &syncVersion)
		if err != nil {
			return nil, err
		}
		entry.LastBatchAccInputHash = common.HexToHash(lastBatchAccInputHash)
		var sequencesDeleted []entities.SequencedBatches
		err = json.Unmarshal([]byte(sequencesDeletedJson), &sequencesDeleted)
		if err != nil {
			return nil, err
		}
		entry.SequencesDeleted = sequencesDeleted
		if entry.ID().String() != id {
			return nil, fmt.Errorf("queryRollbackBatchesLogEntries: id mismatch database:%s calculated:%s", id, entry.ID().String())
		}
		entry.SetSyncVersion(syncVersion)
		entries = append(entries, entry)
	}
	return entries, nil
}
