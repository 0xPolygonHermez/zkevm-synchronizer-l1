package sqlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	block300 = sqlstorage.L1Block{
		BlockNumber: 300,
		BlockHash:   common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		ParentHash:  common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678906789"),
		ReceivedAt:  time.Now().Truncate(time.Second).UTC(),
		Checked:     false,
		SyncVersion: "test",
	}
	block301 = sqlstorage.L1Block{
		BlockNumber: 301,
		BlockHash:   common.HexToHash("0x00"),
		ParentHash:  common.HexToHash("0x12"),
		ReceivedAt:  time.Now().Truncate(time.Second).UTC(),
		Checked:     true,
		SyncVersion: "test_2",
	}
	block310 = sqlstorage.L1Block{
		BlockNumber: 310,
		BlockHash:   common.HexToHash("0x00"),
		ParentHash:  common.HexToHash("0x12"),
		ReceivedAt:  time.Now().Truncate(time.Second).UTC(),
		Checked:     true,
		SyncVersion: "test_2",
	}
)

func TestBlockSetTimeToUTC(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block310.ReceivedAt = time.Now().Truncate(time.Second)
	err = storage.AddBlock(ctx, &block310, dbTx)
	require.NoError(t, err)
	dbBlock, err := storage.GetBlockByNumber(ctx, block310.BlockNumber, dbTx)
	require.NoError(t, err)
	require.Equal(t, block310.ReceivedAt.UTC(), dbBlock.ReceivedAt, "GetBlockByNumber tstamp is UTC")
}

func TestBlockAddAndGets(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	tests := []struct {
		addBlock           *pgstorage.L1Block
		queryLastBlock     *pgstorage.L1Block
		queryPreviousBlock *pgstorage.L1Block
	}{
		{&block300, &block300, nil},
		{&block301, &block301, &block300},
	}
	for _, testCase := range tests {
		err = storage.AddBlock(ctx, testCase.addBlock, dbTx)
		require.NoError(t, err)
		dbBlock, err := storage.GetBlockByNumber(ctx, testCase.addBlock.BlockNumber, dbTx)
		require.NoError(t, err)
		require.Equal(t, testCase.addBlock.String(), dbBlock.String(), "GetBlockByNumber")

		lastBlock, err := storage.GetLastBlock(ctx, dbTx)
		require.NoError(t, err)
		require.Equal(t, testCase.queryLastBlock.String(), lastBlock.String(), "GetLastBlock")

		prevBlock, err := storage.GetPreviousBlock(ctx, 1, dbTx)
		if testCase.queryPreviousBlock == nil {
			require.ErrorIs(t, err, entities.ErrNotFound)
		} else {
			require.NoError(t, err)
			require.Equal(t, testCase.queryPreviousBlock.String(), prevBlock.String(), "GetPreviousBlock")
		}
	}
}

func TestGetPreviousBlockFromBlock(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()

	err = storage.AddBlock(ctx, &block300, dbTx)
	require.NoError(t, err)
	err = storage.AddBlock(ctx, &block301, dbTx)
	require.NoError(t, err)
	err = storage.AddBlock(ctx, &block310, dbTx)
	require.NoError(t, err)

	block, err := storage.GetPreviousBlock(ctx, 0, dbTx)
	require.NoError(t, err)
	require.Equal(t, block310.String(), block.String(), "offset 0 must return latest block")
}

func TestUpdateCheckedBlockByNumber(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block := block300
	block.Checked = true
	err = storage.AddBlock(ctx, &block, dbTx)
	require.NoError(t, err)
	err = storage.UpdateCheckedBlockByNumber(ctx, 300, false, dbTx)
	require.NoError(t, err)
	blockRead, err := storage.GetBlockByNumber(ctx, 300, dbTx)
	require.NoError(t, err)
	require.False(t, blockRead.Checked)
	err = storage.UpdateCheckedBlockByNumber(ctx, 300, true, dbTx)
	require.NoError(t, err)
	blockRead, err = storage.GetBlockByNumber(ctx, 300, dbTx)
	require.NoError(t, err)
	require.True(t, blockRead.Checked)
	err = storage.UpdateCheckedBlockByNumber(ctx, 300, false, dbTx)
	require.NoError(t, err)

	blockRead, err = storage.GetFirstUncheckedBlock(ctx, 0, dbTx)
	require.NoError(t, err)
	require.Equal(t, block300.String(), blockRead.String())
}
