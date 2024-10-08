package syncinterfaces

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/ethereum/go-ethereum/common"
)

type StateTxProvider interface {
	BeginTransaction(ctx context.Context) (stateTxType, error)
}

type StateForkidQuerier interface {
	GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) uint64
	GetForkIDByBlockNumber(ctx context.Context, blockNumber uint64, dbTx stateTxType) uint64
}

type stateOnSequencedBatchesManager interface {
	OnSequencedBatchesOnL1(ctx context.Context, seq model.SequenceOfBatches, dbTx stateTxType) error
}

type stateReorgManager interface {
	AddOnReorgCallback(f model.ReorgCallbackType)
	ExecuteReorg(ctx context.Context, reorgRequest model.ReorgRequest, dbTx stateTxType) model.ReorgExecutionResult
}

type StateForkIdQuerier interface {
	GetForkIDByBatchNumber(ctx context.Context, batchNumber uint64, dbTx stateTxType) uint64
	GetForkIDByBlockNumber(ctx context.Context, blockNumber uint64, dbTx stateTxType) uint64
}

type StateRollbackBatchesExecutor interface {
	ExecuteRollbackBatches(ctx context.Context, rollbackBatchesRequest model.RollbackBatchesRequest, dbTx stateTxType) (*model.RollbackBatchesExecutionResult, error)
	AddOnRollbackBatchesCallback(f model.RollbackBatchesCallbackType)
}

type StateInterface interface {
	AddL1InfoTreeLeafAndAssignIndex(ctx context.Context, exitRoot *entities.L1InfoTreeLeaf, dbTx stateTxType) (*entities.L1InfoTreeLeaf, error)

	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx stateTxType) ([]entities.L1InfoTreeLeaf, error)
	GetL1InfoRootPerLeafIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx stateTxType) (common.Hash, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx stateTxType) (*entities.L1InfoTreeLeaf, error)
	GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32, dbTx stateTxType) (map[uint32]entities.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx stateTxType) (*entities.L1InfoTreeLeaf, error)

	AddForkID(ctx context.Context, newForkID entities.ForkIDInterval, dbTx stateTxType) error

	StateForkidQuerier
	StateTxProvider
	stateOnSequencedBatchesManager
	StorageBlockReaderInterface
	stateReorgManager
	StateRollbackBatchesExecutor
}

type StateStorageCompatibilityCheckerInterface interface {
	CheckAndUpdateStorage(ctx context.Context, runBoundData storageContentsBoundData, overrideStorageCheck bool, dbTx stateTxType) error
}
