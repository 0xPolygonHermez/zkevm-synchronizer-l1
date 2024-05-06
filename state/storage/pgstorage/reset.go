package pgstorage

import (
	"context"
	"log"
)

func (p *PostgresStorage) Reset(ctx context.Context, blockNumber uint64, dbTx dbTxType) error {
	log.Fatal("Not implemented")
	return nil
}
