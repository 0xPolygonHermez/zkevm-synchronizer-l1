package sqlstorage_test

import (
	"context"
	"testing"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/stretchr/testify/require"
)

const (
	testKey = "fake_key_used_for_unittest"
)

func TestKVSet(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	exists, err := storage.KVExists(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.False(t, exists)
	err = storage.KVSetString(ctx, testKey, "fake_value", nil, dbTx)
	require.NoError(t, err)
	exists, err = storage.KVExists(ctx, testKey, dbTx)
	require.NoError(t, err)
	require.True(t, exists)
	value, err := storage.KVGetString(ctx, testKey, nil, dbTx)
	require.NoError(t, err)
	require.Equal(t, "fake_value", value)
	_ = dbTx.Commit(ctx)
}

type kvTestStruct struct {
	A int
	B string
}

func TestKVJson(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	data := kvTestStruct{A: 1, B: "test"}
	err = storage.KVSetJson(ctx, testKey, data, nil, dbTx)
	require.NoError(t, err)
	var dataRead kvTestStruct
	err = storage.KVGetJson(ctx, testKey, &dataRead, nil, dbTx)
	require.NoError(t, err)
	require.Equal(t, data, dataRead)

	_ = dbTx.Commit(ctx)
}

func TestKVUint64(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	data := uint64(1234)
	err = storage.KVSetJson(ctx, testKey, data, nil, dbTx)
	require.NoError(t, err)
	var dataRead uint64
	err = storage.KVGetJson(ctx, testKey, &dataRead, nil, dbTx)
	require.NoError(t, err)
	require.Equal(t, data, dataRead)

	err = storage.KVSetString(ctx, testKey, "not a number", nil, dbTx)
	require.NoError(t, err)
	err = storage.KVGetJson(ctx, testKey, &dataRead, nil, dbTx)
	require.Error(t, err)

	_ = dbTx.Commit(ctx)
}

func TestKVSetMetadataDefault(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.KVSetString(ctx, testKey, "fake_value", nil, dbTx)
	require.NoError(t, err)
	timeNow := time.Now()
	metadata := sqlstorage.KVMetadataEntry{}
	_, err = storage.KVGetString(ctx, testKey, &metadata, dbTx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, metadata.CreatedAt.Unix(), timeNow.Unix())
	require.GreaterOrEqual(t, metadata.CreatedAt.Unix(), timeNow.Unix())
	require.Equal(t, metadata.CreatedAt, metadata.UpdatedAt)
	require.Equal(t, zkevm_synchronizer_l1.Version, metadata.SyncVersion)
}

func TestKVSetMetadataForced(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	timeNow := time.Now().Round(time.Second)
	metadata := sqlstorage.KVMetadataEntry{
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
		SyncVersion: "forced_version",
	}
	err = storage.KVSetString(ctx, testKey, "fake_value", &metadata, dbTx)
	require.NoError(t, err)
	readMetadata := sqlstorage.KVMetadataEntry{}
	_, err = storage.KVGetString(ctx, testKey, &readMetadata, dbTx)
	require.NoError(t, err)
	require.Equal(t, metadata.SyncVersion, readMetadata.SyncVersion)
	// The DB store time.Time with no timezone to UTC+0
	require.Equal(t, metadata.CreatedAt.UTC(), readMetadata.CreatedAt)
	require.Equal(t, metadata.UpdatedAt.UTC(), readMetadata.UpdatedAt)

	// Update Metadata:
	timeUpdate := timeNow.Add(time.Hour * 4)
	metadata = sqlstorage.KVMetadataEntry{
		CreatedAt:   timeUpdate,
		UpdatedAt:   timeUpdate,
		SyncVersion: "another_vesion",
	}
	err = storage.KVSetString(ctx, testKey, "fake_value2", &metadata, dbTx)
	require.NoError(t, err)
	_, err = storage.KVGetString(ctx, testKey, &readMetadata, dbTx)
	require.NoError(t, err)
	// CreatedAt should be ignored
	metadata.CreatedAt = timeNow
	require.Equal(t, metadata.SyncVersion, readMetadata.SyncVersion)
	require.Equal(t, metadata.CreatedAt.UTC(), readMetadata.CreatedAt)
	require.Equal(t, metadata.UpdatedAt.UTC(), readMetadata.UpdatedAt)

}
