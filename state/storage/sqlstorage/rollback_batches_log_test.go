package sqlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestAddRollbackBatches(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	err = storage.AddBlock(ctx, &sqlstorage.L1Block{BlockNumber: 123}, dbTx)
	require.NoError(t, err)

	entry := &sqlstorage.RollbackBatchesLogEntry{
		BlockNumber:           123,
		LastBatchNumber:       123,
		LastBatchAccInputHash: common.HexToHash("0x123"),
		L1EventAt:             time.Now().Truncate(time.Second).UTC(),
		ReceivedAt:            time.Now().Truncate(time.Second).UTC(),
		UndoFirstBlockNumber:  123,
		Description:           "this is a unittest entry",
		SequencesDeleted: []entities.SequencedBatches{
			{
				FromBatchNumber: 100,
				ToBatchNumber:   300,
			},
			{
				FromBatchNumber: 200,
				ToBatchNumber:   400,
			},
		},
	}
	entry.SetSyncVersion("version-unittest")
	err = storage.AddRollbackBatchesLogEntry(ctx, entry, dbTx)
	require.NoError(t, err)

	entries, err := storage.GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx, 123, dbTx)
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, entry, &entries[0])
}
