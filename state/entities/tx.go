package entities

//go:generate bash -c "rm -Rf mocks"
//go:generate mockery --all --case snake --dir . --output ./mocks --outpkg mock_entities --disable-version-string --with-expecter

import (
	"context"
)

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	AddRollbackCallback(func())
}
