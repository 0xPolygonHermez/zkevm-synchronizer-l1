package pgstorage

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/jackc/pgx/v4"
)

type L1Block = entities.L1Block
type ForkIDInterval = entities.ForkIDInterval
type dbTxType = entities.Tx

func getPgTx(tx dbTxType) pgx.Tx {
	res, ok := tx.(pgx.Tx)
	if !ok {
		return nil
	}
	return res
}
