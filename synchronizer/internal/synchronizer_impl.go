package internal

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/elderberry"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/etrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/incaberry"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions/processor_manager"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_check_block"
	l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces"
)

// SynchronizerImpl connects L1 and L2
type SynchronizerImpl struct {
	etherMan       syncinterfaces.EthermanFullInterface
	storage        syncinterfaces.StorageInterface
	state          syncinterfaces.StateInterface
	ctx            context.Context
	cancelCtx      context.CancelFunc
	genBlockNumber uint64
	cfg            syncconfig.Config
	networkID      uint
	synced         bool

	blockRangeProcessor syncinterfaces.BlockRangeProcessor
	l1Sync              syncinterfaces.L1Syncer
	storageChecker      syncinterfaces.StorageCompatibilityChecker

	reorgCallback func(nreorgData ReorgExecutionResult)
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizerImpl(
	ctx context.Context,
	storage syncinterfaces.StorageInterface,
	state syncinterfaces.StateInterface,
	ethMan syncinterfaces.EthermanFullInterface,
	storageChecker syncinterfaces.StorageCompatibilityChecker,
	cfg syncconfig.Config) (*SynchronizerImpl, error) {
	ctx, cancel := context.WithCancel(ctx)
	networkID := uint(0)
	genesisBlockNumber, err := getGenesisBlockNumber(ctx, cfg.GenesisBlockNumber, ethMan)
	if err != nil {
		defer cancel()
		return nil, err
	}
	cfg.GenesisBlockNumber = genesisBlockNumber
	l1EventProcessors := newL1EventProcessor(state)
	blockRangeProcessor := NewBlockRangeProcessLegacy(storage, state, state, l1EventProcessors)

	finalizedBlockNumberFetcher := l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.FinalizedBlockNumber, 0)
	syncPointBlockNumberFecther := l1_check_block.NewSafeL1BlockNumberFetch(l1_check_block.StringToL1BlockPoint(cfg.SyncBlockProtection), cfg.SyncBlockProtectionOffset)
	reorgManager := l1sync.NewCheckReorgManager(ctx, ethMan, state)
	blocksRetriever := l1sync.NewBlockPointsRetriever(
		syncPointBlockNumberFecther,
		finalizedBlockNumberFetcher,
		ethMan,
	)
	checkl1blocks := l1_check_block.NewCheckL1BlockHash(ethMan, storage, finalizedBlockNumberFetcher)

	l1SequentialSync := l1sync.NewL1SequentialSync(blocksRetriever, ethMan, state,
		blockRangeProcessor, reorgManager,
		checkl1blocks,
		l1sync.L1SequentialSyncConfig{
			SyncChunkSize:                 cfg.SyncChunkSize,
			GenesisBlockNumber:            genesisBlockNumber,
			AllowEmptyBlocksAsCheckPoints: true,
		})

	sync := SynchronizerImpl{
		storage:             storage,
		state:               state,
		etherMan:            ethMan,
		ctx:                 ctx,
		cancelCtx:           cancel,
		genBlockNumber:      genesisBlockNumber,
		cfg:                 cfg,
		networkID:           networkID,
		storageChecker:      storageChecker,
		l1Sync:              l1SequentialSync,
		blockRangeProcessor: blockRangeProcessor,
	}
	state.AddOnReorgCallback(sync.OnReorgExecuted)

	err = sync.CheckStorage(ctx)
	if err != nil {
		defer cancel()
		return nil, err
	}
	return &sync, nil
}

func getGenesisBlockNumber(ctx context.Context, cfgGenesisBlockNumber uint64, ethMan syncinterfaces.EthermanFullInterface) (uint64, error) {
	genesisBlockNumber := cfgGenesisBlockNumber
	if genesisBlockNumber == 0 {
		firstBlock, err := ethMan.GetL1BlockUpgradeLxLy(ctx, nil)
		if err != nil {
			log.Errorf("Error getting the first block from the blockchain. Error: %v", err)
			return 0, err
		}
		log.Infof("First block from the blockchain: %d (ETROG)", firstBlock)
		genesisBlockNumber = firstBlock
	}
	return genesisBlockNumber, nil
}

func newL1EventProcessor(state syncinterfaces.StateInterface) *processor_manager.L1EventProcessors {
	builder := processor_manager.NewL1EventProcessorsBuilder()
	builder.Register(etrog.NewProcessorL1InfoTreeUpdate(state))
	etrogSequenceBatchesProcessor := etrog.NewProcessorL1SequenceBatches(state)
	builder.Register(etrogSequenceBatchesProcessor)
	builder.Register(incaberry.NewProcessorForkId(state))
	builder.Register(etrog.NewProcessorL1InitialSequenceBatches(state))
	builder.Register(elderberry.NewProcessorL1SequenceBatchesElderberry(etrogSequenceBatchesProcessor))
	return builder.Build()
}

var waitDuration = time.Duration(0)

// IsSynced returns true if the synchronizer is synced or false if it's not
func (s *SynchronizerImpl) IsSynced() bool {
	return s.synced
}

func (s *SynchronizerImpl) SetCallbackOnReorgDone(callback func(reorgData ReorgExecutionResult)) {
	s.reorgCallback = callback
}

// OnReorgExecuted this is a CB setted to state reorg
func (s *SynchronizerImpl) OnReorgExecuted(reorg model.ReorgExecutionResult) {
	log.Infof("Reorg executed! %s", reorg.String())
	if s.reorgCallback != nil {
		param := ReorgExecutionResult{
			FirstL1BlockNumberValidAfterReorg: &reorg.Request.FirstL1BlockNumberToKeep,
			ReasonError:                       reorg.Request.ReasonError,
		}
		log.Infof("Executing reorg callback in a goroutine")
		go s.reorgCallback(param)
	}

}

func (s *SynchronizerImpl) CheckStorage(ctx context.Context) error {
	if s.storageChecker == nil {
		return nil
	}
	return s.storageChecker.CheckAndUpdateStorage(ctx)
}

func (s *SynchronizerImpl) getLastL1BlockOnStorage(ctx context.Context) (*entities.L1Block, error) {
	l1block, err := s.storage.GetLastBlock(ctx, nil)
	if errors.Is(err, entities.ErrNotFound) {
		log.Infof("networkID: %d, error getting the latest block. No data on stored.", s.networkID)
		return nil, nil
	}
	return l1block, err
}

type SyncExecutionFlags uint64

const (
	FlagReturnOnSync SyncExecutionFlags = 1 << iota
	FlagReturnBeforeReorg
	FlagReturnAfterReorg
)

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *SynchronizerImpl) Sync(executionFlags SyncExecutionFlags) error {
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Infof("Synchronization started")
	s.synced = false

	lastBlockSynced, err := s.getLastL1BlockOnStorage(s.ctx)
	if err != nil {
		log.Fatalf("networkID: %d, unexpected error getting the latest block. Error: %s", s.networkID, err.Error())
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
			if lastBlockSynced, isSynced, err = s.l1Sync.SyncBlocks(s.ctx, lastBlockSynced); err != nil {
				log.Warnf("networkID: %d, error syncing blocks: %v", s.networkID, err)
				reorgError := common.CastReorgError(err)
				if reorgError != nil {
					if (executionFlags & FlagReturnBeforeReorg) != 0 {
						log.Infof("NetworkID: %d, Synchronization finished, returning before executing reorg %+v", s.networkID, reorgError)
						return err
					}
					err = s.executeReorg(reorgError)
					if err != nil {
						log.Errorf("networkID: %d, error resetting the state to a previous block. Error: %v", s.networkID, err)
						continue
					}
					if (executionFlags & FlagReturnAfterReorg) != 0 {
						log.Infof("NetworkID: %d, Synchronization finished, returning after reorg %+v", s.networkID, reorgError)
						return err
					}
				}

				lastBlockSynced, err = s.getLastL1BlockOnStorage(s.ctx)
				if err != nil {
					log.Fatalf("networkID: %d, error getting lastBlockSynced to resume the synchronization... Error: ", s.networkID, err)
				}
				if s.ctx.Err() != nil {
					continue
				}
			}
			s.setSyncedStatus(isSynced)
			if s.synced {
				log.Infof("NetworkID %d Synced!   lastBlockSynced:%d ", s.networkID, lastBlockSynced.BlockNumber)
				if (executionFlags & FlagReturnOnSync) != 0 {
					log.Infof("NetworkID: %d, Synchronization finished, returning because returnOnSync=true", s.networkID)
					return nil
				}
				waitDuration = s.cfg.SyncInterval.Duration
			}
		}
	}
}

func (s *SynchronizerImpl) setSyncedStatus(synced bool) {
	s.synced = synced
}

// Stop function stops the synchronizer
func (s *SynchronizerImpl) Stop() {
	s.cancelCtx()
}
func (s *SynchronizerImpl) executeReorg(reorgError *common.ReorgError) error {
	if reorgError == nil {
		return nil
	}
	dbTx, err := s.state.BeginTransaction(s.ctx)
	if err != nil {
		log.Errorf("networkID: %d, error starting a db transaction to execute reorg. Error: %v", s.networkID, err)
		return err
	}

	req := model.ReorgRequest{
		FirstL1BlockNumberToKeep: reorgError.BlockNumber - 1, // Previous block to last bad block
		ReasonError:              reorgError,
	}

	result := s.state.ExecuteReorg(s.ctx, req, dbTx)
	if !result.IsSuccess() {
		log.Errorf("networkID: %d, error executing reorg. Error: %v", s.networkID, result.ExecutionError)
		// I don't care about result of Rollback
		_ = dbTx.Rollback(s.ctx)
		return result.ExecutionError
	}
	// This commit is going to launch the callback to OnReorgExecuted
	errCommit := dbTx.Commit(s.ctx)
	if errCommit != nil {
		log.Errorf("networkID: %d, error committing reorg. Error: %v", s.networkID, errCommit)
		return errCommit
	}
	return nil

}
