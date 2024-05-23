package l1sync

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	// This is a special case that means that all blocks are bad on DB, is a bit weird
	// for that reason returns a specific error to allow the caller to take action
	ErrReorgAllBlocksOnDBAreBad = fmt.Errorf("reorg reached the genesis block")
)

type StateReorgInterface interface {
	BeginTransaction(ctx context.Context) (dbTxType, error)
	GetPreviousBlock(ctx context.Context, depth uint64, tx dbTxType) (*stateBlockType, error)
}

type EthermanReorgManager interface {
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*ethTypes.Block, error)
}

type CheckReorgManager struct {
	ctx                context.Context
	etherMan           EthermanReorgManager
	state              StateReorgInterface
	GenesisBlockNumber uint64
}

func NewCheckReorgManager(ctx context.Context, etherMan EthermanReorgManager, state StateReorgInterface, genesisBlockNumber uint64) *CheckReorgManager {
	return &CheckReorgManager{
		ctx:                ctx,
		etherMan:           etherMan,
		state:              state,
		GenesisBlockNumber: genesisBlockNumber,
	}
}

// CheckReorg checks consistency of blocks
// Returns:
// - first block ok
// - last bad block number
// - error
func (s *CheckReorgManager) CheckReorg(latestBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, uint64, error) {
	if latestBlock == nil {
		err := fmt.Errorf("lastEthBlockSynced is nil calling checkReorgAndExecuteReset")
		log.Errorf("%s, it never have to happens", err.Error())
		return nil, 0, err
	}
	blockOk, badBlockNumber, errReturnedReorgFunction := s.NewCheckReorg(latestBlock, syncedBlock)
	return blockOk, badBlockNumber, errReturnedReorgFunction
}

/*
This function will check if there is a reorg.
As input param needs the last ethereum block synced. Retrieve the block info from the blockchain
to compare it with the stored info. If hash and hash parent matches, then no reorg is detected and return a nil.
If hash or hash parent don't match, reorg detected and the function will return the block until the sync process
must be reverted. Then, check the previous ethereum block synced, get block info from the blockchain and check
hash and has parent. This operation has to be done until a match is found.
*/

// Returns:
// - first block ok
// - last bad block number
// - error
func (s *CheckReorgManager) NewCheckReorg(latestStoredBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, uint64, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestStoredEthBlock := *latestStoredBlock
	reorgedBlock := *latestStoredBlock
	var depth uint64
	block := syncedBlock
	lastBadBlockNumber := uint64(0)
	for {
		if block == nil {
			log.Infof("[checkReorg function] Checking Block %d in L1", reorgedBlock.BlockNumber)
			b, err := s.etherMan.EthBlockByNumber(s.ctx, reorgedBlock.BlockNumber)
			if err != nil {
				log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %v", reorgedBlock.BlockNumber, err)
				return nil, 0, err
			}
			block = &etherman.Block{
				BlockNumber: b.Number().Uint64(),
				BlockHash:   b.Hash(),
				ParentHash:  b.ParentHash(),
			}
			if block.BlockNumber != reorgedBlock.BlockNumber {
				err := fmt.Errorf("wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
					reorgedBlock.BlockNumber, block.BlockNumber)
				log.Error("error: ", err)
				return nil, 0, err
			}
		} else {
			log.Infof("[checkReorg function] Using block %d from GetRollupInfoByBlockRange", block.BlockNumber)
		}
		log.Infof("[checkReorg function] BlockNumber: %d BlockHash got from L1 provider: %s", block.BlockNumber, block.BlockHash.String())
		log.Infof("[checkReorg function] reorgedBlockNumber: %d reorgedBlockHash already synced: %s", reorgedBlock.BlockNumber, reorgedBlock.BlockHash.String())

		// Compare hashes
		if block.BlockHash != reorgedBlock.BlockHash || block.ParentHash != reorgedBlock.ParentHash {
			log.Infof("checkReorg: Bad block %d hashOk %t parentHashOk %t", reorgedBlock.BlockNumber, block.BlockHash == reorgedBlock.BlockHash, block.ParentHash == reorgedBlock.ParentHash)
			log.Debug("[checkReorg function] => latestBlockNumber: ", reorgedBlock.BlockNumber)
			log.Debug("[checkReorg function] => latestBlockHash: ", reorgedBlock.BlockHash)
			log.Debug("[checkReorg function] => latestBlockHashParent: ", reorgedBlock.ParentHash)
			log.Debug("[checkReorg function] => BlockNumber: ", reorgedBlock.BlockNumber, block.BlockNumber)
			log.Debug("[checkReorg function] => BlockHash: ", block.BlockHash)
			log.Debug("[checkReorg function] => BlockHashParent: ", block.ParentHash)
			depth++
			log.Debug("REORG: Looking for the latest correct ethereum block. Depth: ", depth)
			// Reorg detected. Getting previous block
			lastBadBlockNumber = reorgedBlock.BlockNumber
			lb, err := s.state.GetPreviousBlock(s.ctx, depth, nil)

			if errors.Is(err, entities.ErrNotFound) {
				log.Warnf("error checking reorg: previous block not found in db. Reorg reached las block on DB (usually the genesis block): %v.That is very unusual. It returns an error", reorgedBlock)
				return nil, reorgedBlock.BlockNumber, ErrReorgAllBlocksOnDBAreBad
			} else if err != nil {
				log.Error("error getting previousBlock from db. Error: ", err)
				return nil, 0, err
			}
			reorgedBlock = *lb
		} else {
			log.Debugf("checkReorg: Block %d hashOk %t parentHashOk %t", reorgedBlock.BlockNumber, block.BlockHash == reorgedBlock.BlockHash, block.ParentHash == reorgedBlock.ParentHash)
			break
		}
		// This forces to get the block from L1 in the next iteration of the loop
		block = nil
	}
	if latestStoredEthBlock.BlockHash != reorgedBlock.BlockHash {
		latestStoredBlock = &reorgedBlock
		log.Info("Reorg detected in block: ", latestStoredEthBlock.BlockNumber, " last block OK: ", latestStoredBlock.BlockNumber, " first bad Block: ", lastBadBlockNumber)
		return latestStoredBlock, lastBadBlockNumber, nil
	}
	log.Debugf("No reorg detected in block: %d. BlockHash: %s", latestStoredEthBlock.BlockNumber, latestStoredEthBlock.BlockHash.String())
	return nil, 0, nil
}
