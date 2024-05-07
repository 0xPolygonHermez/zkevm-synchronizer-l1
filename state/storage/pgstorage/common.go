package pgstorage

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type L1Block = entities.L1Block
type VirtualBatch = entities.VirtualBatch
type ForkIDInterval = entities.ForkIDInterval
type dbTxType = entities.Tx

func getPgTx(tx dbTxType) pgx.Tx {
	res, ok := tx.(pgx.Tx)
	if !ok {
		return nil
	}
	return res
}

const UniqueViolationErr = "23505"
const ForeignKeyViolationErr = "23503"

func translatePgxError(err error, contextDescription string) error {
	if err == nil {
		return nil
	}
	newErr := err
	if err == pgx.ErrNoRows {
		newErr = entities.ErrNotFound
	}
	pgErr, ok := err.(*pgconn.PgError)
	if ok {
		switch pgErr.Code {
		case UniqueViolationErr:
			newErr = fmt.Errorf("%w : pgError:%w ", entities.ErrAlreadyExists, err)
		case ForeignKeyViolationErr:
			newErr = fmt.Errorf("%w : pgError:%w ", entities.ErrForeignKeyViolation, err)
		}
	}
	return fmt.Errorf("storage error: %s: Err: %w", contextDescription, newErr)
}
