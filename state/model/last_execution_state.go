package model

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type LastExecutionState struct {
	storage            storageKVInterface
	lastExecutionCache *lastExecutionBoundData
	currentExecution   *lastExecutionBoundData
}

type lastExecutionBoundData = entities.LastExecutionData

func NewLastExecutionState(storage storageKVInterface,
	configString string, syncVesion string) *LastExecutionState {
	res := &LastExecutionState{
		storage: storage,
	}
	last, err := res.getCurrentExecutionData(context.Background(), nil, storage)
	if err != nil {
		log.Error("Error getting current execution data", "error", err)
		return nil
	}
	res.lastExecutionCache = last
	res.currentExecution = &lastExecutionBoundData{
		Configuration: configString,
		SyncVersion:   syncVesion,
	}

	return res
}

func (l *LastExecutionState) GetPreviousExecutionData(ctx context.Context, dbTx storageTxType, storage storageKVInterface) (*lastExecutionBoundData, error) {
	return l.lastExecutionCache, nil
}

func (l *LastExecutionState) StartingSynchronization(ctx context.Context, dbTx storageTxType) error {
	l.currentExecution.LastStart = time.Now()
	return l.SetCurrentExecutionData(ctx, dbTx, l.storage, l.currentExecution)
}

func (l *LastExecutionState) SetCurrentExecutionData(ctx context.Context, dbTx storageTxType, storage storageKVInterface, data *lastExecutionBoundData) error {
	err := l.storage.KVSetJson(ctx, keyStorageLastExecution, data, nil, dbTx)
	if err != nil {
		return err
	}
	l.lastExecutionCache = data
	return nil
}

func (l *LastExecutionState) getCurrentExecutionData(ctx context.Context, dbTx storageTxType, storage storageKVInterface) (*lastExecutionBoundData, error) {
	return GetKVHelper[lastExecutionBoundData](ctx, keyStorageLastExecution, storage, dbTx)
}
