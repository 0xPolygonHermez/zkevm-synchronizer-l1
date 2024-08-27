package etrog

import (
	"context"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

// stateProcessorL1InfoTreeInterface interface required from state
type stateProcessorL1InfoTreeInterface interface {
	AddL1InfoTreeLeafAndAssignIndex(ctx context.Context, exitRoot *L1InfoTreeLeaf, dbTx stateTxType) (*L1InfoTreeLeaf, error)
}

// ProcessorL1InfoTreeUpdate implements L1EventProcessor for GlobalExitRootsOrder
type ProcessorL1InfoTreeUpdate struct {
	actions.ProcessorBase[ProcessorL1InfoTreeUpdate]
	state stateProcessorL1InfoTreeInterface
}

// NewProcessorL1InfoTreeUpdate new processor for GlobalExitRootsOrder
func NewProcessorL1InfoTreeUpdate(state stateProcessorL1InfoTreeInterface) *ProcessorL1InfoTreeUpdate {
	return &ProcessorL1InfoTreeUpdate{
		ProcessorBase: actions.ProcessorBase[ProcessorL1InfoTreeUpdate]{
			SupportedEvent:    []etherman.EventOrder{etherman.L1InfoTreeOrder},
			SupportedForkdIds: &actions.ForksIdAll},
		state: state}
}

// Process process event
func (p *ProcessorL1InfoTreeUpdate) Process(ctx context.Context, forkId ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx stateTxType) error {
	l1InfoTree := l1Block.L1InfoTree[order.Pos]
	leaf := L1InfoTreeLeaf{
		BlockNumber:       l1InfoTree.BlockNumber,
		MainnetExitRoot:   l1InfoTree.MainnetExitRoot,
		RollupExitRoot:    l1InfoTree.RollupExitRoot,
		GlobalExitRoot:    l1InfoTree.GlobalExitRoot,
		Timestamp:         l1InfoTree.Timestamp,
		PreviousBlockHash: l1InfoTree.PreviousBlockHash,
	}

	entry, err := p.state.AddL1InfoTreeLeafAndAssignIndex(ctx, &leaf, dbTx)
	if err != nil {
		log.Errorf("error storing the l1InfoTree(etrog). BlockNumber: %d, error: %v", l1Block.BlockNumber, err)
		return err
	}
	log.Infof("L1InfoTree(etrog) stored. BlockNumber: %d,GER:%s L1InfoTreeIndex: %d L1InfoRoot:%s event_data:%s",
		l1Block.BlockNumber, entry.GlobalExitRoot, entry.L1InfoTreeIndex, entry.L1InfoTreeRoot,
		l1InfoTree.String())
	return nil
}
