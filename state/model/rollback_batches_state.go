package model

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
)

const maxRewindBlocks = 1000

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
	allowReorg                 bool
}

func NewRollbackBatchesState(storage StorageRollbackBatchesInterface, allowReorg bool) *RollbackBatchesState {
	return &RollbackBatchesState{
		storage:    storage,
		allowReorg: allowReorg,
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

	err = sanityCheckAffectedSequences(affectedSequences, rollbackBatchesRequest)
	if err != nil {
		log.Errorf("error sanity check : %s", err.Error())
		return nil, err
	}

	description := fmt.Sprintf("RollbackBatchesState: %d sequences affected, %d batches affected",
		affectedSequences.Len(), affectedSequences.NumBatchesIncluded())
	log.Debug(description)
	if affectedSequences == nil {
		affectedSequences = &SequencesBatchesSlice{}
	}
	rollbackBatchesEntry := &RollbackBatchesLogEntry{
		BlockNumber:           rollbackBatchesRequest.L1BlockNumber,
		LastBatchNumber:       rollbackBatchesRequest.LastBatchNumber,
		LastBatchAccInputHash: rollbackBatchesRequest.LastBatchAccInputHash,
		L1EventAt:             rollbackBatchesRequest.L1BlockTimestamp,
		ReceivedAt:            time.Now(),
		UndoFirstBlockNumber:  affectedSequences.GetMinimumBlockNumber(rollbackBatchesRequest.L1BlockNumber),
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

func sanityCheckAffectedSequences(affectedSequences *SequencesBatchesSlice, rollbackBatchesRequest RollbackBatchesRequest) error {
	if affectedSequences.Len() == 0 {
		log.Warnf("No sequences affected for rollback batches (LastBatchNumber=%d)", rollbackBatchesRequest.LastBatchNumber)
		return nil
	}
	// The first batch after rollback must be a first batch of a sequence
	for _, seq := range *affectedSequences {
		if seq.FromBatchNumber == rollbackBatchesRequest.LastBatchNumber+1 {
			return nil
		}
	}
	err := fmt.Errorf("cant break a sequence, batchNumber %d is not a initial one (%v)", rollbackBatchesRequest.LastBatchNumber+1, affectedSequences)
	log.Error(err.Error())
	return err
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

	if !s.allowReorg {
		err = fmt.Errorf("pre_execute_reorg: reorg not possible (L1BlockNumber>=%d), affect to rollback batches! (%v)", reorgRequest.FirstL1BlockNumberToKeep+1, entries)
		log.Error(err.Error())
		return err
	}
	// Todo a real reorg you must get the lowest UndoFirstBlockNumber and execute the reorg from this point
	oldFirstL1BlockNumberToKeep := reorgRequest.FirstL1BlockNumberToKeep
	// Get the undo_first_block_number that is the block where the first seq deleted where declared
	//   so we need to delete this block (that is the reason of -1)
	newFirstL1BlockNumberToKeep := getMinimumUndoFirstBlockNumber(entries) - 1
	if oldFirstL1BlockNumberToKeep-newFirstL1BlockNumberToKeep > maxRewindBlocks {
		err = fmt.Errorf("pre_execute_reorg: reorg not possible (L1BlockNumber>=%d), too many blocks to rewind! (%v) toBlock:%d", reorgRequest.FirstL1BlockNumberToKeep+1, oldFirstL1BlockNumberToKeep-newFirstL1BlockNumberToKeep, newFirstL1BlockNumberToKeep)
		log.Error(err.Error())
		return err
	}
	if newFirstL1BlockNumberToKeep < oldFirstL1BlockNumberToKeep {
		reorgRequest.FirstL1BlockNumberToKeep = newFirstL1BlockNumberToKeep
		log.Infof("pre_execute_reorg: To be able to reorg we need to keep block_number=%d instead of %d ", newFirstL1BlockNumberToKeep, oldFirstL1BlockNumberToKeep)
		return nil
	}
	return nil

}

func getMinimumUndoFirstBlockNumber(entries []RollbackBatchesLogEntry) uint64 {
	min := uint64(0)
	for _, entry := range entries {
		if min == 0 || entry.UndoFirstBlockNumber < min {
			min = entry.UndoFirstBlockNumber
		}
	}
	return min
}
