package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Tx interface {
	pgx.Tx
	AddRollbackCallback(func())
}

type stateImplementationTx struct {
	pgx.Tx
	rollbackCallbacks []func()
}

func (s *stateImplementationTx) AddRollbackCallback(cb func()) {
	s.rollbackCallbacks = append(s.rollbackCallbacks, cb)
}
func (tx *stateImplementationTx) Rollback(ctx context.Context) error {
	for _, cb := range tx.rollbackCallbacks {
		cb()
	}
	return tx.Tx.Rollback(ctx)
}
