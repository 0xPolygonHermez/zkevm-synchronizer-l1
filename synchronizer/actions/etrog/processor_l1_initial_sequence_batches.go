package etrog

import (
	"context"
	"fmt"
	"time"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

// ProcessorL1InitialSequenceBatches implements L1EventProcessor
type ProcessorL1InitialSequenceBatches struct {
	actions.ProcessorBase[ProcessorL1InitialSequenceBatches]
	state stateOnSequencedBatchesInterface
}

// NewProcessorL1InitialSequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1InitialSequenceBatches(state stateOnSequencedBatchesInterface) *ProcessorL1InitialSequenceBatches {
	return &ProcessorL1InitialSequenceBatches{
		ProcessorBase: actions.ProcessorBase[ProcessorL1InitialSequenceBatches]{
			SupportedEvent:    []etherman.EventOrder{etherman.InitialSequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdEtrogElderberryBanana},
		state: state,
	}
}

// Process process event
func (g *ProcessorL1InitialSequenceBatches) Process(ctx context.Context, forkId ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx stateTxType) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.ProcessSequenceBatches(ctx, forkId, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1InitialSequenceBatches) ProcessSequenceBatches(ctx context.Context, forkId ForkIdType, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, dbTx stateTxType) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	if sequencedBatches[0].BatchNumber != 1 {
		return fmt.Errorf("invalid initial batch number, expected 1 , received %d", sequencedBatches[0].BatchNumber)
	}

	l1inforoot := sequencedBatches[0].EtrogSequenceData.ForcedGlobalExitRoot
	l1BlockTimestamp := time.Unix(int64(sequencedBatches[0].EtrogSequenceData.ForcedTimestamp), 0)
	seq := SequenceOfBatches{}

	seq.Sequence = *entities.NewSequencedBatches(
		sequencedBatches[0].BatchNumber, sequencedBatches[len(sequencedBatches)-1].BatchNumber,
		blockNumber, uint64(forkId),
		l1BlockTimestamp, time.Now(),
		l1inforoot, string(etherman.InitialSequenceBatchesOrder))

	for _, sequencedBatch := range sequencedBatches {
		virtualBatch := entities.NewVirtualBatchFromL1(blockNumber, seq.Sequence.FromBatchNumber,
			seq.Sequence.ForkID, sequencedBatch)
		seq.Batches = append(seq.Batches, virtualBatch)
	}

	return p.state.OnSequencedBatchesOnL1(ctx, seq, dbTx)
}
