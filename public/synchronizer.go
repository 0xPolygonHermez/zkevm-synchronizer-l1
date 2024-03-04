package public

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx"
)

type Config struct {
	configDB       db.Config
	configSync     synchronizer.Config
	configEtherman etherman.Config
}

type Synchronizer interface {
	SynchronizerRunner
	SynchornizerStatusQuery
	SynchronizerL1InfoTreeQuery
}

type SynchronizerRunner interface {
	// Sync is blocking call, must be launched as a goroutine
	Sync() error
	// Stop stops the synchronizer
	Stop()
}

type SynchornizerStatusQuery interface {
	// IsSynced returns true if the synchronizer is synced or false if it's not
	IsSynced() bool
}

type SynchronizerL1InfoTreeQuery interface {
	GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error)
}

func NewSynchronizer(ctx context.Context, cfg Config) (Synchronizer, error) {
	err := db.RunMigrations(cfg.configDB)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	configStorage := pgstorage.Config{
		Name:     cfg.configDB.Name,
		User:     cfg.configDB.User,
		Password: cfg.configDB.Password,
		Host:     cfg.configDB.Host,
		Port:     cfg.configDB.Port,
		MaxConns: cfg.configDB.MaxConns,
	}
	storage, err := pgstorage.NewPostgresStorage(configStorage)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	etherman, err := etherman.NewClient(cfg.configEtherman)
	if err != nil {
		log.Error("Error creating etherman", err)
		return nil, err
	}
	sync, err := synchronizer.NewSynchronizer(ctx, storage, etherman, cfg.configSync)
	if err != nil {
		log.Error("Error creating synchronizer", err)
		return nil, err
	}
	return NewSynchronizerAdapter(sync, sync, storage), nil
}
