package l1sync

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_check_block"
)

type EthermanInterface interface {
}

type stateL1SeqInterface interface {
	BeginStateTransaction(ctx context.Context) (dbTxType, error)
	GetPreviousBlock(ctx context.Context, depth uint64, tx dbTxType) (*stateBlockType, error)
}

type L1SequentialSync struct {
	syncBlockProtection l1_check_block.SafeL1BlockNumberFetcher
	l1Client            l1_check_block.L1Requester
	etherMan            syncinterfaces.EthermanFullInterface
	state               stateL1SeqInterface
	SyncChunkSize       uint64
	ctx                 context.Context
	blockRangeProcessor syncinterfaces.BlockRangeProcessor
}

func (s *L1SequentialSync) resetState(blockNumber uint64) error {
	// Reset the state to the previous block
	return nil
}

func (s *L1SequentialSync) checkReorg(latestBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, error) {
	// Check if there is a reorg
	return nil, nil
}

// This function syncs the node from a specific block to the latest
func (s *L1SequentialSync) SyncBlocksSequential(ctx context.Context, lastEthBlockSynced *stateBlockType) (*stateBlockType, error) {
	// Call the blockchain to retrieve data
	lastKnownBlock, err := s.syncBlockProtection.GetSafeBlockNumber(ctx, s.l1Client)
	if err != nil {
		log.Error("error getting header of the latest block in L1. Error: ", err)
		return lastEthBlockSynced, err
	}

	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber
	}
	toBlock := fromBlock + s.SyncChunkSize

	for {
		if toBlock > lastKnownBlock {
			log.Debug("Setting toBlock to the lastKnownBlock: ", lastKnownBlock)
			toBlock = lastKnownBlock
		}
		if fromBlock > toBlock {
			log.Debug("FromBlock is higher than toBlock. Skipping...")
			return lastEthBlockSynced, nil
		}
		log.Infof("Syncing block %d of %d", fromBlock, lastKnownBlock)
		log.Infof("Getting rollup info from block %d to block %d", fromBlock, toBlock)
		// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
		// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
		// Name can be different in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
		// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
		// array index where this value is.

		blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(ctx, fromBlock, &toBlock)

		if err != nil {
			return lastEthBlockSynced, err
		}

		var initBlockReceived *etherman.Block
		if len(blocks) != 0 {
			initBlockReceived = &blocks[0]
			// First position of the array must be deleted
			blocks = removeBlockElement(blocks, 0)
		} else {
			// Reorg detected
			log.Infof("Reorg detected in block %d while querying GetRollupInfoByBlockRange. Rolling back to at least the previous block", fromBlock)
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
		// Check reorg again to be sure that the chain has not changed between the previous checkReorg and the call GetRollupInfoByBlockRange
		block, err := s.checkReorg(lastEthBlockSynced, initBlockReceived)
		if err != nil {
			log.Errorf("error checking reorgs. Retrying... Err: %v", err)
			return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
		}
		if block != nil {
			err = s.resetState(block.BlockNumber)
			if err != nil {
				log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
				return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
			}
			return block, nil
		}

		err = s.blockRangeProcessor.ProcessBlockRange(s.ctx, blocks, order, lastKnownBlock)

		if err != nil {
			return lastEthBlockSynced, err
		}
		if len(blocks) > 0 {
			lastEthBlockSynced = &stateBlockType{
				BlockNumber: blocks[len(blocks)-1].BlockNumber,
				BlockHash:   blocks[len(blocks)-1].BlockHash,
				ParentHash:  blocks[len(blocks)-1].ParentHash,
				ReceivedAt:  blocks[len(blocks)-1].ReceivedAt,
			}
			for i := range blocks {
				log.Info("Position: ", i, ". New block. BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
			}
		}

		//if lastKnownBlock.Cmp(new(big.Int).SetUint64(toBlock)) < 1 {
		if lastKnownBlock < toBlock {
			//  	waitDuration = s.cfg.SyncInterval.Duration
			break
		}

		fromBlock = lastEthBlockSynced.BlockNumber
		toBlock = toBlock + s.SyncChunkSize
	}

	return lastEthBlockSynced, nil
}

func removeBlockElement(slice []etherman.Block, s int) []etherman.Block {
	ret := make([]etherman.Block, 0)
	ret = append(ret, slice[:s]...)
	return append(ret, slice[s+1:]...)
}
