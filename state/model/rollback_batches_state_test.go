package model_test

import (
	"context"
	"testing"
	"time"

	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	mock_model "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRollbackBatchesHappyPath(t *testing.T) {
	storageMock := mock_model.NewStorageRollbackBatchesInterface(t)
	sut := model.NewRollbackBatchesState(storageMock)
	request := model.RollbackBatchesRequest{
		LastBatchNumber:       123,
		LastBatchAccInputHash: [32]byte{},
		L1BlockNumber:         60,
		L1BlockTimestamp:      time.Now(),
	}
	seqs := &model.SequencesBatchesSlice{
		model.SequencedBatches{
			FromBatchNumber: 1,
			ToBatchNumber:   123,
			L1BlockNumber:   50,
		},
	}
	dbTxMock := mock_entities.NewTx(t)
	storageMock.EXPECT().GetSequencesGreatestOrEqualBatchNumber(context.TODO(), request.LastBatchNumber+1, dbTxMock).Return(seqs, nil)
	storageMock.EXPECT().AddRollbackBatchesLogEntry(context.TODO(), mock.Anything, dbTxMock).Return(nil)
	storageMock.EXPECT().DeleteSequencesGreatestOrEqualBatchNumber(context.TODO(), request.LastBatchNumber+1, dbTxMock).Return(nil)
	dbTxMock.EXPECT().AddCommitCallback(mock.Anything)
	res, err := sut.ExecuteRollbackBatches(context.TODO(), request, dbTxMock)
	require.NoError(t, err)
	require.Equal(t, request.LastBatchNumber, res.RollbackEntry.LastBatchNumber)
	require.Equal(t, uint64(50), res.RollbackEntry.UndoFirstBlockNumber)
}
