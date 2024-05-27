package pgstorage_test

import (
	"context"
	"testing"
	"time"

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
	skipDatabaseTestIfNeeded(t)
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

func skipDatabaseTestIfNeeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
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
