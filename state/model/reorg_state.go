package model

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ReorgRequest is a struct that contains the information needed to execute a reorg
type ReorgRequest struct {
	FirstL1BlockNumberToKeep uint64
	ReasonError              error
}

func (r *ReorgRequest) String() string {
	return fmt.Sprintf("FirstL1BlockNumberToKeep: %d, ReasonError: %s", r.FirstL1BlockNumberToKeep, r.ReasonError)
}

// ReorgExecutionResult is a struct that contains the information of the reorg execution
type ReorgExecutionResult struct {
	Request           ReorgRequest
	ExecutionCounter  uint64 // Number of reorg in this execution (is not a global unique ID!!!!)
	ExecutionError    error
	ExecutionTime     time.Time
	ExecutionDuration time.Duration
}

func (r ReorgExecutionResult) IsSuccess() bool {
	return r.ExecutionError == nil
}

func (r *ReorgExecutionResult) String() string {
	return fmt.Sprintf("Request: %s, ExecutionCounter: %d, ExecutionError: %v, ExecutionTime: %s, ExecutionDuration: %s",
		r.Request.String(), r.ExecutionCounter, r.ExecutionError, r.ExecutionTime.String(), r.ExecutionDuration.String())
}

type ReorgCallbackType = func(ReorgExecutionResult)

type StorageReorgInterface interface {
	ResetToL1BlockNumber(ctx context.Context, firstBlockNumberToKeep uint64, dbTx storageTxType) error
}

type ReorgState struct {
	mutex            sync.Mutex
	storage          StorageReorgInterface
	onReorgCallbacks []ReorgCallbackType
	lastReorgResult  *ReorgExecutionResult
}

func NewReorgState(storage StorageReorgInterface) *ReorgState {
	return &ReorgState{
		storage: storage,
	}
}

func (s *ReorgState) AddOnReorgCallback(f ReorgCallbackType) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onReorgCallbacks = append(s.onReorgCallbacks, f)
}

func (s *ReorgState) ExecuteReorg(ctx context.Context, reorgRequest ReorgRequest, dbTx storageTxType) ReorgExecutionResult {
	startTime := time.Now()
	err := s.storage.ResetToL1BlockNumber(ctx, reorgRequest.FirstL1BlockNumberToKeep, dbTx)
	res := s.createNewResult(reorgRequest, err, startTime)
	dbTx.AddCommitCallback(s.onTxCommit)
	dbTx.AddCommitCallback(s.onTxRollback)
	return res
}

func (s *ReorgState) onTxCommit(dbTx storageTxType, err error) {
	if err != nil {
		for _, f := range s.onReorgCallbacks {
			f(*s.lastReorgResult)
		}
	}
}

func (s *ReorgState) onTxRollback(dbTx storageTxType, err error) {
}

func (s *ReorgState) createNewResult(reorgRequest ReorgRequest, err error, startTime time.Time) ReorgExecutionResult {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	res := ReorgExecutionResult{
		Request:           reorgRequest,
		ExecutionCounter:  1,
		ExecutionError:    err,
		ExecutionTime:     startTime,
		ExecutionDuration: time.Since(startTime),
	}
	if s.lastReorgResult != nil {
		res.ExecutionCounter = s.lastReorgResult.ExecutionCounter + 1
	}
	s.lastReorgResult = &res
	return res
}
