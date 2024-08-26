package synchronizer

import (
	internal "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/internal"
)

type SynchronizerAdapter struct {
	*SyncrhronizerQueries
	internalSyncrhonizer *internal.SynchronizerImpl
}

func NewSynchronizerAdapter(queries *SyncrhronizerQueries, sync *internal.SynchronizerImpl) *SynchronizerAdapter {
	return &SynchronizerAdapter{
		SyncrhronizerQueries: queries,
		internalSyncrhonizer: sync,
	}
}

func (s *SynchronizerAdapter) SetCallbackOnReorgDone(callback func(reorgData ReorgExecutionResult)) {
	s.internalSyncrhonizer.SetCallbackOnReorgDone(
		func(nreorgData internal.ReorgExecutionResult) {
			callback(ReorgExecutionResult{
				FirstL1BlockNumberValidAfterReorg: nreorgData.FirstL1BlockNumberValidAfterReorg,
				ReasonError:                       nreorgData.ReasonError,
			})
		})
}

func (s *SynchronizerAdapter) SetCallbackOnRollbackBatches(callback func(data RollbackBatchesData)) {
}

func (s *SynchronizerAdapter) Sync(returnOnSync bool) error {
	var flags internal.SyncExecutionFlags
	if returnOnSync {
		flags = internal.FlagReturnOnSync
	}
	return s.internalSyncrhonizer.Sync(flags)
}

func (s *SynchronizerAdapter) Stop() {
	s.internalSyncrhonizer.Stop()
}

func (s *SynchronizerAdapter) IsSynced() bool {
	return s.internalSyncrhonizer.IsSynced()
}
