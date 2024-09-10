package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestReset(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	blocks := []sqlstorage.L1Block{block300, block301, block310}
	for _, block := range blocks {
		err = storage.AddBlock(ctx, &block, dbTx)
		require.NoError(t, err)
	}
	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	lastBlock, err := storage.GetLastBlock(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, block310.BlockNumber, lastBlock.BlockNumber)

	dbTx, err = storage.BeginTransaction(ctx)
	require.NoError(t, err)
	err = storage.ResetToL1BlockNumber(ctx, 300, dbTx)
	require.NoError(t, err)
	err = dbTx.Commit(ctx)
	require.NoError(t, err)

	lastBlock, err = storage.GetLastBlock(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, block300.BlockNumber, lastBlock.BlockNumber)
}
