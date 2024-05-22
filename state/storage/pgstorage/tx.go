package pgstorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/jackc/pgx/v4"
)

type Tx interface {
	pgx.Tx
	entities.Tx
}

type TxCallbackType = entities.TxCallbackType

type stateImplementationTx struct {
	pgx.Tx
	rollbackCallbacks []TxCallbackType
	commitCallbacks   []TxCallbackType
}

func NewTxImpl(tx pgx.Tx) Tx {
	return &stateImplementationTx{
		Tx: tx,
	}
}

func (s *stateImplementationTx) AddRollbackCallback(cb TxCallbackType) {
	s.rollbackCallbacks = append(s.rollbackCallbacks, cb)
}
func (s *stateImplementationTx) AddCommitCallback(cb TxCallbackType) {
	s.commitCallbacks = append(s.commitCallbacks, cb)
}

func (s *stateImplementationTx) Commit(ctx context.Context) error {
	err := s.Tx.Commit(ctx)
	for _, cb := range s.commitCallbacks {
		cb(s, err)
	}
	return err
}

func (s *stateImplementationTx) Rollback(ctx context.Context) error {
	err := s.Tx.Rollback(ctx)
	for _, cb := range s.rollbackCallbacks {
		cb(s, err)
	}
	return err
}
