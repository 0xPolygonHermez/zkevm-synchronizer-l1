package syncinterfaces

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/ethereum/go-ethereum/common"
)

type StorageBlockWriterInterface interface {
	AddBlock(ctx context.Context, block *entities.L1Block, dbTx stateTxType) error
	UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx stateTxType) error
}

type StorageBlockReaderInterface interface {
	GetLastBlock(ctx context.Context, dbTx stateTxType) (*entities.L1Block, error)
	AddBlock(ctx context.Context, block *entities.L1Block, dbTx stateTxType) error
	GetPreviousBlock(ctx context.Context, offset uint64, fromBlockNumber *uint64, dbTx stateTxType) (*entities.L1Block, error)
	GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx stateTxType) (*entities.L1Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx stateTxType) (*entities.L1Block, error)
}

type StorageForkIDInterface interface {
	AddForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx stateTxType) error
	GetForkIDs(ctx context.Context, dbTx stateTxType) ([]pgstorage.ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID pgstorage.ForkIDInterval, dbTx stateTxType) error
}

type StorageL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx stateTxType) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx stateTxType) ([]pgstorage.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx stateTxType) (*pgstorage.L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx stateTxType) (*pgstorage.L1InfoTreeLeaf, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx stateTxType) ([]pgstorage.L1InfoTreeLeaf, error)
}

type StorageVirtualBatchInterface interface {
	GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) (*pgstorage.VirtualBatch, error)
	GetLastestVirtualBatchNumber(ctx context.Context, constrains *pgstorage.VirtualBatchConstraints, dbTx stateTxType) (uint64, error)
}

type StorageSequenceBatchesInterface interface {
	AddSequencedBatches(ctx context.Context, sequence *pgstorage.SequencedBatches, dbTx stateTxType) error
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) (*pgstorage.SequencedBatches, error)
}

// StorageKVInterface is an interface for key-value storage

type StateForkIdQuerier interface {
	GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) uint64
	GetForkIDByBlockNumber(ctx context.Context, blockNumber uint64, dbTx stateTxType) uint64
}

type StorageInterface interface {
	StorageBlockWriterInterface
	StorageBlockReaderInterface
	StorageForkIDInterface
	StorageL1InfoTreeInterface
	StorageSequenceBatchesInterface
}
