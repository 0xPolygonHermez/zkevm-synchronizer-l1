package banana

import (
	"context"
	"time"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type ProcessorRollbackBatches struct {
	actions.ProcessorBase[ProcessorRollbackBatches]
	state stateOnRollbackBatchesInterface
}

// NewProcessorRollbackBatches returns instance of a processor for RollbackBatchesOrder
func NewProcessorRollbackBatches(state stateOnRollbackBatchesInterface) *ProcessorRollbackBatches {
	return &ProcessorRollbackBatches{
		ProcessorBase: actions.ProcessorBase[ProcessorRollbackBatches]{
			SupportedEvent: []etherman.EventOrder{etherman.RollbackBatchesOrder},
			// This event is processed for all forks, if the meaning or the way to execute depends on forkid need to be adapted
			SupportedForkdIds: &actions.ForksIdAll},
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
	log.Warnf("Processing RollbackBatches: %s", rollbackData.String())
	req := RollbackBatchesRequest{
		LastBatchNumber:       rollbackData.TargetBatch,
		LastBatchAccInputHash: rollbackData.AccInputHashToRollback,
		L1BlockNumber:         blockNumber,
		L1BlockTimestamp:      l1BlockTimestamp,
	}
	_, err := p.state.ExecuteRollbackBatches(ctx, req, dbTx)
	if err != nil {
		return err
	}

	return nil
}
