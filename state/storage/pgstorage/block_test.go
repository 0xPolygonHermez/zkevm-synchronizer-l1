package pgstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestBlockAddAndGets(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block300 := pgstorage.L1Block{
		BlockNumber: 300,
		BlockHash:   common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		ParentHash:  common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678906789"),
		ReceivedAt:  time.Now().Truncate(time.Second),
		Checked:     false,
		SyncVersion: "test",
	}
	block301 := pgstorage.L1Block{
		BlockNumber: 301,
		BlockHash:   common.HexToHash("0x00"),
		ParentHash:  common.HexToHash("0x12"),
		ReceivedAt:  time.Now().Truncate(time.Second),
		Checked:     true,
		SyncVersion: "test_2",
	}
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

		prevBlock, err := storage.GetPreviousBlock(ctx, 1, nil, dbTx)
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
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block300 := pgstorage.L1Block{
		BlockNumber: 300,
		BlockHash:   common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678901234"),
		ParentHash:  common.HexToHash("0x1234567890123456789012345678901234567890123456789012345678906789"),
		ReceivedAt:  time.Now().Truncate(time.Second),
		Checked:     false,
		SyncVersion: "test",
	}
	block301 := pgstorage.L1Block{
		BlockNumber: 301,
		BlockHash:   common.HexToHash("0x00"),
		ParentHash:  common.HexToHash("0x12"),
		ReceivedAt:  time.Now().Truncate(time.Second),
		Checked:     true,
		SyncVersion: "test_2",
	}
	block310 := pgstorage.L1Block{
		BlockNumber: 310,
		BlockHash:   common.HexToHash("0x00"),
		ParentHash:  common.HexToHash("0x12"),
		ReceivedAt:  time.Now().Truncate(time.Second),
		Checked:     true,
		SyncVersion: "test_2",
	}
	err = storage.AddBlock(ctx, &block300, dbTx)
	require.NoError(t, err)
	err = storage.AddBlock(ctx, &block301, dbTx)
	require.NoError(t, err)
	err = storage.AddBlock(ctx, &block310, dbTx)
	require.NoError(t, err)

	block, err := storage.GetPreviousBlock(ctx, 0, nil, dbTx)
	require.NoError(t, err)
	require.Equal(t, block310.String(), block.String(), "offset 0 must return latest block")
	blockNumber := uint64(310)
	block, err = storage.GetPreviousBlock(ctx, 0, &blockNumber, dbTx)
	require.NoError(t, err)
	require.Equal(t, block310.String(), block.String(), "offset 0 and fromBlock latest must return latest block")

	blockNumber = uint64(309)
	block, err = storage.GetPreviousBlock(ctx, 0, &blockNumber, dbTx)
	require.NoError(t, err)
	require.Equal(t, block301.String(), block.String(), "offset 0 and fromBlock latest-1 must return previous to latest")

	blockNumber = uint64(309)
	block, err = storage.GetPreviousBlock(ctx, 1, &blockNumber, dbTx)
	require.NoError(t, err)
	require.Equal(t, block300.String(), block.String(), "offset 1 and fromBlock latest-1 must return 2 before to latest")

}
