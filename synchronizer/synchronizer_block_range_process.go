package synchronizer

import (
	"context"
	"errors"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1event_orders"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateBlockRangeProcessor interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	AddBlock(ctx context.Context, block *pgstorage.L1Block, dbTx pgx.Tx) error
}

// BlockRangeProcess is the struct that process the block range that implements syncinterfaces.BlockRangeProcessor
type BlockRangeProcess struct {
	state             stateBlockRangeProcessor
	stateForkId       StateForkIdQuerier
	l1EventProcessors syncinterfaces.L1EventProcessorManager
	flushIdManager    syncinterfaces.SynchronizerFlushIDManager
}

// NewBlockRangeProcessLegacy creates a new BlockRangeProcess
func NewBlockRangeProcessLegacy(
	state stateBlockRangeProcessor,
	stateForkId StateForkIdQuerier,
	l1EventProcessors syncinterfaces.L1EventProcessorManager,
	flushIdManager syncinterfaces.SynchronizerFlushIDManager,
) *BlockRangeProcess {
	return &BlockRangeProcess{
		state:             state,
		stateForkId:       stateForkId,
		l1EventProcessors: l1EventProcessors,
		flushIdManager:    flushIdManager,
	}
}

// ProcessBlockRangeSingleDbTx process the L1 events and stores the information in the db reusing same DbTx
func (s *BlockRangeProcess) ProcessBlockRangeSingleDbTx(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, finalizedBlockNumber uint64, storeBlocks syncinterfaces.ProcessBlockRangeL1BlocksMode, dbTx pgx.Tx) error {
	return s.internalProcessBlockRange(ctx, blocks, order, finalizedBlockNumber, storeBlocks, &dbTx)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) ProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, finalizedBlockNumber uint64) error {
	return s.internalProcessBlockRange(ctx, blocks, order, finalizedBlockNumber, syncinterfaces.StoreL1Blocks, nil)
}

func isBlockFinalized(blockNumber uint64, finalizedBlockNumber uint64) bool {
	return blockNumber <= finalizedBlockNumber
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) addBlock(ctx context.Context, block *etherman.Block, isFinalized bool, dbTx pgx.Tx) error {
	b := pgstorage.L1Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
		SyncVersion: zkevm_synchronizer_l1.Version,
		Checked:     isFinalized,
	}
	// Add block information
	return s.state.AddBlock(ctx, &b, dbTx)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) internalProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order,
	finalizedBlockNumber uint64,
	storeBlocksMode syncinterfaces.ProcessBlockRangeL1BlocksMode,
	dbTxExt *pgx.Tx) error {

	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		var dbTx pgx.Tx
		var err error
		if dbTxExt == nil {
			log.Debugf("Starting dbTx for BlockNumber:%d", blocks[i].BlockNumber)
			dbTx, err = s.state.BeginStateTransaction(ctx)
			if err != nil {
				return err
			}
		} else {
			dbTx = *dbTxExt
		}
		// Process event received from l1
		err = s.processBlock(ctx, blocks, i, dbTx, order, storeBlocksMode, finalizedBlockNumber)
		if err != nil {
			if dbTxExt == nil {
				// Rollback db transaction
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					if !errors.Is(rollbackErr, pgx.ErrTxClosed) {
						log.Errorf("error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
						return rollbackErr
					} else {
						log.Warnf("error rolling back state because is already closed. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
						return err
					}
				}
				return err
			}
			return err
		}
		if dbTxExt == nil {
			// Commit db transaction
			err = dbTx.Commit(ctx)
			if err != nil {
				log.Errorf("error committing state. BlockNumber: %d, Error: %v", blocks[i].BlockNumber, err)
			}
		}
	}
	return nil
}

func (s *BlockRangeProcess) processBlock(ctx context.Context, blocks []etherman.Block, blockIndex int, dbTx pgx.Tx, order map[common.Hash][]etherman.Order, storeBlock syncinterfaces.ProcessBlockRangeL1BlocksMode, finalizedBlockNumber uint64) error {
	var err error
	if storeBlock == syncinterfaces.StoreL1Blocks {
		err = s.addBlock(ctx, &blocks[blockIndex], isBlockFinalized(blocks[blockIndex].BlockNumber, finalizedBlockNumber), dbTx)
		if err != nil {
			log.Errorf("error adding block to db. BlockNumber: %d, error: %v", blocks[blockIndex].BlockNumber, err)
			return err
		}
	} else {
		log.Debugf("Skip storing block BlockNumber:%d", blocks[blockIndex].BlockNumber)
	}
	for _, element := range order[blocks[blockIndex].BlockHash] {
		err := s.processElement(ctx, element, blocks, blockIndex, dbTx)
		if err != nil {
			return err
		}
	}
	if s.flushIdManager != nil {
		log.Debug("Checking FlushID to commit L1 data to db")
		err = s.flushIdManager.CheckFlushID(dbTx)
		if err != nil {
			log.Errorf("error checking flushID. BlockNumber: %d, Error: %v", blocks[blockIndex].BlockNumber, err)
			return err
		}
	}
	return nil
}

func (s *BlockRangeProcess) processElement(ctx context.Context, element etherman.Order, blocks []etherman.Block, i int, dbTx pgx.Tx) error {
	batchSequence := l1event_orders.GetSequenceFromL1EventOrder(element.Name, &blocks[i], element.Pos)
	forkId := state.FORKID_ZERO
	if s.stateForkId != nil {
		if batchSequence != nil {
			forkId = s.stateForkId.GetForkIDByBatchNumber(ctx, batchSequence.FromBatchNumber, dbTx)
			log.Debug("EventOrder: ", element.Name, ". Batch Sequence: ", batchSequence, "forkId: ", forkId)
		} else {
			forkId = s.stateForkId.GetForkIDByBlockNumber(ctx, blocks[i].BlockNumber, dbTx)
			log.Debug("EventOrder: ", element.Name, ". BlockNumber: ", blocks[i].BlockNumber, "forkId: ", forkId)
		}
	}
	forkIdTyped := actions.ForkIdType(forkId)

	err := s.l1EventProcessors.Process(ctx, forkIdTyped, element, &blocks[i], dbTx)
	if err != nil {
		log.Error("error l1EventProcessors.Process: ", err)
		return err
	}
	return nil
}
