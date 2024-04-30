package etrog

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// ProcessorL1InitialSequenceBatches implements L1EventProcessor
type ProcessorL1InitialSequenceBatches struct {
	actions.ProcessorBase[ProcessorL1InitialSequenceBatches]
	state stateProcessSequenceBatches
}

// NewProcessorL1InitialSequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1InitialSequenceBatches(state stateProcessSequenceBatches) *ProcessorL1InitialSequenceBatches {
	return &ProcessorL1InitialSequenceBatches{
		ProcessorBase: actions.ProcessorBase[ProcessorL1InitialSequenceBatches]{
			SupportedEvent:    []etherman.EventOrder{etherman.InitialSequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyEtrogAndElderberry},
		state: state,
	}
}

// Process process event
func (g *ProcessorL1InitialSequenceBatches) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1InitialSequenceBatches) ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	if sequencedBatches[0].BatchNumber != 1 {
		return fmt.Errorf("invalid initial batch number, expected 1 , received %d", sequencedBatches[0].BatchNumber)
	}

	l1inforoot := sequencedBatches[0].PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot
	l1BlockTimestamp = time.Unix(int64(sequencedBatches[0].PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0)

	seq := pgstorage.SequencedBatches{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
		L1BlockNumber:   blockNumber,
		Timestamp:       l1BlockTimestamp,
		L1InfoRoot:      l1inforoot,
	}

	return p.state.AddSequencedBatches(ctx, seq, dbTx)
}
