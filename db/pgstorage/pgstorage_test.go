package pgstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestQueries(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	storage, err := pgstorage.NewPostgresStorage(getStorageConfig())
	require.NoError(t, err)
	dbTx, err := storage.Begin(context.Background())
	require.NoError(t, err)
	defer func() { _ = dbTx.Rollback(context.Background()) }()
	err = storage.AddBlock(context.Background(), &pgstorage.L1Block{BlockNumber: 300}, dbTx)
	require.NoError(t, err)
	sb := pgstorage.SequencedBatches{
		FromBatchNumber: 100,
		ToBatchNumber:   200,
		L1BlockNumber:   300,
		Timestamp:       time.Now(),
		L1InfoRoot:      common.Hash{},
	}
	err = storage.AddSequencedBatches(context.Background(), sb, dbTx)
	require.NoError(t, err)
	data, err := storage.GetSequenceByBatchNumber(context.Background(), uint64(102), dbTx)
	require.NoError(t, err)
	require.Equal(t, sb.FromBatchNumber, data.FromBatchNumber)
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
