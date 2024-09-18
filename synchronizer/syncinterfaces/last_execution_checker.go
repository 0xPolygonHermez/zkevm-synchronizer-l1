package syncinterfaces

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type LastExecutionChecker interface {
	StartingSynchronization(ctx context.Context, dbTx entities.Tx) error
}
