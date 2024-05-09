package pgstorage_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/stretchr/testify/require"
)

const (
	testKey = "fake_key_used_for_unittest"
)

func TestKVSet(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	exists, err := storage.KVExists(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.False(t, exists)
	err = storage.KVSetString(ctx, testKey, "fake_value", dbTx)
	require.NoError(t, err)
	exists, err = storage.KVExists(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.True(t, exists)
	value, err := storage.KVGetString(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.Equal(t, "fake_value", value)
	dbTx.Commit(ctx)
}

type kvTestStruct struct {
	A int
	B string
}

func TestKVJson(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	data := kvTestStruct{A: 1, B: "test"}
	err = storage.KVSetJson(ctx, testKey, data, dbTx)
	require.NoError(t, err)
	var dataRead kvTestStruct
	err = storage.KVGetJson(ctx, testKey, &dataRead, dbTx)
	require.NoError(t, err)
	require.Equal(t, data, dataRead)

	dbTx.Commit(ctx)
}

func TestKVUint64(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	data := uint64(1234)
	err = storage.KVSetUint64(ctx, testKey, data, dbTx)
	require.NoError(t, err)
	dataRead, err := storage.KVGetUint64(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.Equal(t, data, dataRead)

	err = storage.KVSetString(ctx, testKey, "not a number", dbTx)
	require.NoError(t, err)
	_, err = storage.KVGetUint64(ctx, testKey, dbTx)
	require.Error(t, err)

	dbTx.Commit(ctx)
}
