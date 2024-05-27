package pgstorage

import (
	"context"
)

// ResetToL1BlockNumber resets the state to a block for the given DB tx
func (p *PostgresStorage) ResetToL1BlockNumber(ctx context.Context, firstBlockNumberToKeep uint64, dbTx dbTxType) error {
	e := p.getExecQuerier(getPgTx(dbTx))
	const resetSQL = "DELETE FROM sync.block WHERE block_num > $1"
	if _, err := e.Exec(ctx, resetSQL, firstBlockNumberToKeep); err != nil {
		return err
	}
	return nil
}
