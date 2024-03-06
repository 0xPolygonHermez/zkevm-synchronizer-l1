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
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	GetForkIDByBlockNumber(blockNumber uint64) uint64
}

type StorageL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx pgx.Tx) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
}

type storageTransactionInterface interface {
	Rollback(ctx context.Context, dbTx pgx.Tx) error
	BeginDBTransaction(ctx context.Context) (pgx.Tx, error)
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context, dbTx pgx.Tx) error
}

type storageInterface interface {
	storageTransactionInterface
	StorageBlockInterface
	StorageResetInterface
	StorageForkIDInterface
	StorageL1InfoTreeInterface
}
