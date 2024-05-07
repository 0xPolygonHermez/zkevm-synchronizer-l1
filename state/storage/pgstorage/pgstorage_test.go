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

func initDbForTest(t *testing.T) *pgstorage.PostgresStorage {
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	return storage
}

func TestQueries(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	dbConfig := getStorageConfig()
	err := pgstorage.ResetDB(dbConfig)
	require.NoError(t, err)
	storage, err := pgstorage.NewPostgresStorage(dbConfig)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(context.Background())
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(context.Background()) }()
	err = storage.AddBlock(context.Background(), &pgstorage.L1Block{BlockNumber: 300}, dbTx)
	require.NoError(t, err)
	sb := pgstorage.SequencedBatches{
		FromBatchNumber: 100,
		ToBatchNumber:   200,
		L1BlockNumber:   300,
		Timestamp:       time.Now(),
		L1InfoRoot:      common.Hash{},
	}
	err = storage.AddSequencedBatches(context.Background(), &sb, dbTx)
	require.NoError(t, err)
	data, err := storage.GetSequenceByBatchNumber(context.Background(), uint64(102), dbTx)
	require.NoError(t, err)
	require.Equal(t, sb.FromBatchNumber, data.FromBatchNumber)
}

func TestBlockAddAndGets(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
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

		prevBlock, err := storage.GetPreviousBlock(ctx, 1, dbTx)
		if testCase.queryPreviousBlock == nil {
			require.ErrorIs(t, err, entities.ErrNotFound)
		} else {
			require.NoError(t, err)
			require.Equal(t, testCase.queryPreviousBlock.String(), prevBlock.String(), "GetPreviousBlock")
		}

	}
}

func getStorageConfig() pgstorage.Config {
	return pgstorage.Config{
		Host:     "localhost",
		Port:     "5436",
		User:     "test_user",
		Password: "test_password",
		Name:     "sync",
		MaxConns: 10,
	}
}
