package banana

import (
	"context"
	"time"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/ethereum/go-ethereum/common"
)

type ProcessorL1SequenceBatchesBanana struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesBanana]
	state stateOnSequencedBatchesInterface
}

// NewProcessorL1SequenceBatchesElderberry returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatchesBanana(state stateOnSequencedBatchesInterface) *ProcessorL1SequenceBatchesBanana {
	return &ProcessorL1SequenceBatchesBanana{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesBanana]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyBanana},
		state: state,
	}
}

func (g *ProcessorL1SequenceBatchesBanana) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}

	if len(l1Block.SequencedBatches[order.Pos]) == 0 {
		log.Warnf("No sequenced batches for position")
		return nil
	}
	err := g.ProcessSequenceBatches(ctx, forkId, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1SequenceBatchesBanana) ProcessSequenceBatches(ctx context.Context, forkId ForkIdType, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx stateTxType) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	l1inforoot := common.Hash{}
	if sequencedBatches[0].L1InfoRoot != nil {
		l1inforoot = *sequencedBatches[0].L1InfoRoot
	}
	seq := SequenceOfBatches{}
	seqSource := string(etherman.SequenceBatchesOrder)
	if sequencedBatches[0].Metadata != nil {
		seqSource = seqSource + "/" + sequencedBatches[0].Metadata.RollupFlavor + "/" + sequencedBatches[0].Metadata.ForkName
	}
	seqTimeStamp := time.Unix(int64(sequencedBatches[0].BananaData.MaxSequenceTimestamp), 0)
	seq.Sequence = *entities.NewSequencedBatches(
		sequencedBatches[0].BatchNumber, sequencedBatches[len(sequencedBatches)-1].BatchNumber,
		blockNumber, uint64(forkId),
		seqTimeStamp, time.Now(),
		l1inforoot, seqSource)

	for _, sequencedBatch := range sequencedBatches {
		virtualBatch := entities.NewVirtualBatchFromL1(blockNumber, seq.Sequence.FromBatchNumber,
			seq.Sequence.ForkID, sequencedBatch)
		seq.Batches = append(seq.Batches, virtualBatch)
	}
	return p.state.OnSequencedBatchesOnL1(ctx, seq, dbTx)
}
