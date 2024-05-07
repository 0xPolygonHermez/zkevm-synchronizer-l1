package pgstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/stretchr/testify/require"
)

func TestAddSequence(t *testing.T) {
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
