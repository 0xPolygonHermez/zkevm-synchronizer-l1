package synchronizer

import (
	"context"
	"errors"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/rpcsync"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage"
	internal "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/internal"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// ErrNotFound is used when the object is not found
	ErrNotFound = errors.New("not found")
)

type L1InfoTreeLeaf struct {
	L1InfoTreeRoot    common.Hash
	L1InfoTreeIndex   uint32
	PreviousBlockHash common.Hash
	BlockNumber       uint64
	Timestamp         time.Time
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
}

type SynchronizerRunner interface {
	// Sync is blocking call, must be launched as a goroutine
	// If returnOnSync is true, it will return when the synchronizer is synced,
	//  otherwise it will keep running
	Sync(returnOnSync bool) error
	// Stop stops the synchronizer
	Stop()
}

type SynchornizerStatusQuerier interface {
	// IsSynced returns true if the synchronizer is synced or false if it's not
	IsSynced() bool
}

type SynchronizerL1InfoTreeQuerier interface {
	// GetL1InfoRootPerIndex returns the L1InfoTree root hash for a given index
	// if not found returns ErrNotFound
	GetL1InfoRootPerIndex(ctx context.Context, L1InfoTreeIndex uint32) (common.Hash, error)
	GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32) (map[uint32]L1InfoTreeLeaf, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash) ([]L1InfoTreeLeaf, error)
}

type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool // The block is safe (have past the safe point, e.g. Finalized in L1)
	HasEvents   bool // This block have events from the rollup
	SyncVersion string
}

type SynchronizerBlockQuerier interface {
	GetLastL1Block(ctx context.Context) (*L1Block, error)
	GetL1BlockByNumber(ctx context.Context, blockNumber uint64) (*L1Block, error)
}

type SequencedBatches struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	L1BlockNumber   uint64
	ForkID          uint64
	Timestamp       time.Time
	ReceivedAt      time.Time
	L1InfoRoot      common.Hash
	Source          string
}
type SynchronizerSequencedBatchesQuerier interface {
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64) (*SequencedBatches, error)
}

type VirtualBatch struct {
	BatchNumber             uint64
	ForkID                  uint64
	BatchL2Data             []byte
	VlogTxHash              common.Hash // Hash of tx inside L1Block that emit this log
	Coinbase                common.Address
	SequencerAddr           common.Address
	SequenceFromBatchNumber uint64 // Linked to sync.sequenced_batches table
	BlockNumber             uint64 // Linked to sync.block table
	L1InfoRoot              *common.Hash
	ReceivedAt              time.Time
	BatchTimestamp          *time.Time // This is optional depend on ForkID
	ExtraInfo               *string
}

type SynchronizerVirtualBatchesQuerier interface {
	GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64) (*VirtualBatch, error)
	GetLastestVirtualBatchNumber(ctx context.Context) (uint64, error)
}

type ReorgExecutionResult struct {
	// FirstL1BlockNumberValidAfterReorg is the first block or nil if the reorg have delete all blocks
	FirstL1BlockNumberValidAfterReorg *uint64
	ReasonError                       error
}

// SynchronizerReorgSupporter is an interface that give support to the reorgs detected on L1
type SynchronizerReorgSupporter interface {
	// SetCallbackOnReorgDone sets a callback that will be called when the reorg is done
	// to disable it you can set nil
	SetCallbackOnReorgDone(callback func(reorgData ReorgExecutionResult))
}

type RollbackBatchesData struct {
	LastBatchNumber       uint64
	LastBatchAccInputHash common.Hash
}

// SynchronizerRollbackSupporter is an interface that give support to the banna rollbackBatches
type SynchronizerRollbackBatchesSupporter interface {
	// SetCallbackOnRollbackBatches sets a callback that will be called when the rollbackBatches  is done
	// to disable it you can set nil
	SetCallbackOnRollbackBatches(callback func(data RollbackBatchesData))
}

type Synchronizer interface {
	SynchronizerRunner
	SynchornizerStatusQuerier
	SynchronizerL1InfoTreeQuerier
	SynchronizerSequencedBatchesQuerier
	SynchronizerReorgSupporter
	SynchronizerVirtualBatchesQuerier
	SynchronizerBlockQuerier
	SynchronizerRollbackBatchesSupporter
}

func NewSynchronizerFromConfigfile(ctx context.Context, configFile string) (Synchronizer, error) {
	config, err := config.LoadFile(configFile)
	if err != nil || config == nil {
		log.Error("Error loading config", err)
		return nil, err
	}
	log.Init(config.Log)
	return NewSynchronizer(ctx, *config)
}

func NewSynchronizer(ctx context.Context, config config.Config) (Synchronizer, error) {
	log.Init(config.Log)
	log.Debugf("Creating storage")
	storage, err := storage.NewStorage(config.SQLDB)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debugf("Creating etherman")
	etherman, err := etherman.NewClient(config.Etherman)
	if err != nil {
		log.Error("Error creating etherman", err)
		return nil, err
	}
	log.Debugf("Creating state")
	state := state.NewState(storage)
	storageCompatibilityChecker := internal.NewSanityStorageCheckerImpl(state, etherman, config.Synchronizer.OverrideStorageCheck)
	LastExecutionChecker := newLastExecutionState(storage, config)
	sync, err := internal.NewSynchronizerImpl(ctx, storage, state, etherman, storageCompatibilityChecker, LastExecutionChecker, config.Synchronizer)
	if err != nil {
		log.Error("Error creating synchronizer", err)
		return nil, err
	}
	log.Debugf("Creating synchronizer adapter")
	syncAdapter := NewSynchronizerAdapter(NewSyncrhronizerQueries(state, storage, ctx), sync)
	log.Debugf("Starting RPC if enabled")
	rpcsync.StartRPC(state)

	return syncAdapter, nil
}

func newLastExecutionState(storage storage.KvStorer, cfg config.Config) *model.LastExecutionState {
	configStr, err := config.SaveConfigToString(cfg)
	if err != nil {
		log.Fatal("Error saving config to string", err)
	}
	return model.NewLastExecutionState(storage, configStr, zkevm_synchronizer_l1.Version)
}
