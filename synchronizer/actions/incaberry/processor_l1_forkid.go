package incaberry

import (
	"context"
	"math"

	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type stateForkIdInterface interface {
	AddForkID(ctx context.Context, newForkID entities.ForkIDInterval, dbTx entities.Tx) error
}

// ProcessorForkId implements L1EventProcessor
type ProcessorForkId struct {
	actions.ProcessorBase[ProcessorForkId]
	stateForkId stateForkIdInterface
}

// NewProcessorForkId returns instance of a processor for ForkIDsOrder
func NewProcessorForkId(stateForkId stateForkIdInterface) *ProcessorForkId {
	return &ProcessorForkId{
		ProcessorBase: actions.ProcessorBase[ProcessorForkId]{
			SupportedEvent:    []etherman.EventOrder{etherman.ForkIDsOrder},
			SupportedForkdIds: &actions.ForksIdAll,
		},
		stateForkId: stateForkId,
	}
}

// Process process event
func (p *ProcessorForkId) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error {
	return p.processForkID(ctx, l1Block.ForkIDs[order.Pos], l1Block.BlockNumber, dbTx)
}

func (s *ProcessorForkId) processForkID(ctx context.Context, forkID etherman.ForkID, blockNumber uint64, dbTx entities.Tx) error {
	fID := entities.ForkIDInterval{
		FromBatchNumber: forkID.BatchNumber + 1,
		ToBatchNumber:   math.MaxUint64,
		ForkId:          forkID.ForkID,
		Version:         forkID.Version,
		BlockNumber:     blockNumber,
	}

	// If forkID affects to a batch from the past. State must be reseted.
	log.Infof("ForkID: %d, synchronization must use the new forkID since batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	err := s.stateForkId.AddForkID(ctx, fID, dbTx)
	if err != nil {
		log.Errorf("Fails to add ForkID: %d from BatchNumer:%d", forkID.ForkID, forkID.BatchNumber+1)
		return err
	}
	//TODO: Figure out why it returns an error 	return fmt.Errorf("new ForkID detected, reseting synchronizarion")
	return nil
}
