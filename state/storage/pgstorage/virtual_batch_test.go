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

func TestAddVirtualBatch(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, dbTx)
	require.NoError(t, err)

	err = storage.AddSequencedBatches(ctx, &pgstorage.SequencedBatches{FromBatchNumber: 100, ToBatchNumber: 300, L1BlockNumber: 123}, dbTx)
	require.NoError(t, err)
	l1InfoRoot := common.HexToHash("0xad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5")
	extraInfo := "test batch"
	virtualBatch := pgstorage.VirtualBatch{BatchNumber: 300, BlockNumber: 123, SequenceFromBatchNumber: 100,
		ReceivedAt:  time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local),
		L1InfoRoot:  &l1InfoRoot,
		BatchL2Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06},
		ExtraInfo:   &extraInfo,
	}
	err = storage.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.NoError(t, err)
}

func TestAddVirtualBatchDuplicated(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, dbTx)
	require.NoError(t, err)

	err = storage.AddSequencedBatches(ctx, &pgstorage.SequencedBatches{FromBatchNumber: 100, ToBatchNumber: 300, L1BlockNumber: 123}, dbTx)
	require.NoError(t, err)
	virtualBatch := pgstorage.VirtualBatch{BatchNumber: 300, BlockNumber: 123, SequenceFromBatchNumber: 100}
	err = storage.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.NoError(t, err)
	err = storage.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.ErrorIs(t, err, entities.ErrAlreadyExists)
}

func TestAddVirtualBatchMissingSequence(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, dbTx)
	require.NoError(t, err)

	virtualBatch := pgstorage.VirtualBatch{BatchNumber: 300, BlockNumber: 123, SequenceFromBatchNumber: 100}
	err = storage.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.ErrorIs(t, err, entities.ErrForeignKeyViolation)

}

func TestGetVirtualBatch(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123,
		ReceivedAt: time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local)}, dbTx)
	require.NoError(t, err)

	err = storage.AddSequencedBatches(ctx, &pgstorage.SequencedBatches{FromBatchNumber: 100, ToBatchNumber: 300, L1BlockNumber: 123}, dbTx)
	require.NoError(t, err)
	virtualBatch := pgstorage.VirtualBatch{BatchNumber: 300, BlockNumber: 123, SequenceFromBatchNumber: 100, ReceivedAt: time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local)}
	err = storage.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	require.NoError(t, err)

	dbVirtualBatch, err := storage.GetVirtualBatchByBatchNumber(ctx, 300, dbTx)
	require.NoError(t, err)
	require.Equal(t, virtualBatch, *dbVirtualBatch)
}

func TestGetVirtualBatchNotFound(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()

	_, err = storage.GetVirtualBatchByBatchNumber(ctx, 300, dbTx)
	require.ErrorIs(t, err, entities.ErrNotFound)
}
