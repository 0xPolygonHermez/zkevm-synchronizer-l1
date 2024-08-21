package etrog

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/ethereum/go-ethereum/common"
)

// ProcessorL1UpdateEtrogSequence implements L1EventProcessor
type ProcessorL1UpdateEtrogSequence struct {
	actions.ProcessorBase[ProcessorL1UpdateEtrogSequence]
	state stateOnSequencedBatchesInterface
}

// NewProcessorL1UpdateEtrogSequence returns instance of a processor for UpdateEtrogSequenceOrder
func NewProcessorL1UpdateEtrogSequence(state stateOnSequencedBatchesInterface,
) *ProcessorL1UpdateEtrogSequence {
	return &ProcessorL1UpdateEtrogSequence{
		ProcessorBase: actions.ProcessorBase[ProcessorL1UpdateEtrogSequence]{
			SupportedEvent:    []etherman.EventOrder{etherman.UpdateEtrogSequenceOrder},
			SupportedForkdIds: &actions.ForksIdOnlyEtrog},
		state: state,
	}
}

// Process process event
func (p *ProcessorL1UpdateEtrogSequence) Process(ctx context.Context, forkId ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx stateTxType) error {
	if l1Block == nil || l1Block.UpdateEtrogSequence.BatchNumber == 0 {
		return actions.ErrInvalidParams
	}
	err := p.processUpdateEtrogSequence(ctx, forkId, order, l1Block.UpdateEtrogSequence, l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

func (p *ProcessorL1UpdateEtrogSequence) processUpdateEtrogSequence(ctx context.Context, forkId ForkIdType, order etherman.Order, updateEtrogSequence etherman.UpdateEtrogSequence, blockNumber uint64, l1BlockTimestamp time.Time, dbTx stateTxType) error {
	l1inforoot := common.Hash(updateEtrogSequence.EtrogSequenceData.ForcedGlobalExitRoot)
	seq := SequenceOfBatches{}
	seq.Sequence = *entities.NewSequencedBatches(
		updateEtrogSequence.BatchNumber, updateEtrogSequence.BatchNumber,
		blockNumber, uint64(forkId),
		l1BlockTimestamp, time.Now(),
		l1inforoot, string(order.Name))
	ethSeqBatch := etherman.SequencedBatch{
		BatchNumber:       updateEtrogSequence.BatchNumber,
		L1InfoRoot:        &l1inforoot,
		SequencerAddr:     updateEtrogSequence.SequencerAddr,
		TxHash:            updateEtrogSequence.TxHash,
		EtrogSequenceData: updateEtrogSequence.EtrogSequenceData,
	}

	batch := entities.NewVirtualBatchFromL1(blockNumber, seq.Sequence.FromBatchNumber, uint64(forkId), ethSeqBatch)

	seq.Batches = append(seq.Batches, batch)

	return p.state.OnSequencedBatchesOnL1(ctx, seq, dbTx)
}
