package etrog

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessorL1SequenceBatchesEtrog implements L1EventProcessor
type ProcessorL1SequenceBatchesEtrog struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]
	state stateOnSequencedBatchesInterface
}

// NewProcessorL1SequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatches(state stateOnSequencedBatchesInterface) *ProcessorL1SequenceBatchesEtrog {
	return &ProcessorL1SequenceBatchesEtrog{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyEtrog},
		state: state,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatchesEtrog) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx stateTxType) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.ProcessSequenceBatches(ctx, forkId, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1SequenceBatchesEtrog) ProcessSequenceBatches(ctx context.Context, forkId ForkIdType, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx stateTxType) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	l1inforoot := common.Hash{}
	if sequencedBatches[0].L1InfoRoot != nil {
		l1inforoot = *sequencedBatches[0].L1InfoRoot
	}
	seq := SequenceOfBatches{}
	seq.Sequence = SequencedBatches{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
		L1BlockNumber:   blockNumber,
		Timestamp:       l1BlockTimestamp,
		L1InfoRoot:      l1inforoot,
		ForkID:          uint64(forkId),
		Source:          string(etherman.SequenceBatchesOrder),
	}
	for _, sequencedBatch := range sequencedBatches {
		virtualBatch := entities.NewVirtualBatchFromL1(blockNumber, seq.Sequence.FromBatchNumber,
			seq.Sequence.ForkID, sequencedBatch)
		seq.Batches = append(seq.Batches, virtualBatch)
	}
	return p.state.OnSequencedBatchesOnL1(ctx, seq, dbTx)
}
