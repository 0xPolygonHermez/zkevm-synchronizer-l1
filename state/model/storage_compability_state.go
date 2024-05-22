package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage"
)

// storageContentsBoundData is a struct that contains the RollupID and L1ChainID
// basically is the data that is coupled to the rest of the contents of storage
// If any of this values changes from an execution to another means that all data
// in storage is invalid

type storageContentsBoundData = entities.StorageContentsBoundData
type storageKVInterface = storage.KvStorer

type StorageCompatibilityState struct {
	storage storageKVInterface
}

func NewStorageCompatibilityState(storage storageKVInterface) *StorageCompatibilityState {
	return &StorageCompatibilityState{
		storage: storage,
	}
}

func (s *StorageCompatibilityState) GetStorageContentsBoundData(ctx context.Context, dbTx storageTxType) (*storageContentsBoundData, error) {
	storageData, err := GetKVHelper[storageContentsBoundData](ctx, keyStorageContentBound, s.storage, dbTx)
	if err != nil {
		log.Error("Error getting sanity storage data", "err", err)
		return nil, err
	}
	return storageData, nil

}

func (s *StorageCompatibilityState) CheckAndUpdateStorage(ctx context.Context, runBoundData storageContentsBoundData, overrideStorageCheck bool, dbTx storageTxType) error {
	err := s.CheckSanity(ctx, runBoundData, dbTx)
	if err != nil && overrideStorageCheck && errors.Is(err, entities.ErrMismatchStorageAndExecutionEnvironment) {
		log.Warnf("mismatch checking storage sanity, but the error is bypass by configuration (OverrideStorageCheck) . Error: %v", err)
		return nil
	}
	if err != nil {
		log.Errorf("Error checking storage sanity. Error: %v", err)
		return err
	}
	log.Infof("Storage sanity check passed, storing data on DB: %s", runBoundData.String())
	err = s.setStorageContentsBoundData(ctx, runBoundData, dbTx)
	if err != nil {
		log.Errorf("Error saving sanity data to storage. Error: %v", err)
		return err
	}
	return nil
}

func (s *StorageCompatibilityState) CheckSanity(ctx context.Context, runBoundData storageContentsBoundData, dbTx storageTxType) error {
	storageData, err := s.GetStorageContentsBoundData(ctx, dbTx)
	if err != nil {
		return err
	}
	return s.compareData(runBoundData, storageData)
}

func (s *StorageCompatibilityState) setStorageContentsBoundData(ctx context.Context, data storageContentsBoundData, dbTx storageTxType) error {
	return s.storage.KVSetJson(ctx, keyStorageContentBound, data, nil, dbTx)
}

func (s *StorageCompatibilityState) compareData(runBoundData storageContentsBoundData, storageData *storageContentsBoundData) error {
	if storageData == nil {
		log.Debug("No Sanity storage data in storage")
		return nil
	}
	if runBoundData != *storageData {
		err := fmt.Errorf("sanity storage data mismatch: run: %s, storage: %s  err:%w", runBoundData.String(), storageData.String(), entities.ErrMismatchStorageAndExecutionEnvironment)
		log.Warnf(err.Error())
		return err
	}
	return nil
}
