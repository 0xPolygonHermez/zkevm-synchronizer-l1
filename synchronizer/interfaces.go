package synchronizer

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// EthermanInterface contains the methods required to interact with ethereum.
type EthermanInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
	//GetNetworkID(ctx context.Context) (uint, error)
	GetRollupID() uint
	GetL1BlockUpgradeLxLy(ctx context.Context, genesisBlock *uint64) (uint64, error)
	GetForks(ctx context.Context, genBlockNumber uint64, lastL1BlockSynced uint64) ([]etherman.ForkIDInterval, error)
}

type StorageBlockInterface interface {
	GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*pgstorage.L1Block, error)
	AddBlock(ctx context.Context, block *pgstorage.L1Block, dbTx pgx.Tx) error
	GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*pgstorage.L1Block, error)
}

type StorageResetInterface interface {
	Reset(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error
}

type StorageForkIDInterface interface {
	AddForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx pgx.Tx) error
	GetForkIDs(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx pgx.Tx) error
}

type StorageL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx pgx.Tx) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
}

type StorageTransactionInterface interface {
	Rollback(ctx context.Context, dbTx pgx.Tx) error
	BeginDBTransaction(ctx context.Context) (pgx.Tx, error)
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, dbTx pgx.Tx) error
}

type StorageSequenceBatchesInterface interface {
	AddSequencedBatches(ctx context.Context, sequence pgstorage.SequencedBatches, dbTx pgx.Tx) error
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*pgstorage.SequencedBatches, error)
}

type StorageInterface interface {
	StorageTransactionInterface
	StorageBlockInterface
	StorageResetInterface
	StorageForkIDInterface
	StorageL1InfoTreeInterface
	StorageSequenceBatchesInterface
}
