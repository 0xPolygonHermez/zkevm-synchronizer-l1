package sqlstorage

import (
	"context"
)

// ResetToL1BlockNumber resets the state to a block for the given DB tx
func (p *SqlStorage) ResetToL1BlockNumber(ctx context.Context, firstBlockNumberToKeep uint64, dbTx dbTxType) error {
	e := p.getExecQuerier(getSqlTx(dbTx))
	resetSQL := "DELETE FROM " + p.BuildTableName(blockTable) + " WHERE block_num > $1"
	if _, err := e.ExecContext(ctx, resetSQL, firstBlockNumberToKeep); err != nil {
		return err
	}
	return nil
}
