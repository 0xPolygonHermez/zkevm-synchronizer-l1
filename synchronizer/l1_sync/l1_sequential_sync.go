package l1sync

import (
	"context"
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
}

type StateL1SeqInterface interface {
	BeginTransaction(ctx context.Context) (dbTxType, error)
	GetPreviousBlock(ctx context.Context, depth uint64, fromBlockNumber *uint64, tx dbTxType) (*stateBlockType, error)
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

type L1SequentialSync struct {
	blockPointsRetriever BlockPointsRetriever
	etherMan             EthermanInterface
	state                StateL1SeqInterface
	//ctx                  context.Context
	blockRangeProcessor BlockRangeProcessor
	reorgManager        ReorgManager
	cfg                 L1SequentialSyncConfig
}

type L1SequentialSyncConfig struct {
	SyncChunkSize      uint64
	GenesisBlockNumber uint64
}

func NewL1SequentialSync(blockPointsRetriever BlockPointsRetriever,
	etherMan EthermanInterface,
	state StateL1SeqInterface,
	blockRangeProcessor BlockRangeProcessor,
	reorgManager ReorgManager,
	cfg L1SequentialSyncConfig) *L1SequentialSync {
	return &L1SequentialSync{
		blockPointsRetriever: blockPointsRetriever,
		etherMan:             etherMan,
		state:                state,
		blockRangeProcessor:  blockRangeProcessor,
		reorgManager:         reorgManager,
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
	return BlockPoints{
		L1LastBlockToSync:      lastKnownBlock,
		L1FinalizedBlockNumber: finalizedBlockNumber,
	}, nil
}

type IterateBlockRange struct {
	FromBlock, ToBlock uint64
	SyncChunkSize      uint64
	MaximumBlock       uint64
}

func NewIterateBlockRange(fromBlock, syncChunkSize uint64, maximumBlock uint64) *IterateBlockRange {
	res := &IterateBlockRange{
		FromBlock:     fromBlock,
		ToBlock:       fromBlock,
		SyncChunkSize: syncChunkSize,
		MaximumBlock:  maximumBlock,
	}
	res = res.NextRange(fromBlock)
	return res
}
func (i *IterateBlockRange) IsLastRange() bool {
	return i.FromBlock >= i.MaximumBlock
}

func (i *IterateBlockRange) NextRange(fromBlock uint64) *IterateBlockRange {
	// The FromBlock is the new block (can be the previous one if no blocks found in the range)
	i.FromBlock = fromBlock
	// Extend toBlock by sync chunk size
	i.ToBlock = i.ToBlock + i.SyncChunkSize

	if i.ToBlock > i.MaximumBlock {
		i.ToBlock = i.MaximumBlock
	}
	if i.FromBlock > i.ToBlock {
		return nil
	}
	return i
}
func (i *IterateBlockRange) String() string {
	return fmt.Sprintf("FromBlock: %d, ToBlock: %d, MaximumBlock: %d", i.FromBlock, i.ToBlock, i.MaximumBlock)
}

// This function syncs the node from a specific block to the latest
// returns the last block synced and an error if any
// returns true if the sync is completed
func (s *L1SequentialSync) SyncBlocksSequential(ctx context.Context, lastEthBlockSynced *stateBlockType) (*stateBlockType, bool, error) {
	// Call the blockchain to retrieve data
	blockPoints, err := s.blockPointsRetriever.GetL1BlockPoints(ctx)
	if err != nil {
		return lastEthBlockSynced, false, err
	}

	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber
	}
	blockRange := NewIterateBlockRange(fromBlock, s.cfg.SyncChunkSize, blockPoints.L1LastBlockToSync)

	for {
		if blockRange == nil {
			log.Debugf("Nothing to do starting from %d to %d. Skipping...", fromBlock, blockPoints.L1LastBlockToSync)
			return lastEthBlockSynced, true, nil
		}

		log.Infof("Syncing %s", blockRange.String())
		lastEthBlockSynced, synced, err := s.iteration(ctx, blockRange.FromBlock, blockRange.ToBlock, blockPoints.L1FinalizedBlockNumber, lastEthBlockSynced)
		if err != nil {
			return lastEthBlockSynced, false, err
		}
		if synced {
			return lastEthBlockSynced, true, nil
		}
		if blockRange.IsLastRange() {
			break
		}
		blockRange = blockRange.NextRange(lastEthBlockSynced.BlockNumber)
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
func (s *L1SequentialSync) iteration(ctx context.Context, fromBlock, toBlock, finalizedBlockNumber uint64, lastEthBlockSynced *stateBlockType) (*stateBlockType, bool, error) {

	log.Infof("Getting rollup info from block %d to block %d", fromBlock, toBlock)
	// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
	// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
	// Name can be different in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
	// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
	// array index where this value is.
	firstBlockRequestIsOverlapped := lastEthBlockSynced.HasEvents
	if !firstBlockRequestIsOverlapped {
		log.Infof("Request is not overlapped, so is not going to detect reorgs. FromBlock: %d toBlock:%d", fromBlock)
	}
	blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(ctx, fromBlock, &toBlock)
	if err != nil {
		log.Errorf("error getting rollup info by block range.  Err: %v", err)
		return lastEthBlockSynced, false, err
	}
	if firstBlockRequestIsOverlapped {
		err = s.checkResponseGetRollupInfoByBlockRangeForOverlappedFirstBlock(blocks, fromBlock)
		if err != nil {
			return lastEthBlockSynced, false, err
		}
	}

	var initBlockReceived *etherman.Block
	if len(blocks) != 0 && firstBlockRequestIsOverlapped {
		initBlockReceived = &blocks[0]
		// First position of the array must be deleted
		blocks = removeBlockElement(blocks, 0)
	}
	// Check reorg again to be sure that the chain has not changed between the previous checkReorg and the call GetRollupInfoByBlockRange
	block, lastBadBlockNumber, err := s.reorgManager.CheckReorg(lastEthBlockSynced, initBlockReceived)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, false, fmt.Errorf("error checking reorgs. Err:%w", err)
	}
	if block != nil {
		// In fact block.BlockNumber is the first ok block, so  add 1 to be the first block wrong
		// maybe doesnt exists
		err := syncommon.NewReorgError(lastBadBlockNumber, fmt.Errorf("reorg detected. First valid block is %d, lastBadBlock is %d by CheckReorg func", block.BlockNumber, lastBadBlockNumber))
		return block, false, err
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
	}
	return lastEthBlockSynced, false, nil
}

/*
	func (s *L1SequentialSync) reorgDetected(ctx context.Context, lastEthBlockSynced *stateBlockType) (*stateBlockType, error) {
		prevBlock, err := s.state.GetPreviousBlock(ctx, 1, nil)
		if errors.Is(err, entities.ErrNotFound) {
			log.Warn("error checking reorg: previous block not found in db: ", err)
			prevBlock = &stateBlockType{}
		} else if err != nil {
			log.Error("error getting previousBlock from db. Error: ", err)
			return lastEthBlockSynced, err
		}
		blockReorged, err := s.checkReorg(prevBlock, nil)
		if err != nil {
			log.Error("error checking reorgs in previous blocks. Error: ", err)
			return lastEthBlockSynced, err
		}
		if blockReorged == nil {
			blockReorged = prevBlock
		}
		err = s.resetState(blockReorged.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return blockReorged, nil
	}
*/
func removeBlockElement(slice []etherman.Block, s int) []etherman.Block {
	ret := make([]etherman.Block, 0)
	ret = append(ret, slice[:s]...)
	return append(ret, slice[s+1:]...)
}
