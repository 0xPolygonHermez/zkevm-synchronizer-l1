package model

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage"
)

type KVKey = string

type StorageKVInterface = storage.KvStorer

const (
	// keyStorageContentBound is the name of the KEY for sanity storage
	keyStorageContentBound KVKey = "ContentBound"
)

type StateKVInterface interface {
	SetKV(ctx context.Context, key KVKey, value interface{}, dbTx storageTxType) error
	GetKV(ctx context.Context, key KVKey, value interface{}, dbTx storageTxType) error
	// In the future this class is going to check version of lib
	// and call the rest when updating
}

type KVState struct {
	storage StorageKVInterface
}

func NewKVState(storage StorageKVInterface) *KVState {
	return &KVState{
		storage: storage,
	}
}

func (s *KVState) SetKV(ctx context.Context, key KVKey, value interface{}, dbTx storageTxType) error {
	return s.storage.KVSetJson(ctx, key, value, nil, dbTx)
}

func (s *KVState) GetKV(ctx context.Context, key KVKey, value interface{}, dbTx storageTxType) error {
	err := s.storage.KVGetJson(ctx, key, value, nil, dbTx)

	return err
}

// GetKVHelper is a helper function to get a value from the KV storage, if not found returns nil
func GetKVHelper[T any](ctx context.Context, key KVKey, getter StateKVInterface, dbTx storageTxType) (*T, error) {
	var value T
	err := getter.GetKV(ctx, key, &value, dbTx)
	if err != nil && errors.Is(err, entities.ErrNotFound) {
		return nil, nil
	}
	return &value, err
}
