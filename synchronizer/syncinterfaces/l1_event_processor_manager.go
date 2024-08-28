package syncinterfaces

import (
	"context"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type L1EventProcessorManager interface {
	Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, block *etherman.Block, dbTx entities.Tx) error
	Get(forkId actions.ForkIdType, event etherman.EventOrder) actions.L1EventProcessor
}
