package sqlstorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

const forkidTable = "fork_id"

type ForkIDInterval = entities.ForkIDInterval

// AddForkID adds a new forkID to the storage
func (p *SqlStorage) AddForkID(ctx context.Context, forkID ForkIDInterval, dbTx dbTxType) error {
	addForkIDSQL := "INSERT INTO " + p.BuildTableName(forkidTable) + " (from_batch_num, to_batch_num, fork_id, version, block_num) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (fork_id) DO UPDATE SET block_num = $5 WHERE fork_id.fork_id = $3;"
	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, addForkIDSQL, forkID.FromBatchNumber, forkID.ToBatchNumber, forkID.ForkId, forkID.Version, forkID.BlockNumber)
	err = translateSqlError(err, "AddForkID")
	return err
}

// GetForkIDs get all the forkIDs stored
func (p *SqlStorage) GetForkIDs(ctx context.Context, dbTx dbTxType) ([]ForkIDInterval, error) {
	getForkIDsSQL := "SELECT from_batch_num, to_batch_num, fork_id, version, block_num FROM " + p.BuildTableName(forkidTable) + " ORDER BY from_batch_num ASC"
	q := p.getExecQuerier(getSqlTx(dbTx))

	rows, err := q.QueryContext(ctx, getForkIDsSQL)
	if err != nil {
		return nil, translateSqlError(err, "GetForkIDs")
	}
	defer rows.Close()

	var forkIDs []ForkIDInterval

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
func (p *SqlStorage) UpdateForkID(ctx context.Context, forkID ForkIDInterval, dbTx dbTxType) error {
	updateForkIDSQL := "UPDATE " + p.BuildTableName(forkidTable) + " SET to_batch_num = $1 WHERE fork_id = $2"
	e := p.getExecQuerier(getSqlTx(dbTx))
	if _, err := e.ExecContext(ctx, updateForkIDSQL, forkID.ToBatchNumber, forkID.ForkId); err != nil {
		return translateSqlError(err, "UpdateForkID")
	}
	return nil
}

func (s *SqlStorage) GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) uint64 {
	return 0
}
