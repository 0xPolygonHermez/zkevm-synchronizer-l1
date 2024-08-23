package banana

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type ProcessorUpdateL1InfoTreeV2 struct {
	actions.ProcessorBase[ProcessorUpdateL1InfoTreeV2]
	state stateOnSequencedBatchesInterface
}

// ProcessorUpdateL1InfoTreeV2 returns instance of a processor for UpdateL1InfoTreeV2Order
func NewProcessorUpdateL1InfoTreeV2(state stateOnSequencedBatchesInterface) *ProcessorUpdateL1InfoTreeV2 {
	return &ProcessorUpdateL1InfoTreeV2{
		ProcessorBase: actions.ProcessorBase[ProcessorUpdateL1InfoTreeV2]{
			SupportedEvent:    []etherman.EventOrder{etherman.UpdateL1InfoTreeV2Order},
			SupportedForkdIds: &actions.ForksIdAll},
		state: state,
	}
}

func (g *ProcessorUpdateL1InfoTreeV2) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error {
	if l1Block == nil || len(l1Block.L1InfoTreeV2) <= order.Pos {
		return actions.ErrInvalidParams
	}

	err := g.ProcessUpdateL1InfoTreeV2(ctx, forkId, l1Block.L1InfoTreeV2[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorUpdateL1InfoTreeV2) ProcessUpdateL1InfoTreeV2(ctx context.Context, forkId ForkIdType, data etherman.L1InfoTreeV2Data, blockNumber uint64, l1BlockTimestamp time.Time, dbTx stateTxType) error {
	log.Debugf("Processing UpdateL1InfoTreeV2: %s", data.String())
	return fmt.Errorf("not implemented")
}
