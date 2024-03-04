package public

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx"
)

type SynchronizerAdapter struct {
	SyncRunner        SynchronizerRunner
	SyncStatus        SynchornizerStatusQuery
	storageL1InfoTree SynchronizerL1InfoTreeQuery
}

func NewSynchronizerAdapter(syncRunner SynchronizerRunner, syncStatus SynchornizerStatusQuery, storageL1InfoTree SynchronizerL1InfoTreeQuery) *SynchronizerAdapter {
	return &SynchronizerAdapter{
		SyncRunner:        syncRunner,
		SyncStatus:        syncStatus,
		storageL1InfoTree: storageL1InfoTree,
	}
}

func (s *SynchronizerAdapter) Sync() error {
	return s.SyncRunner.Sync()
}

func (s *SynchronizerAdapter) Stop() {
	s.SyncRunner.Stop()
}

func (s *SynchronizerAdapter) IsSynced() bool {
	return s.SyncStatus.IsSynced()
}

func (s *SynchronizerAdapter) GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	return s.storageL1InfoTree.GetL1InfoRootLeafByIndex(ctx, l1InfoTreeIndex, dbTx)
}

func (s *SynchronizerAdapter) GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	return s.storageL1InfoTree.GetLeafsByL1InfoRoot(ctx, l1InfoRoot, dbTx)
}
