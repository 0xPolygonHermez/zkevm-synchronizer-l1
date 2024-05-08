package entities

import (
	"context"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	AddRollbackCallback(func())
}
