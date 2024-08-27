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

type stateVerifyL1InfoTreeInterface interface {
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx stateTxType) (*L1InfoTreeLeaf, error)
}

type ProcessorUpdateL1InfoTreeV2 struct {
	actions.ProcessorBase[ProcessorUpdateL1InfoTreeV2]
	state stateVerifyL1InfoTreeInterface
}

// ProcessorUpdateL1InfoTreeV2 returns instance of a processor for UpdateL1InfoTreeV2Order
func NewProcessorUpdateL1InfoTreeV2(state stateVerifyL1InfoTreeInterface) *ProcessorUpdateL1InfoTreeV2 {
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
	stateLeaf, err := p.state.GetL1InfoLeafPerIndex(ctx, data.LeafCount-1, dbTx)
	if err != nil {
		log.Errorf("error getting the state leaf. Error: %v", err)
		return err
	}
	err = compareL1InfoTreeLeaf(*stateLeaf, data)
	if err != nil {
		log.Errorf("error comparing the state leaf. Error: %v", err)
		return err
	}
	return nil
}

func compareL1InfoTreeLeaf(leaf L1InfoTreeLeaf, event etherman.L1InfoTreeV2Data) error {
	if leaf.L1InfoTreeRoot != event.CurrentL1InfoRoot {
		return fmt.Errorf("L1InfoTreeRoot mismatch: %v != %v", leaf.L1InfoTreeRoot, event.CurrentL1InfoRoot)
	}
	return nil
}
