package synchronizer

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
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
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]L1InfoTreeLeaf, error)
}

type SequencedBatches struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	L1BlockNumber   uint64
	Timestamp       time.Time
}
type SynchronizerSequencedBatchesQuerier interface {
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64) (*SequencedBatches, error)
}

type Synchronizer interface {
	SynchronizerRunner
	SynchornizerStatusQuerier
	SynchronizerL1InfoTreeQuerier
	SynchronizerSequencedBatchesQuerier
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
	configStorage := pgstorage.Config{
		Name:     config.DB.Name,
		User:     config.DB.User,
		Password: config.DB.Password,
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		MaxConns: config.DB.MaxConns,
	}

	storage, err := pgstorage.NewPostgresStorage(configStorage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	etherman, err := etherman.NewClient(config.Etherman)
	if err != nil {
		log.Error("Error creating etherman", err)
		return nil, err
	}
	sync, err := NewSynchronizerImpl(ctx, storage, etherman, config.Synchronizer)
	if err != nil {
		log.Error("Error creating synchronizer", err)
		return nil, err
	}
	return sync, nil
}
