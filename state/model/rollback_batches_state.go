package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
)

type RollbackBatchesRequest struct {
	LastBatchNumber       uint64
	LastBatchAccInputHash common.Hash

	L1BlockNumber    uint64    // L1BlockNumber  block where the log where register
	L1BlockTimestamp time.Time // L1BlockTimestamp timestamp of the L1 block
}

type StorageRollbackBatchesInterface interface {
	GetSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx storageTxType) (*SequencesBatchesSlice, error)
	DeleteSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx storageTxType) error
	AddRollbackBatchesLogEntry(ctx context.Context, entry *RollbackBatchesLogEntry, dbTx dbTxType) error
	GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx context.Context, l1BlockNumber uint64, dbTx dbTxType) ([]RollbackBatchesLogEntry, error)
}

type RollbackBatchesExecutionResult struct {
	Request       RollbackBatchesRequest
	RollbackEntry *RollbackBatchesLogEntry
}

func (r *RollbackBatchesExecutionResult) String() string {
	return fmt.Sprintf("RollbackBatchesExecutionResult{Request: %v, RollbackEntry: %v}", r.Request, r.RollbackEntry)
}

type RollbackBatchesCallbackType = func(RollbackBatchesExecutionResult)

type RollbackBatchesState struct {
	mutex                      sync.Mutex
	storage                    StorageRollbackBatchesInterface
	onRollbackBatchesCallbacks []RollbackBatchesCallbackType
}

func NewRollbackBatchesState(storage StorageRollbackBatchesInterface) *RollbackBatchesState {
	return &RollbackBatchesState{
		storage: storage,
	}
}

func (s *RollbackBatchesState) AddOnRollbackBatchesCallback(f RollbackBatchesCallbackType) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onRollbackBatchesCallbacks = append(s.onRollbackBatchesCallbacks, f)
}

func (s *RollbackBatchesState) onTxCommit(data *RollbackBatchesExecutionResult, _ storageTxType, err error) {
	if err == nil {
		for _, f := range s.onRollbackBatchesCallbacks {
			f(*data)
		}
	}
}

func (s *RollbackBatchesState) ExecuteRollbackBatches(ctx context.Context, rollbackBatchesRequest RollbackBatchesRequest, dbTx storageTxType) (*RollbackBatchesExecutionResult, error) {
	if dbTx == nil {
		return nil, fmt.Errorf("for execute rollback batches, dbTx must be not nil because is used for callback")
	}
	response := &RollbackBatchesExecutionResult{
		Request: rollbackBatchesRequest,
	}

	affectedSequences, err := s.storage.GetSequencesGreatestOrEqualBatchNumber(ctx, rollbackBatchesRequest.LastBatchNumber+1, dbTx)
	if err != nil {
		err = fmt.Errorf("error getting affected sequences (batchNumber>=%d): %w", rollbackBatchesRequest.LastBatchNumber+1, err)
		log.Error(err.Error())
		return nil, err
	}

	description := fmt.Sprintf("RollbackBatchesState: %d sequences affected, %d batches affected",
		affectedSequences.Len(), affectedSequences.NumBatchesIncluded())
	rollbackBatchesEntry := &RollbackBatchesLogEntry{
		BlockNumber:           rollbackBatchesRequest.L1BlockNumber,
		LastBatchNumber:       rollbackBatchesRequest.LastBatchNumber,
		LastBatchAccInputHash: rollbackBatchesRequest.LastBatchAccInputHash,
		L1EventAt:             rollbackBatchesRequest.L1BlockTimestamp,
		ReceivedAt:            time.Now(),
		UndoFirstBlockNumber:  affectedSequences.GetMinimumBlockNumber(),
		Description:           description,
		SequencesDeleted:      *affectedSequences,
	}
	err = s.storage.AddRollbackBatchesLogEntry(ctx, rollbackBatchesEntry, dbTx)
	if err != nil {
		err = fmt.Errorf("error writing rollback batches log entry: %w", err)
		log.Error(err.Error())
		return nil, err
	}
	// Delete sequences delete also virtual batches (on cascade)
	err = s.storage.DeleteSequencesGreatestOrEqualBatchNumber(ctx, rollbackBatchesRequest.LastBatchNumber+1, dbTx)
	if err != nil {
		err = fmt.Errorf("error deleting affected sequences (batchNumber>=%d): %w", rollbackBatchesRequest.LastBatchNumber+1, err)
		log.Error(err.Error())
		return nil, err
	}
	response.RollbackEntry = rollbackBatchesEntry
	// Add commit callback to execute the onRollbackBatchesCallbacks
	dbTx.AddCommitCallback(func(dbTx storageTxType, err error) { s.onTxCommit(response, dbTx, err) })
	return response, nil
}

// PreExecuteReorg check that is possible to execute the reorg
// or add more info to the reorg request
func (s *RollbackBatchesState) PreExecuteReorg(ctx context.Context, reorgRequest *ReorgRequest, dbTx storageTxType) error {
	entries, err := s.storage.GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx, reorgRequest.FirstL1BlockNumberToKeep+1, dbTx)
	if err != nil {
		err = fmt.Errorf("pre_execute_reorg: error getting rollback batches log entries (L1BlockNumber>=%d): %w", reorgRequest.FirstL1BlockNumberToKeep+1, err)
		log.Error(err.Error())
		return err
	}
	if len(entries) == 0 {
		log.Infof("pre_execute_reorg: no rollback batches log entries affected for reorg (L1BlockNumber>=%d)", reorgRequest.FirstL1BlockNumberToKeep+1)
		return nil
	}
	// Todo a real reorg you must get the lowest UndoFirstBlockNumber and execute the reorg from this point
	err = fmt.Errorf("pre_execute_reorg: reorg not possible (L1BlockNumber>=%d), affect to rollback batches! (%v)", reorgRequest.FirstL1BlockNumberToKeep+1, entries)
	log.Error(err.Error())
	return err

}
