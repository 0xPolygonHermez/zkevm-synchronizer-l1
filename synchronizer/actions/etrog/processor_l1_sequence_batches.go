package etrog

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type stateProcessSequenceBatches interface {
	AddSequencedBatches(ctx context.Context, sequence pgstorage.SequencedBatches, dbTx pgx.Tx) error
}

// ProcessorL1SequenceBatchesEtrog implements L1EventProcessor
type ProcessorL1SequenceBatchesEtrog struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]
	state stateProcessSequenceBatches
}

// NewProcessorL1SequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatches(state stateProcessSequenceBatches) *ProcessorL1SequenceBatchesEtrog {
	return &ProcessorL1SequenceBatchesEtrog{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder, etherman.InitialSequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdAll},
		state: state,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatchesEtrog) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1SequenceBatchesEtrog) ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	seq := pgstorage.SequencedBatches{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
		L1BlockNumber:   blockNumber,
		Timestamp:       l1BlockTimestamp,
	}
	return p.state.AddSequencedBatches(ctx, seq, dbTx)
}
