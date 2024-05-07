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

type stateInterface interface {
	BeginStateTransaction(ctx context.Context) (dbTxType, error)
	GetPreviousBlock(ctx context.Context, depth uint64, tx dbTxType) (*stateBlockType, error)
}

type CheckReorgManager struct {
	asyncL1BlockChecker l1_check_block.L1BlockCheckerIntegrator
	ctx                 context.Context
	etherMan            syncinterfaces.EthermanFullInterface
	state               stateInterface
	GenesisBlockNumber  uint64
}

func NewCheckReorgManager(ctx context.Context, etherMan syncinterfaces.EthermanFullInterface, state stateInterface, asyncL1BlockChecker l1_check_block.L1BlockCheckerIntegrator, genesisBlockNumber uint64) *CheckReorgManager {
	return &CheckReorgManager{
		asyncL1BlockChecker: asyncL1BlockChecker,
		ctx:                 ctx,
		etherMan:            etherMan,
		state:               state,
		GenesisBlockNumber:  genesisBlockNumber,
	}
}

func (s *CheckReorgManager) CheckReorg(latestBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, error) {
	if latestBlock == nil {
		err := fmt.Errorf("lastEthBlockSynced is nil calling checkReorgAndExecuteReset")
		log.Errorf("%s, it never have to happens", err.Error())
		return nil, err
	}
	block, errReturnedReorgFunction := s.NewCheckReorg(latestBlock, syncedBlock)
	if s.asyncL1BlockChecker != nil {
		return s.asyncL1BlockChecker.CheckReorgWrapper(s.ctx, block, errReturnedReorgFunction)
	}
	return block, errReturnedReorgFunction
}

/*
This function will check if there is a reorg.
As input param needs the last ethereum block synced. Retrieve the block info from the blockchain
to compare it with the stored info. If hash and hash parent matches, then no reorg is detected and return a nil.
If hash or hash parent don't match, reorg detected and the function will return the block until the sync process
must be reverted. Then, check the previous ethereum block synced, get block info from the blockchain and check
hash and has parent. This operation has to be done until a match is found.
*/

func (s *CheckReorgManager) NewCheckReorg(latestStoredBlock *stateBlockType, syncedBlock *etherman.Block) (*stateBlockType, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestStoredEthBlock := *latestStoredBlock
	reorgedBlock := *latestStoredBlock
	var depth uint64
	block := syncedBlock
	for {
		if block == nil {
			log.Infof("[checkReorg function] Checking Block %d in L1", reorgedBlock.BlockNumber)
			b, err := s.etherMan.EthBlockByNumber(s.ctx, reorgedBlock.BlockNumber)
			if err != nil {
				log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %v", reorgedBlock.BlockNumber, err)
				return nil, err
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
				return nil, err
			}
		} else {
			log.Infof("[checkReorg function] Using block %d from GetRollupInfoByBlockRange", block.BlockNumber)
		}
		log.Infof("[checkReorg function] BlockNumber: %d BlockHash got from L1 provider: %s", block.BlockNumber, block.BlockHash.String())
		log.Infof("[checkReorg function] reorgedBlockNumber: %d reorgedBlockHash already synced: %s", reorgedBlock.BlockNumber, reorgedBlock.BlockHash.String())

		// Compare hashes
		if (block.BlockHash != reorgedBlock.BlockHash || block.ParentHash != reorgedBlock.ParentHash) && reorgedBlock.BlockNumber > s.GenesisBlockNumber {
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
			dbTx, err := s.state.BeginStateTransaction(s.ctx)
			if err != nil {
				log.Errorf("error creating db transaction to get prevoius blocks")
				return nil, err
			}
			lb, err := s.state.GetPreviousBlock(s.ctx, depth, dbTx)
			errC := dbTx.Commit(s.ctx)
			if errC != nil {
				log.Errorf("error committing dbTx, err: %v", errC)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return nil, rollbackErr
				}
				log.Errorf("error committing dbTx, err: %v", errC)
				return nil, errC
			}
			if errors.Is(err, entities.ErrNotFound) {
				log.Warn("error checking reorg: previous block not found in db. Reorg reached the genesis block: %v.Genesis block can't be reorged, using genesis block as starting point. Error: %v", reorgedBlock, err)
				return &reorgedBlock, nil
			} else if err != nil {
				log.Error("error getting previousBlock from db. Error: ", err)
				return nil, err
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
		log.Info("Reorg detected in block: ", latestStoredEthBlock.BlockNumber, " last block OK: ", latestStoredBlock.BlockNumber)
		return latestStoredBlock, nil
	}
	log.Debugf("No reorg detected in block: %d. BlockHash: %s", latestStoredEthBlock.BlockNumber, latestStoredEthBlock.BlockHash.String())
	return nil, nil
}
