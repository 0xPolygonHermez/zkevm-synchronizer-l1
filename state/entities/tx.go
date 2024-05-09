package entities

import (
	"context"
)

type TxCallbackType func(Tx, error)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	AddRollbackCallback(TxCallbackType)
	AddCommitCallback(TxCallbackType)
}
