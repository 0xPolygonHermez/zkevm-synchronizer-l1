package pgstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/stretchr/testify/require"
)

func TestAddSequence(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	defer func() { _ = dbTx.Commit(ctx) }()
	err = storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, dbTx)
	require.NoError(t, err)

	seq := &pgstorage.SequencedBatches{
		FromBatchNumber: 100,
		ToBatchNumber:   300,
		L1BlockNumber:   123,
		ForkID:          123,
		Timestamp:       time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local),
		ReceivedAt:      time.Date(2024, 1, 14, 14, 30, 45, 0, time.Local),
		Source:          "test"}
	err = storage.AddSequencedBatches(ctx, seq, dbTx)
	require.NoError(t, err)

	seqDb, err := storage.GetSequenceByBatchNumber(ctx, seq.FromBatchNumber, dbTx)
	require.NoError(t, err)
	require.Equal(t, seq.FromBatchNumber, seqDb.FromBatchNumber)
	require.Equal(t, seq, seqDb)
}

func TestGetSequenceByBatchNumber(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	err := storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, nil)
	require.NoError(t, err)

	seq := &pgstorage.SequencedBatches{
		FromBatchNumber: 100,
		ToBatchNumber:   300,
		L1BlockNumber:   123,
		ForkID:          123,
		Timestamp:       time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local),
		ReceivedAt:      time.Date(2024, 1, 14, 14, 30, 45, 0, time.Local),
		Source:          "test"}
	err = storage.AddSequencedBatches(ctx, seq, nil)
	require.NoError(t, err)

	seqDb, err := storage.GetSequenceByBatchNumber(ctx, 300, nil)
	require.NoError(t, err)
	require.Equal(t, seq.ToBatchNumber, seqDb.ToBatchNumber)
	seqDb, err = storage.GetSequenceByBatchNumber(ctx, 100, nil)
	require.NoError(t, err)
	require.Equal(t, seq.ToBatchNumber, seqDb.ToBatchNumber)
	// If a seq doesnt exists returns nil, nil
	seqDb, err = storage.GetSequenceByBatchNumber(ctx, 301, nil)
	require.NoError(t, err)
	require.Nil(t, seqDb)
}

func TestGetSequencesGreatestOrEqualBatchNumber(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	err := storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, nil)
	require.NoError(t, err)

	// This generate:
	// 100-300, 301-400, 401-500
	populateSequences(t, ctx, 123, 3, storage, nil)

	seqDb, err := storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 1, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 3)

	// Must return 301-400 and 401-500
	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 325, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 2)
	require.Equal(t, uint64(301), (*seqDb)[0].FromBatchNumber)
	require.Equal(t, uint64(401), (*seqDb)[1].FromBatchNumber)

	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 301, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 2)

	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 400, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 2)

	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 401, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 1)
	// Returns 0 rows so it returns nil, nil
	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 501, nil)
	require.NoError(t, err)
	require.Nil(t, seqDb)
}

func TestDeleteSequencesGreatestOrEqualBatchNumber(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	storage := initDbForTest(t)
	ctx := context.TODO()
	err := storage.AddBlock(ctx, &pgstorage.L1Block{BlockNumber: 123}, nil)
	require.NoError(t, err)

	// This generate:
	// 100-300, 301-400, 401-500
	populateSequences(t, ctx, 123, 3, storage, nil)

	seqDb, err := storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 1, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 3)

	// Returns 0 rows so it returns nil, nil
	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 501, nil)
	require.NoError(t, err)
	require.Nil(t, seqDb)
	// This must delete: 301-400, 401-500
	err = storage.DeleteSequencesGreatestOrEqualBatchNumber(ctx, 325, nil)
	require.NoError(t, err)

	seqDb, err = storage.GetSequencesGreatestOrEqualBatchNumber(ctx, 1, nil)
	require.NoError(t, err)
	require.NotNil(t, seqDb)
	require.Len(t, *seqDb, 1)
	require.Equal(t, uint64(100), (*seqDb)[0].FromBatchNumber)
}

func populateSequences(t *testing.T, ctx context.Context, l1BlockNumber uint64, num_seq int, storage *pgstorage.PostgresStorage, dbTx entities.Tx) {
	if num_seq == 0 {
		return
	}
	seq := &pgstorage.SequencedBatches{
		FromBatchNumber: 100,
		ToBatchNumber:   300,
		L1BlockNumber:   l1BlockNumber,
		ForkID:          123,
		Timestamp:       time.Date(2023, 12, 14, 14, 30, 45, 0, time.Local),
		ReceivedAt:      time.Date(2024, 1, 14, 14, 30, 45, 0, time.Local),
		Source:          "test"}
	err := storage.AddSequencedBatches(ctx, seq, dbTx)
	require.NoError(t, err)
	for i := 1; i < num_seq; i++ {
		seq.FromBatchNumber = seq.ToBatchNumber + 1
		seq.ToBatchNumber = seq.FromBatchNumber + 99
		err = storage.AddSequencedBatches(ctx, seq, dbTx)
		require.NoError(t, err)
	}

}
