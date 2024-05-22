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

// GetKVHelper is a helper function to get a value from the KV storage, if not found returns nil
func GetKVHelper[T any](ctx context.Context, key KVKey, getter StorageKVInterface, dbTx storageTxType) (*T, error) {
	var value T
	err := getter.KVGetJson(ctx, key, &value, nil, dbTx)
	if err != nil && errors.Is(err, entities.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &value, nil
}
