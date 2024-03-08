package pgstorage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func (p *PostgresStorage) Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	log.Fatal("Not implemented")
	return nil
}
