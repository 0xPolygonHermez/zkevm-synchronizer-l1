package synchronizer

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/elderberry"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_check_block"
	l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync"
)

// SynchronizerImpl connects L1 and L2
type SynchronizerImpl struct {
	etherMan EthermanInterface
	// TODO: remove

	storage        syncinterfaces.StorageInterface
	state          syncinterfaces.StateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            syncconfig.Config
	networkID      uint
	synced         bool

	l1EventProcessors   *processor_manager.L1EventProcessors
	blockRangeProcessor syncinterfaces.BlockRangeProcessor
	l1Sync              *l1sync.L1SequentialSync
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizerImpl(
	ctx context.Context,
	storage syncinterfaces.StorageInterface,
	state syncinterfaces.StateInterface,
	ethMan EthermanInterface,
	cfg syncconfig.Config) (*SynchronizerImpl, error) {
	ctx, cancel := context.WithCancel(ctx)
	networkID := uint(0)

	sync := SynchronizerImpl{
		storage:        storage,
		state:          state,
		etherMan:       ethMan,
		ctx:            ctx,
		cancelCtx:      cancel,
		genBlockNumber: cfg.GenesisBlockNumber,
		cfg:            cfg,
		networkID:      networkID,
	}

	builder := processor_manager.NewL1EventProcessorsBuilder()
	builder.Register(etrog.NewProcessorL1InfoTreeUpdate(state))
	etrogSequenceBatchesProcessor := etrog.NewProcessorL1SequenceBatches(state)
	builder.Register(etrogSequenceBatchesProcessor)
	builder.Register(incaberry.NewProcessorForkId(state))
	builder.Register(etrog.NewProcessorL1InitialSequenceBatches(state))
	builder.Register(elderberry.NewProcessorL1SequenceBatchesElderberry(etrogSequenceBatchesProcessor))
	sync.l1EventProcessors = builder.Build()

	sync.blockRangeProcessor = NewBlockRangeProcessLegacy(storage, state, state, sync.l1EventProcessors)
	if cfg.GenesisBlockNumber == 0 {
		firstBlock, err := ethMan.GetL1BlockUpgradeLxLy(ctx, nil)
		if err != nil {
			log.Errorf("Error getting the first block from the blockchain. Error: %v", err)
			return nil, err
		}
		log.Infof("First block from the blockchain: %d (ETROG)", firstBlock)
		sync.genBlockNumber = firstBlock
	}
	//TODO: Add blockRetriever
	// TODO: Add Reorg object
	reorgManager := l1sync.NewCheckReorgManager(ctx, ethMan, state, nil, sync.genBlockNumber)
	blocksRetriever := l1sync.NewBlockPointsRetriever(
		l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.FinalizedBlockNumber, 0),
		l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.SafeBlockNumber, 0),
		ethMan,
	)

	sync.l1Sync = l1sync.NewL1SequentialSync(blocksRetriever, ethMan, state, sync.blockRangeProcessor, reorgManager,
		l1sync.L1SequentialSyncConfig{
			SyncChunkSize:      cfg.SyncChunkSize,
			GenesisBlockNumber: sync.genBlockNumber,
		})
	return &sync, nil

}

var waitDuration = time.Duration(0)

// IsSynced returns true if the synchronizer is synced or false if it's not
func (s *SynchronizerImpl) IsSynced() bool {
	return s.synced
}

func (s *SynchronizerImpl) SetCallbackOnReorgDone(callback func(newFirstL1BlockNumberValid uint64)) {
	//TODO: Implement this function
	log.Fatal("Not implemented")
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *SynchronizerImpl) Sync(returnOnSync bool) error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Infof("Synchronization started")
	s.synced = false

	lastBlockSynced, err := s.storage.GetLastBlock(s.ctx, nil)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			//log.Infof("networkID: %d, error getting the latest ethereum block. No data stored. Setting genesis block. Error: %v", s.networkID, err)
			lastBlockSynced = &entities.L1Block{
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
			log.Infof("synchronizer ctx done")
			return nil
		case <-time.After(waitDuration):
			log.Debugf("syncing...")
			//Sync L1Blocks

			var isSynced bool
			//if lastBlockSynced, isSynced, err = s.syncBlocks(lastBlockSynced); err != nil {
			if lastBlockSynced, isSynced, err = s.l1Sync.SyncBlocksSequential(s.ctx, lastBlockSynced); err != nil {
				log.Warnf("networkID: %d, error syncing blocks: %v", s.networkID, err)
				lastBlockSynced, err = s.storage.GetLastBlock(s.ctx, nil)
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
				log.Debugf("NetworkID: %d, lastBlockSynced: %d, lastKnownBlock: %d", s.networkID, lastBlockSynced.BlockNumber, lastKnownBlock)
				if isSynced && !s.synced {
					log.Infof("NetworkID %d Synced!  lastL1Block: %d lastBlockSynced:%d ", s.networkID, lastKnownBlock, lastBlockSynced.BlockNumber)
					waitDuration = s.cfg.SyncInterval.Duration
					s.synced = true
					if returnOnSync {
						log.Infof("NetworkID: %d, Synchronization finished, returning because returnOnSync=true", s.networkID)
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
			} else {
				s.synced = isSynced
			}
		}
	}
}

// Stop function stops the synchronizer
func (s *SynchronizerImpl) Stop() {
	s.cancelCtx()
}

// This function allows reset the state until an specific ethereum block
func (s *SynchronizerImpl) resetState(blockNumber uint64) error {
	log.Infof("NetworkID: %d. Reverting synchronization to block: %d", s.networkID, blockNumber)
	dbTx, err := s.state.BeginTransaction(s.ctx)
	if err != nil {
		log.Errorf("networkID: %d, Error starting a db transaction to reset the state. Error: %v", s.networkID, err)
		return err
	}
	err = s.storage.Reset(s.ctx, blockNumber, dbTx)
	if err != nil {
		log.Errorf("networkID: %d, error resetting the state. Error: %v", s.networkID, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("networkID: %d, error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %s",
				s.networkID, blockNumber, rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}

	err = dbTx.Commit(s.ctx)
	if err != nil {
		log.Errorf("networkID: %d, error committing the resetted state. Error: %v", s.networkID, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("networkID: %d, error rolling back state to store block. BlockNumber: %d, rollbackErr: %v, error : %s",
				s.networkID, blockNumber, rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}

	return nil
}
