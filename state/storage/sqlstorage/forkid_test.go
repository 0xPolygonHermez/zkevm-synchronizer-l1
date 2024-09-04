package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestAddForkID(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	forkInterval := sqlstorage.ForkIDInterval{
		FromBatchNumber: 1,
		ToBatchNumber:   2,
		ForkId:          1,
		Version:         "test",
		BlockNumber:     1,
	}

	err = storage.AddForkID(ctx, forkInterval, dbTx)
	require.NoError(t, err)
	readForkIntervals, err := storage.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Len(t, readForkIntervals, 1)
	require.Equal(t, forkInterval, readForkIntervals[0])
}

func TestAddForkIDOnConlict(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	forkInterval := sqlstorage.ForkIDInterval{
		FromBatchNumber: 1,
		ToBatchNumber:   2,
		ForkId:          1,
		Version:         "test",
		BlockNumber:     1,
	}

	err = storage.AddForkID(ctx, forkInterval, dbTx)
	require.NoError(t, err)
	// Now it's a conflict, must update blockNumber
	forkInterval.BlockNumber = 2
	err = storage.AddForkID(ctx, forkInterval, dbTx)
	require.NoError(t, err)
	readForkIntervals, err := storage.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Len(t, readForkIntervals, 1)
	require.Equal(t, forkInterval, readForkIntervals[0])
}

func TestUpdateForkID(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	forkInterval := sqlstorage.ForkIDInterval{
		FromBatchNumber: 1,
		ToBatchNumber:   2,
		ForkId:          1,
		Version:         "test",
		BlockNumber:     1,
	}
	err = storage.AddForkID(ctx, forkInterval, dbTx)
	require.NoError(t, err)
	forkInterval.ToBatchNumber = 3
	storage.UpdateForkID(ctx, forkInterval, dbTx)
	readForkIntervals, err := storage.GetForkIDs(ctx, dbTx)
	require.NoError(t, err)
	require.Len(t, readForkIntervals, 1)
	require.Equal(t, forkInterval, readForkIntervals[0])
}
