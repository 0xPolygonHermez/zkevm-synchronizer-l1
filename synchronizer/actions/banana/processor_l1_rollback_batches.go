package banana

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type ProcessorRollbackBatches struct {
	actions.ProcessorBase[ProcessorRollbackBatches]
	state stateOnSequencedBatchesInterface
}

// NewProcessorRollbackBatches returns instance of a processor for RollbackBatchesOrder
func NewProcessorRollbackBatches(state stateOnSequencedBatchesInterface) *ProcessorRollbackBatches {
	return &ProcessorRollbackBatches{
		ProcessorBase: actions.ProcessorBase[ProcessorRollbackBatches]{
			SupportedEvent:    []etherman.EventOrder{etherman.RollbackBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyBanana},
		state: state,
	}
}

func (g *ProcessorRollbackBatches) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error {
	if l1Block == nil || len(l1Block.RollbackBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}

	err := g.ProcessRollbackBatches(ctx, forkId, l1Block.RollbackBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorRollbackBatches) ProcessRollbackBatches(ctx context.Context, forkId ForkIdType, rollbackData etherman.RollbackBatchesData, blockNumber uint64, l1BlockTimestamp time.Time, dbTx stateTxType) error {

	return fmt.Errorf("not implemented")
}
