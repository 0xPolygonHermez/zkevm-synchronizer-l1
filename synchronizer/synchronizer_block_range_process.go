package synchronizer

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1event_orders"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateTxProvider = syncinterfaces.StateTxProvider

// BlockRangeProcess is the struct that process the block range that implements syncinterfaces.BlockRangeProcessor
type BlockRangeProcess struct {
	storage           syncinterfaces.StorageBlockReaderInterface
	stateForkId       syncinterfaces.StateForkIdQuerier
	stateTxProvider   stateTxProvider
	l1EventProcessors syncinterfaces.L1EventProcessorManager
}

// NewBlockRangeProcessLegacy creates a new BlockRangeProcess
func NewBlockRangeProcessLegacy(
	state syncinterfaces.StorageBlockReaderInterface,
	stateForkId syncinterfaces.StateForkIdQuerier,
	stateTxProvider stateTxProvider,
	l1EventProcessors syncinterfaces.L1EventProcessorManager,
) *BlockRangeProcess {
	return &BlockRangeProcess{
		storage:           state,
		stateForkId:       stateForkId,
		stateTxProvider:   stateTxProvider,
		l1EventProcessors: l1EventProcessors,
	}
}

// ProcessBlockRangeSingleDbTx process the L1 events and stores the information in the db reusing same DbTx
func (s *BlockRangeProcess) ProcessBlockRangeSingleDbTx(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order,
	finalizedBlockNumber uint64, storeBlocks syncinterfaces.ProcessBlockRangeL1BlocksMode, dbTx stateTxType) error {
	txProvider := NewReuseStateTxProvider(s.stateTxProvider, dbTx)
	return s.internalProcessBlockRange(ctx, blocks, order, finalizedBlockNumber, storeBlocks, txProvider)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) ProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, finalizedBlockNumber uint64) error {
	return s.internalProcessBlockRange(ctx, blocks, order, finalizedBlockNumber, syncinterfaces.StoreL1Blocks, s.stateTxProvider)
}

func isBlockFinalized(blockNumber uint64, finalizedBlockNumber uint64) bool {
	return entities.IsBlockFinalized(blockNumber, finalizedBlockNumber)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) addBlock(ctx context.Context, block *etherman.Block, isFinalized bool, dbTx stateTxType) error {
	b := entities.NewL1BlockFromEthermanBlock(block, isFinalized)
	// Add block information
	return s.storage.AddBlock(ctx, b, dbTx)
}

func (s *BlockRangeProcess) rollback(ctx context.Context, err error, dbTx stateTxType) error {
	// Rollback db transaction
	rollbackErr := dbTx.Rollback(ctx)
	if rollbackErr != nil {
		if !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			log.Errorf("error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
			return rollbackErr
		} else {
			log.Debugf("error rolling back state because is already closed. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
			return err
		}
	}
	return err
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) internalProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order,
	finalizedBlockNumber uint64,
	storeBlocksMode syncinterfaces.ProcessBlockRangeL1BlocksMode,
	txProvider stateTxProvider) error {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		dbTx, err := txProvider.BeginTransaction(ctx)
		if err != nil {
			errExt := fmt.Errorf("error beginStateTransaction. BlockNumber: %d, Error: %w", blocks[i].BlockNumber, err)
			return errExt
		}
		// Process event received from l1
		err = s.processBlock(ctx, blocks, i, dbTx, order, storeBlocksMode, finalizedBlockNumber)
		if err != nil {
			return s.rollback(ctx, err, dbTx)
		}

		err = dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("error committing state. BlockNumber: %d, Error: %v", blocks[i].BlockNumber, err)
			return err
		}

	}
	return nil
}

func (s *BlockRangeProcess) processBlock(ctx context.Context, blocks []etherman.Block, blockIndex int, dbTx stateTxType, order map[common.Hash][]etherman.Order, storeBlock syncinterfaces.ProcessBlockRangeL1BlocksMode, finalizedBlockNumber uint64) error {
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

	return nil
}

func (s *BlockRangeProcess) processElement(ctx context.Context, element etherman.Order, blocks []etherman.Block, i int, dbTx stateTxType) error {
	batchSequence := l1event_orders.GetSequenceFromL1EventOrder(element.Name, &blocks[i], element.Pos)
	forkId := entities.FORKID_ZERO
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
