package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// SynchronizerImpl connects L1 and L2
type SynchronizerImpl struct {
	etherMan EthermanInterface
	// TODO: remove

	storage        StorageInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            syncconfig.Config
	networkID      uint
	// TODO: remove
	//chExitRootEvent chan *etherman.GlobalExitRoot
	// TODO: remove
	//chSynced chan uint
	// TODO: remove
	synced bool

	l1EventProcessors   *processor_manager.L1EventProcessors
	blockRangeProcessor syncinterfaces.BlockRangeProcessor
	l1InfoTreeManager   *state.L1InfoTreeState
	ForkIdState         *state.ForkIdState
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizerImpl(
	ctx context.Context,
	storage StorageInterface,
	ethMan EthermanInterface,
	cfg syncconfig.Config) (*SynchronizerImpl, error) {
	ctx, cancel := context.WithCancel(ctx)
	// TODO: ??? GetNetworkID
	networkID := uint(0)

	sync := SynchronizerImpl{
		storage:        storage,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: cfg.GenesisBlockNumber,
		cfg:            cfg,
		networkID:      networkID,
		//chExitRootEvent:  chExitRootEvent,
		//chSynced: chSynced,
		//zkEVMClient:      zkEVMClient,
		//l1RollupExitRoot: ger.ExitRoots[1],
	}
	sync.l1InfoTreeManager = state.NewL1InfoTreeManager(storage)
	sync.ForkIdState = state.NewForkIdState(storage)
	builder := processor_manager.NewL1EventProcessorsBuilder()
	builder.Register(etrog.NewProcessorL1InfoTreeUpdate(sync.l1InfoTreeManager))
	builder.Register(etrog.NewProcessorL1SequenceBatches(storage))
	builder.Register(incaberry.NewProcessorForkId(sync.ForkIdState))
	sync.l1EventProcessors = builder.Build()
	sync.blockRangeProcessor = NewBlockRangeProcessLegacy(storage, nil, sync.l1EventProcessors, nil)
	if cfg.GenesisBlockNumber == 0 {
		firstBlock, err := ethMan.GetL1BlockUpgradeLxLy(ctx, nil)
		if err != nil {
			log.Errorf("Error getting the first block from the blockchain. Error: %v", err)
			return nil, err
		}
		log.Infof("First block from the blockchain: %d (ETROG)", firstBlock)
		sync.genBlockNumber = firstBlock
	}
	return &sync, nil

}

var waitDuration = time.Duration(0)

// IsSynced returns true if the synchronizer is synced or false if it's not
func (s *SynchronizerImpl) IsSynced() bool {
	return s.synced
}

func (s *SynchronizerImpl) GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]L1InfoTreeLeaf, error) {
	leaves, err := s.l1InfoTreeManager.GetLeafsByL1InfoRoot(ctx, l1InfoRoot, dbTx)
	if err != nil {
		log.Error("error getting leaves by L1InfoRoot. Error: ", err)
		return nil, err
	}
	var res []L1InfoTreeLeaf
	for _, leaf := range leaves {
		tmp := L1InfoTreeLeaf(leaf)
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SynchronizerImpl) GetL1InfoRootPerIndex(ctx context.Context, L1InfoTreeIndex uint32) (common.Hash, error) {
	root, err := s.l1InfoTreeManager.GetL1InfoRootPerLeafIndex(ctx, L1InfoTreeIndex, nil)
	if errors.Is(err, state.ErrNotFound) {
		return common.Hash{}, ErrNotFound
	}
	return root, err
}

func (s *SynchronizerImpl) GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32) (map[uint32]L1InfoTreeLeaf, error) {
	leaves, err := s.l1InfoTreeManager.GetL1InfoTreeLeaves(ctx, indexLeaves, nil)
	if err != nil {
		return nil, err
	}
	// Convert type state.L1InfoTreeLeaf to type L1InfoTreeLeaf
	returnLeaves := make(map[uint32]L1InfoTreeLeaf)
	for _, idx := range indexLeaves {
		returnLeaves[idx] = L1InfoTreeLeaf(leaves[idx])
	}
	return returnLeaves, nil
}

func (s *SynchronizerImpl) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64) (*SequencedBatches, error) {
	sequence, err := s.storage.GetSequenceByBatchNumber(ctx, batchNumber, nil)
	if sequence == nil {
		return nil, err
	}
	res := SequencedBatches(*sequence)
	return &res, err
}

func convertStorageBlock(block *pgstorage.L1Block) *L1Block {
	return &L1Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
	}
}

func (s *SynchronizerImpl) getLastBlockFromStorage(ctx context.Context, dbTx pgx.Tx) (*L1Block, error) {
	lastBlockSynced, err := s.storage.GetLastBlock(s.ctx, dbTx)
	if err != nil {
		return nil, err
	}
	if lastBlockSynced == nil {
		return nil, nil
	}
	return convertStorageBlock(lastBlockSynced), nil
}

func (s *SynchronizerImpl) getPreviousBlockFromStorage(ctx context.Context, offset uint64, dbTx pgx.Tx) (*L1Block, error) {
	dbBlock, err := s.storage.GetPreviousBlock(s.ctx, offset, dbTx)
	if err != nil {
		return nil, err
	}
	if dbBlock == nil {
		return nil, nil
	}
	return convertStorageBlock(dbBlock), nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *SynchronizerImpl) Sync(returnOnSync bool) error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Infof("NetworkID: %d, Synchronization started", s.networkID)
	s.synced = false

	forks, _ := s.etherMan.GetForks(s.ctx, s.genBlockNumber, s.genBlockNumber-10)
	log.Info(forks)
	lastBlockSynced, err := s.getLastBlockFromStorage(s.ctx, nil)
	if err != nil {
		if errors.Is(err, pgstorage.ErrNotFound) {
			//log.Infof("networkID: %d, error getting the latest ethereum block. No data stored. Setting genesis block. Error: %v", s.networkID, err)
			lastBlockSynced = &L1Block{
				BlockNumber: max(0, s.genBlockNumber-1),
			}
			log.Infof("networkID: %d, error getting the latest block. No data stored. Using starting block: %d ",
				s.networkID, lastBlockSynced.BlockNumber)
		} else {
			log.Fatalf("networkID: %d, unexpected error getting the latest block. Error: %s", s.networkID, err.Error())
		}
	} else {
		log.Infof("networkID: %d, continuing from the last block stored on DB. lastBlockSynced: %+v", s.networkID, lastBlockSynced)

	}
	log.Infof("NetworkID: %d, initial lastBlockSynced: %+v", s.networkID, lastBlockSynced)
	for {
		select {
		case <-s.ctx.Done():
			log.Debugf("NetworkID: %d, synchronizer ctx done", s.networkID)
			return nil
		case <-time.After(waitDuration):
			log.Debugf("NetworkID: %d, syncing...", s.networkID)
			//Sync L1Blocks

			if lastBlockSynced, err = s.syncBlocks(lastBlockSynced); err != nil {
				log.Warnf("networkID: %d, error syncing blocks: %v", s.networkID, err)
				lastBlockSynced, err = s.getLastBlockFromStorage(s.ctx, nil)
				if err != nil {
					log.Fatalf("networkID: %d, error getting lastBlockSynced to resume the synchronization... Error: ", s.networkID, err)
				}
				if s.ctx.Err() != nil {
					continue
				}
			}
			if !s.synced {
				// Check latest Block
				header, err := s.etherMan.HeaderByNumber(s.ctx, nil)
				if err != nil {
					log.Warnf("networkID: %d, error getting latest block from. Error: %s", s.networkID, err.Error())
					continue
				}
				lastKnownBlock := header.Number.Uint64()
				if lastBlockSynced.BlockNumber == lastKnownBlock && !s.synced {
					log.Infof("NetworkID %d Synced!", s.networkID)
					waitDuration = s.cfg.SyncInterval.Duration
					s.synced = true
					//s.chSynced <- s.networkID
					if returnOnSync {
						log.Infof("NetworkID: %d, Synchronization finished, finishing", s.networkID)
						return nil
					}
				}
				if lastBlockSynced.BlockNumber > lastKnownBlock {
					if s.networkID == 0 {
						log.Fatalf("networkID: %d, error: latest Synced BlockNumber (%d) is higher than the latest Proposed block (%d) in the network", s.networkID, lastBlockSynced.BlockNumber, lastKnownBlock)
					} else {
						log.Errorf("networkID: %d, error: latest Synced BlockNumber (%d) is higher than the latest Proposed block (%d) in the network", s.networkID, lastBlockSynced.BlockNumber, lastKnownBlock)
						err = s.resetState(lastKnownBlock)
						if err != nil {
							log.Errorf("networkID: %d, error resetting the state to a previous block. Error: %v", s.networkID, err)
							continue
						}
					}
				}
			}
		}
	}
}

// Stop function stops the synchronizer
func (s *SynchronizerImpl) Stop() {
	s.cancelCtx()
}

// This function syncs the node from a specific block to the latest
func (s *SynchronizerImpl) syncBlocks(lastBlockSynced *L1Block) (*L1Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastBlockSynced)
	if err != nil {
		log.Errorf("networkID: %d, error checking reorgs. Retrying... Err: %s", s.networkID, err.Error())
		return lastBlockSynced, fmt.Errorf("networkID: %d, error checking reorgs", s.networkID)
	}
	if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("networkID: %d, error resetting the state to a previous block. Retrying... Error: %s", s.networkID, err.Error())
			return lastBlockSynced, fmt.Errorf("networkID: %d, error resetting the state to a previous block", s.networkID)
		}
		return block, nil
	}
	log.Debugf("NetworkID: %d, after checkReorg: no reorg detected", s.networkID)
	// Call the blockchain to retrieve data
	header, err := s.etherMan.HeaderByNumber(s.ctx, nil)
	if err != nil {
		return lastBlockSynced, err
	}
	lastKnownBlock := header.Number

	var fromBlock uint64
	if lastBlockSynced.BlockNumber > 0 {
		fromBlock = lastBlockSynced.BlockNumber + 1
	}

	for {
		toBlock := fromBlock + s.cfg.SyncChunkSize

		log.Debugf("NetworkID: %d, Getting info from block %d to block %d", s.networkID, fromBlock, toBlock)
		// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
		// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is read.
		// Name can be different in the order struct. This name is an identifier to check if the next info that must be stored in the db.
		// The value pos (position) tells what is the array index where this value is.
		ethBlocks, order, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, fromBlock, &toBlock)
		if err != nil {
			return lastBlockSynced, err
		}
		blocks := convertArrayEthermanBlocks(ethBlocks)
		err = s.blockRangeProcessor.ProcessBlockRange(s.ctx, ethBlocks, order)
		if err != nil {
			return lastBlockSynced, err
		}
		if len(blocks) > 0 {
			lastBlockSynced = &blocks[len(blocks)-1]
			for i := range blocks {
				log.Debug("NetworkID: ", s.networkID, ", Position: ", i, ". BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
			}
		}
		fromBlock = toBlock + 1

		if lastKnownBlock.Cmp(new(big.Int).SetUint64(toBlock)) < 1 {
			if !s.synced {
				log.Infof("NetworkID %d Synced!", s.networkID)
				waitDuration = s.cfg.SyncInterval.Duration
				//s.synced = true
				//s.chSynced <- s.networkID
			}
			break
		}
		if len(blocks) == 0 { // If there is no events in the checked blocks range and lastKnownBlock > fromBlock.
			// Store the latest block of the block range. Get block info and process the block
			fb, err := s.etherMan.EthBlockByNumber(s.ctx, toBlock)
			if err != nil {
				return lastBlockSynced, err
			}
			b := etherman.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			err = s.blockRangeProcessor.ProcessBlockRange(s.ctx, []etherman.Block{b}, order)
			if err != nil {
				return lastBlockSynced, err
			}

			lastBlockSynced = convertEthermanBlock(&b)
			log.Debugf("NetworkID: %d, Storing empty block. BlockNumber: %d. BlockHash: %s",
				s.networkID, b.BlockNumber, b.BlockHash.String())
		}
	}

	return lastBlockSynced, nil
}

// This function allows reset the state until an specific ethereum block
func (s *SynchronizerImpl) resetState(blockNumber uint64) error {
	log.Infof("NetworkID: %d. Reverting synchronization to block: %d", s.networkID, blockNumber)
	dbTx, err := s.storage.BeginDBTransaction(s.ctx)
	if err != nil {
		log.Errorf("networkID: %d, Error starting a db transaction to reset the state. Error: %v", s.networkID, err)
		return err
	}
	err = s.storage.Reset(s.ctx, blockNumber, dbTx)
	if err != nil {
		log.Errorf("networkID: %d, error resetting the state. Error: %v", s.networkID, err)
		rollbackErr := s.storage.Rollback(s.ctx, dbTx)
		if rollbackErr != nil {
			log.Errorf("networkID: %d, error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %s",
				s.networkID, blockNumber, rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}

	err = s.storage.Commit(s.ctx, dbTx)
	if err != nil {
		log.Errorf("networkID: %d, error committing the resetted state. Error: %v", s.networkID, err)
		rollbackErr := s.storage.Rollback(s.ctx, dbTx)
		if rollbackErr != nil {
			log.Errorf("networkID: %d, error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %s",
				s.networkID, blockNumber, rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}

	return nil
}

/*
This function will check if there is a reorg.
As input param needs the last ethereum block synced. Retrieve the block info from the blockchain
to compare it with the stored info. If hash and hash parent matches, then no reorg is detected and return a nil.
If hash or hash parent don't match, reorg detected and the function will return the block until the sync process
must be reverted. Then, check the previous ethereum block synced, get block info from the blockchain and check
hash and has parent. This operation has to be done until a match is found.
*/
func (s *SynchronizerImpl) checkReorg(latestBlock *L1Block) (*L1Block, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestBlockSynced := *latestBlock
	var depth uint64
	for {
		block, err := s.etherMan.EthBlockByNumber(s.ctx, latestBlock.BlockNumber)
		if err != nil {
			log.Errorf("networkID: %d, error getting latest block synced from blockchain. Block: %d, error: %v",
				s.networkID, latestBlock.BlockNumber, err)
			return nil, err
		}
		if block.NumberU64() != latestBlock.BlockNumber {
			err = fmt.Errorf("networkID: %d, wrong ethereum block retrieved from blockchain. Block numbers don't match."+
				" BlockNumber stored: %d. BlockNumber retrieved: %d", s.networkID, latestBlock.BlockNumber, block.NumberU64())
			log.Error("error: ", err)
			return nil, err
		}
		// Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > s.genBlockNumber {
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => latestBlockNumber: ", latestBlock.BlockNumber)
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => latestBlockHash: ", latestBlock.BlockHash)
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => latestBlockHashParent: ", latestBlock.ParentHash)
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => BlockNumber: ", latestBlock.BlockNumber, block.NumberU64())
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => BlockHash: ", block.Hash())
			log.Info("NetworkID: ", s.networkID, ", [checkReorg function] => BlockHashParent: ", block.ParentHash())
			depth++
			log.Info("NetworkID: ", s.networkID, ", REORG: Looking for the latest correct block. Depth: ", depth)
			// Reorg detected. Getting previous block
			dbTx, err := s.storage.BeginDBTransaction(s.ctx)
			if err != nil {
				log.Errorf("networkID: %d, error creating db transaction to get previous blocks. Error: %v", s.networkID, err)
				return nil, err
			}
			latestBlock, err = s.getPreviousBlockFromStorage(s.ctx, depth, dbTx)
			errC := s.storage.Commit(s.ctx, dbTx)
			if errC != nil {
				log.Errorf("networkID: %d, error committing dbTx, err: %v", s.networkID, errC)
				rollbackErr := s.storage.Rollback(s.ctx, dbTx)
				if rollbackErr != nil {
					log.Errorf("networkID: %d, error rolling back state. RollbackErr: %v, err: %s",
						s.networkID, rollbackErr, errC.Error())
					return nil, rollbackErr
				}
				return nil, errC
			}
			if errors.Is(err, pgstorage.ErrStorageNotFound) {
				log.Warnf("networkID: %d, error checking reorg: previous block not found in db: %v", s.networkID, err)
				return &L1Block{}, nil
			} else if err != nil {
				log.Errorf("networkID: %d, error detected getting previous block: %v", s.networkID, err)
				return nil, err
			}
		} else {
			break
		}
	}
	if latestBlockSynced.BlockHash != latestBlock.BlockHash {
		log.Infof("NetworkID: %d, reorg detected in block: %d", s.networkID, latestBlockSynced.BlockNumber)
		return latestBlock, nil
	}
	log.Debugf("NetworkID: %d, no reorg detected", s.networkID)
	return nil, nil
}
