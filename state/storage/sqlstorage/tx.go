package sqlstorage

import (
	"context"
	"database/sql"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type TxCallbackType = entities.TxCallbackType

type stateImplementationTx struct {
	*sql.Tx
	rollbackCallbacks []TxCallbackType
	commitCallbacks   []TxCallbackType
}

func NewTxImpl(tx *sql.Tx) dbTxType {
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
	err := s.Tx.Commit()
	for _, cb := range s.commitCallbacks {
		cb(s, err)
	}
	return err
}

func (s *stateImplementationTx) Rollback(ctx context.Context) error {
	err := s.Tx.Rollback()
	for _, cb := range s.rollbackCallbacks {
		cb(s, err)
	}
	return err
}

func getSqlTx(tx dbTxType) *sql.Tx {
	res, ok := tx.(*stateImplementationTx)
	if !ok {
		return nil
	}
	return res.Tx
}
