package storage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
)

type L1Block = entities.L1Block
type L1InfoTreeLeaf = entities.L1InfoTreeLeaf
type ForkIDInterval = entities.ForkIDInterval
type VirtualBatch = entities.VirtualBatch
type SequencedBatches = entities.SequencedBatches
type storageTxType = entities.Tx

type BlockStorer interface {
	AddBlock(ctx context.Context, block *L1Block, dbTx storageTxType) error
	GetLastBlock(ctx context.Context, dbTx storageTxType) (*L1Block, error)
	GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx storageTxType) (*L1Block, error)
	GetPreviousBlock(ctx context.Context, offset uint64, fromBlockNumber *uint64, dbTx storageTxType) (*L1Block, error)
	GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx storageTxType) (*L1Block, error)
	GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx storageTxType) (*[]L1Block, error)
}

type forkidStorer interface {
	AddForkID(ctx context.Context, forkID ForkIDInterval, dbTx storageTxType) error
	GetForkIDs(ctx context.Context, dbTx storageTxType) ([]ForkIDInterval, error)
	UpdateForkID(ctx context.Context, forkID ForkIDInterval, dbTx storageTxType) error
	GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx storageTxType) uint64
}

type l1infoTreeStorer interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *L1InfoTreeLeaf, dbTx storageTxType) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx storageTxType) ([]L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx storageTxType) (*L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx storageTxType) (*L1InfoTreeLeaf, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx storageTxType) ([]L1InfoTreeLeaf, error)
}

type sequencedBatchStorer interface {
	AddSequencedBatches(ctx context.Context, sequence *SequencedBatches, dbTx storageTxType) error
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx storageTxType) (*SequencedBatches, error)
}

type virtualBatchStorer interface {
	AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx storageTxType) error
	GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx storageTxType) (*VirtualBatch, error)
}

type reorgStorer interface {
	ResetToL1BlockNumber(ctx context.Context, firstBlockNumberToKeep uint64, dbTx storageTxType) error
}

type txStorer interface {
	BeginTransaction(ctx context.Context) (storageTxType, error)
}

type Storer interface {
	txStorer
	BlockStorer
	forkidStorer
	l1infoTreeStorer
	virtualBatchStorer
	sequencedBatchStorer
	reorgStorer
}
