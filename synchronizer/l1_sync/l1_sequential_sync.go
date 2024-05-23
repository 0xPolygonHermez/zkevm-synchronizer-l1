package l1sync

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	syncommon "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_check_block"
	"github.com/ethereum/go-ethereum/common"
)

type EthermanInterface interface {
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	GetL1BlockByNumber(ctx context.Context, blockNumber uint64) (*etherman.Block, error)
}

type StateL1SeqInterface interface {
	BeginTransaction(ctx context.Context) (dbTxType, error)
	//GetPreviousBlock(ctx context.Context, depth uint64, fromBlockNumber *uint64, tx dbTxType) (*stateBlockType, error)
}
type BlockPoints struct {
	L1LastBlockToSync      uint64
	L1FinalizedBlockNumber uint64
}
type BlockPointsRetriever interface {
	GetL1BlockPoints(ctx context.Context) (BlockPoints, error)
}

type BlockRangeProcessor interface {
	ProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, finalizedBlockNumber uint64) error
}

type ReorgManager interface {
	// Returns new first valid block after reorg
	//MissingBlockOnResponseRollup(ctx context.Context, lastEthBlockSynced *stateBlockType) (*stateBlockType, error)
	// Returns first valid block if a reorg is detected
	CheckReorg(latestBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, uint64, error)
}

type BlockChecker interface {
	Step(ctx context.Context) error
}
type L1SequentialSync struct {
	blockPointsRetriever BlockPointsRetriever
	etherMan             EthermanInterface
	state                StateL1SeqInterface
	//ctx                  context.Context
	blockRangeProcessor BlockRangeProcessor
	reorgManager        ReorgManager
	cfg                 L1SequentialSyncConfig
	blockChecker        BlockChecker
}

type L1SequentialSyncConfig struct {
	SyncChunkSize                 uint64
	GenesisBlockNumber            uint64
	AllowEmptyBlocksAsCheckPoints bool
}

func NewL1SequentialSync(blockPointsRetriever BlockPointsRetriever,
	etherMan EthermanInterface,
	state StateL1SeqInterface,
	blockRangeProcessor BlockRangeProcessor,
	reorgManager ReorgManager,
	blockChecker BlockChecker,
	cfg L1SequentialSyncConfig) *L1SequentialSync {
	return &L1SequentialSync{
		blockPointsRetriever: blockPointsRetriever,
		etherMan:             etherMan,
		state:                state,
		blockRangeProcessor:  blockRangeProcessor,
		reorgManager:         reorgManager,
		blockChecker:         blockChecker,
		cfg:                  cfg,
	}
}

type BlockPointsRetrieverImplementation struct {
	syncBlockProtection l1_check_block.SafeL1BlockNumberFetcher
	finalizedBlock      l1_check_block.SafeL1BlockNumberFetcher
	l1Client            l1_check_block.L1Requester
}

func NewBlockPointsRetriever(syncBlockProtection, finalizedBlock l1_check_block.SafeL1BlockNumberFetcher, l1Client l1_check_block.L1Requester) *BlockPointsRetrieverImplementation {
	return &BlockPointsRetrieverImplementation{
		syncBlockProtection: syncBlockProtection,
		finalizedBlock:      finalizedBlock,
		l1Client:            l1Client,
	}
}

func (s *BlockPointsRetrieverImplementation) GetL1BlockPoints(ctx context.Context) (BlockPoints, error) {
	lastKnownBlock, err := s.syncBlockProtection.GetSafeBlockNumber(ctx, s.l1Client)
	if err != nil {
		log.Error("error getting header of the latest block in L1. Error: ", err)
		return BlockPoints{}, err
	}
	finalizedBlockNumber, err := s.finalizedBlock.GetSafeBlockNumber(ctx, s.l1Client)
	if err != nil {
		log.Errorf("error getting finalized block number in L1. Error: %v", err)
		return BlockPoints{}, err
	}
	log.Debugf("Getting block points: syncBlocksProtection: %s = %d, finalizedBlock: %s = %d",
		s.syncBlockProtection.Description(), lastKnownBlock,
		s.finalizedBlock.Description(), finalizedBlockNumber)

	return BlockPoints{
		L1LastBlockToSync:      lastKnownBlock,
		L1FinalizedBlockNumber: finalizedBlockNumber,
	}, nil
}

// It checks L1 block that still are the same in L1
func (s *L1SequentialSync) checkReorgsOnPreviousL1Blocks(ctx context.Context) error {
	if s.blockChecker != nil {
		err := s.blockChecker.Step(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// This function syncs the node from a specific block to the latest
// returns the last block synced and an error if any
// returns true if the sync is completed
func (s *L1SequentialSync) SyncBlocksSequential(ctx context.Context, lastEthBlockSynced *stateBlockType) (*stateBlockType, bool, error) {
	blockPoints, err := s.blockPointsRetriever.GetL1BlockPoints(ctx)
	if err != nil {
		return lastEthBlockSynced, false, err
	}
	if blockPoints.L1FinalizedBlockNumber > blockPoints.L1LastBlockToSync {
		log.Warnf("Finalized block number %d is greater than last block to sync %d", blockPoints.L1FinalizedBlockNumber, blockPoints.L1LastBlockToSync)
	}

	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber
	}
	blockRangeIterator := NewBlockRangeIterator(fromBlock, s.cfg.SyncChunkSize, blockPoints.L1LastBlockToSync)

	for {
		log.Debugf("Check that old blocks haven't changed...")
		err := s.checkReorgsOnPreviousL1Blocks(ctx)
		if err != nil {
			return lastEthBlockSynced, false, err
		}
		if blockRangeIterator == nil {
			log.Debugf("Nothing to do starting from %d to %d. Skipping...", fromBlock, blockPoints.L1LastBlockToSync)
			return lastEthBlockSynced, true, nil
		}
		blockRange := blockRangeIterator.GetRange(lastEthBlockSynced.HasEvents && !lastEthBlockSynced.Checked)
		log.Infof("Syncing %s", blockRangeIterator.String())
		synced := false
		lastEthBlockSynced, synced, err = s.iteration(ctx, blockRange, blockPoints.L1FinalizedBlockNumber, lastEthBlockSynced)
		if err != nil {
			return lastEthBlockSynced, false, err
		}
		if synced {
			return lastEthBlockSynced, true, nil
		}
		if blockRangeIterator.IsLastRange() {
			break
		}
		blockRangeIterator = blockRangeIterator.NextRange(lastEthBlockSynced.BlockNumber)
	}
	return lastEthBlockSynced, true, nil
}

func (s *L1SequentialSync) checkResponseGetRollupInfoByBlockRangeForOverlappedFirstBlock(blocks []etherman.Block, fromBlock uint64) error {
	if len(blocks) != 0 {
		initBlockReceived := &blocks[0]
		if initBlockReceived.BlockNumber != fromBlock {
			log.Errorf("Reorg detected in block %d while querying GetRollupInfoByBlockRange. Expected response first block %d and get %d",
				fromBlock, fromBlock, initBlockReceived.BlockNumber)
			err := syncommon.NewReorgError(fromBlock, fmt.Errorf("reorg detected in block %d while querying GetRollupInfoByBlockRange", fromBlock))
			return err
		}
	} else {
		// Reorg detected
		// The first block that we ask in GetRollupInfoByBlockRange contains events, but in that case returns an empty array
		// so for sure this block have changed
		log.Errorf("Reorg detected in block %d while querying GetRollupInfoByBlockRange. Expected reponse must include overlapped block", fromBlock)
		err := syncommon.NewReorgError(fromBlock, fmt.Errorf("reorg detected in block %d while querying GetRollupInfoByBlockRange. Expected minimum 1 block", fromBlock))
		return err
	}
	return nil
}

// firstBlockRequestIsOverlapped = means that fromBlock is one that we already have and have events
func (s *L1SequentialSync) iteration(ctx context.Context, blockRange BlockRange, finalizedBlockNumber uint64, lastEthBlockSynced *stateBlockType) (*stateBlockType, bool, error) {

	log.Infof("Getting rollup info range: %s finalizedBlock:%d", blockRange.String(), finalizedBlockNumber)
	// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
	// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
	// Name can be different in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
	// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
	// array index where this value is.
	blocks, order, err := s.retrieveDataFromL1AndValidate(ctx, blockRange)
	if err != nil {
		return lastEthBlockSynced, false, err
	}
	blocks, initBlockReceived := s.extractInitialBlock(blockRange, blocks)

	lastEthBlockSynced, err = s.checkReorgs(lastEthBlockSynced, initBlockReceived)
	if err != nil {
		return lastEthBlockSynced, false, err
	}

	err = s.blockRangeProcessor.ProcessBlockRange(ctx, blocks, order, finalizedBlockNumber)
	if err != nil {
		log.Errorf("error processing the block range. Err: %v", err)
		return lastEthBlockSynced, false, err
	}
	if len(blocks) > 0 {
		lastBlockOnSequence := &blocks[len(blocks)-1]
		isFinalized := entities.IsBlockFinalized(lastBlockOnSequence.BlockNumber, finalizedBlockNumber)
		lastEthBlockSynced = entities.NewL1BlockFromEthermanBlock(lastBlockOnSequence, isFinalized)
		for i := range blocks {
			log.Info("Position: ", i, ". New block. BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
		}
	} else {
		// Decide if create a empty block to move the FromBlock
		// We can ignore error because if we can't create empty block we can continue extending toBlock
		emptyBlock, _ := s.CreateBlockWithoutRollupInfoIfNeeded(ctx, blockRange, finalizedBlockNumber)
		if emptyBlock != nil {
			log.Infof("Creating empty block  BlockNumber: %d", emptyBlock.BlockNumber)
			lastEthBlockSynced = emptyBlock

		}
	}
	return lastEthBlockSynced, false, nil
}

func (s *L1SequentialSync) retrieveDataFromL1AndValidate(ctx context.Context, blockRange BlockRange) ([]etherman.Block, map[common.Hash][]etherman.Order, error) {
	toBlock := blockRange.ToBlock
	blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(ctx, blockRange.FromBlock, &toBlock)
	if err != nil {
		log.Errorf("error getting rollup info by block range.  Err: %v", err)
		return nil, nil, err
	}
	if blockRange.OverlappedFirstBlock {
		err = s.checkResponseGetRollupInfoByBlockRangeForOverlappedFirstBlock(blocks, blockRange.FromBlock)
		if err != nil {
			return nil, nil, err
		}
	}
	return blocks, order, nil
}

func (s *L1SequentialSync) extractInitialBlock(blockRange BlockRange, blocks []etherman.Block) ([]etherman.Block, *etherman.Block) {
	var initBlockReceived *etherman.Block
	if len(blocks) != 0 && blockRange.OverlappedFirstBlock {
		// The first block is overlapped, it have been processed we only want
		// it to check reorgs (compare that have not changed between the previous checkReorg and the call GetRollupInfoByBlockRange)
		initBlockReceived = &blocks[0]
		// First position of the array must be deleted
		blocks = removeBlockElement(blocks, 0)
	}
	return blocks, initBlockReceived
}

// checkReorgs returns first good block and error if something is detected
func (s *L1SequentialSync) checkReorgs(lastEthBlockSynced *stateBlockType, initBlockReceived *etherman.Block) (*stateBlockType, error) {

	if !lastEthBlockSynced.Checked || initBlockReceived != nil {
		// Check reorg again to be sure that the chain has not changed between the previous checkReorg and the call GetRollupInfoByBlockRange
		log.Debugf("Checking reorgs between lastEthBlockSynced =%d and initBlockReceived: %v", lastEthBlockSynced.BlockNumber, initBlockReceived)
		block, lastBadBlockNumber, err := s.reorgManager.CheckReorg(lastEthBlockSynced, initBlockReceived)
		log.Debugf("Checking reorgs between lastEthBlockSynced =%d [AFTER]", lastEthBlockSynced.BlockNumber)
		if err != nil && !errors.Is(err, ErrReorgAllBlocksOnDBAreBad) {
			log.Errorf("error checking reorgs. Retrying... Err: %v", err)
			return lastEthBlockSynced, fmt.Errorf("error checking reorgs. Err:%w", err)
		}
		if block != nil {
			// I use lastBadBlockNumber to support the case that all blocks are bad on DB
			err := syncommon.NewReorgError(lastBadBlockNumber, fmt.Errorf("reorg detected. First valid block is %d, lastBadBlock is %d by CheckReorg func", block.BlockNumber, lastBadBlockNumber))
			return block, err
		}
		if errors.Is(err, ErrReorgAllBlocksOnDBAreBad) {
			log.Warn("Reorg detected. Affect all l1block on DB")
			err := syncommon.NewReorgError(lastBadBlockNumber, fmt.Errorf("reorg detected. Affect all l1block on DB, lastBadBlock is %d by CheckReorg func", lastBadBlockNumber))
			return nil, err
		}

	} else {
		log.Debugf("Skipping reorg check because lastEthBlockSynced %d is checked", lastEthBlockSynced.BlockNumber)

	}
	return lastEthBlockSynced, nil
}

// CreateBlockWithoutRollupInfoIfNeeded creates a block without rollup if
// the condition is that the fromBlock + SyncChunkSize < finalizedBlockNumber
// because means that the empty block that we can't check the reorg is safe create it
func (s *L1SequentialSync) CreateBlockWithoutRollupInfoIfNeeded(ctx context.Context, blockRange BlockRange, finalizedBlockNumber uint64) (*stateBlockType, error) {
	if !s.cfg.AllowEmptyBlocksAsCheckPoints {
		return nil, nil
	}
	proposedBlockNumber := blockRange.FromBlock + s.cfg.SyncChunkSize

	if !blockRange.InsideRange(proposedBlockNumber) {
		return nil, nil
	}
	if entities.IsBlockFinalized(proposedBlockNumber, finalizedBlockNumber) {
		// Create a block without rollup info
		emptyBlock, err := s.etherMan.GetL1BlockByNumber(ctx, proposedBlockNumber)
		if err != nil || emptyBlock == nil {
			log.Warnf("error getting block %d from the blockchain. Error: %v", proposedBlockNumber, err)
			return nil, err
		}
		err = s.blockRangeProcessor.ProcessBlockRange(ctx, []etherman.Block{*emptyBlock}, nil, finalizedBlockNumber)
		if err != nil {
			log.Warnf("error processing the block range. Err: %v", err)
			return nil, err
		}
		return entities.NewL1BlockFromEthermanBlock(emptyBlock, true), nil
	}
	return nil, nil
}

func removeBlockElement(slice []etherman.Block, s int) []etherman.Block {
	ret := make([]etherman.Block, 0)
	ret = append(ret, slice[:s]...)
	return append(ret, slice[s+1:]...)
}
