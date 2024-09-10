package model_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	mock_model "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model/mocks"
	"github.com/stretchr/testify/require"
)

func TestOnSequencedBatchesOnL1HappyPath(t *testing.T) {
	mockStorage := mock_model.NewStorageVirtualBatchInterface(t)
	sut := model.NewBatchState(mockStorage)
	ctx := context.TODO()
	dbTx := mock_entities.NewTx(t)

	seq := model.SequenceOfBatches{
		Sequence: entities.SequencedBatches{
			FromBatchNumber: 2,
			ToBatchNumber:   2,
		},
		Batches: []*entities.VirtualBatch{},
	}
	previousSeq := entities.SequencedBatches{
		FromBatchNumber: 1,
		ToBatchNumber:   1,
	}
	seq.Batches = append(seq.Batches, &entities.VirtualBatch{})
	mockStorage.EXPECT().AddSequencedBatches(ctx, &seq.Sequence, dbTx).Return(nil)
	mockStorage.EXPECT().AddVirtualBatch(ctx, seq.Batches[0], dbTx).Return(nil)
	mockStorage.EXPECT().GetLatestSequence(ctx, dbTx).Return(&previousSeq, nil)
	err := sut.OnSequencedBatchesOnL1(ctx, seq, dbTx)
	require.NoError(t, err)
}

func TestOnSequencedBatchesOnFailStoreSeq(t *testing.T) {
	mockStorage := mock_model.NewStorageVirtualBatchInterface(t)
	sut := model.NewBatchState(mockStorage)
	ctx := context.TODO()
	dbTx := mock_entities.NewTx(t)

	seq := model.SequenceOfBatches{
		Sequence: entities.SequencedBatches{
			FromBatchNumber: 2,
			ToBatchNumber:   2,
		},
		Batches: []*entities.VirtualBatch{},
	}
	previousSeq := entities.SequencedBatches{
		FromBatchNumber: 1,
		ToBatchNumber:   1,
	}
	mockStorage.EXPECT().AddSequencedBatches(ctx, &seq.Sequence, dbTx).Return(entities.ErrAlreadyExists)
	mockStorage.EXPECT().GetSequenceByBatchNumber(ctx, seq.Sequence.FromBatchNumber, dbTx).Return(nil, entities.ErrNotFound)
	mockStorage.EXPECT().GetLatestSequence(ctx, dbTx).Return(&previousSeq, nil)
	err := sut.OnSequencedBatchesOnL1(ctx, seq, dbTx)
	require.Error(t, err)
}

func TestOnSequencedBatchesOnSameDataOnDB(t *testing.T) {
	mockStorage := mock_model.NewStorageVirtualBatchInterface(t)
	sut := model.NewBatchState(mockStorage)
	ctx := context.TODO()
	dbTx := mock_entities.NewTx(t)

	seq := model.SequenceOfBatches{
		Sequence: entities.SequencedBatches{
			FromBatchNumber: 1,
			ToBatchNumber:   1,
		},
		Batches: []*entities.VirtualBatch{},
	}

	mockStorage.EXPECT().AddSequencedBatches(ctx, &seq.Sequence, dbTx).Return(entities.ErrAlreadyExists)
	mockStorage.EXPECT().GetSequenceByBatchNumber(ctx, seq.Sequence.FromBatchNumber, dbTx).Return(&seq.Sequence, nil)
	mockStorage.EXPECT().GetLatestSequence(ctx, dbTx).Return(nil, nil)
	err := sut.OnSequencedBatchesOnL1(ctx, seq, dbTx)
	require.NoError(t, err)
}

func TestOnSequencedBatchesNotContigous(t *testing.T) {
	mockStorage := mock_model.NewStorageVirtualBatchInterface(t)
	sut := model.NewBatchState(mockStorage)
	ctx := context.TODO()
	dbTx := mock_entities.NewTx(t)

	seq := model.SequenceOfBatches{
		Sequence: entities.SequencedBatches{
			FromBatchNumber: 3,
			ToBatchNumber:   3,
		},
		Batches: []*entities.VirtualBatch{},
	}
	previousSeq := entities.SequencedBatches{
		FromBatchNumber: 1,
		ToBatchNumber:   1,
	}
	seq.Batches = append(seq.Batches, &entities.VirtualBatch{})
	mockStorage.EXPECT().GetLatestSequence(ctx, dbTx).Return(&previousSeq, nil)
	err := sut.OnSequencedBatchesOnL1(ctx, seq, dbTx)
	require.Error(t, err)
}
