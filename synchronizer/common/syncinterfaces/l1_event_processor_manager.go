package syncinterfaces

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type L1EventProcessorManager interface {
	Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, block *etherman.Block, dbTx pgx.Tx) error
	Get(forkId actions.ForkIdType, event etherman.EventOrder) actions.L1EventProcessor
}
