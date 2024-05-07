package model

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type VirtualBatch = entities.VirtualBatch
type SequencedBatches = entities.SequencedBatches
type dbTxType = entities.Tx

type SequenceOfBatches struct {
	Sequence SequencedBatches
	Batches  []*VirtualBatch
}

type StorageVirtualBatchInterface interface {
	AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx dbTxType) error
	GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*VirtualBatch, error)
	AddSequencedBatches(ctx context.Context, sequence *SequencedBatches, dbTx dbTxType) error
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatches, error)
}

type BatchState struct {
	store StorageVirtualBatchInterface
}

func NewBatchState(store StorageVirtualBatchInterface) *BatchState {
	return &BatchState{
		store: store,
	}
}

// OnSequencedBatchesOnL1 a new sequenceBatch call have been found on L1, add to local database
func (b *BatchState) OnSequencedBatchesOnL1(ctx context.Context, seq SequenceOfBatches, dbTx dbTxType) error {

	err := SetStorageHelper[*SequencedBatches](ctx, &seq.Sequence, dbTx,
		b.store.AddSequencedBatches, b.store.GetSequenceByBatchNumber)
	if err != nil {
		return err
	}
	for _, batch := range seq.Batches {
		err := SetStorageHelper[*VirtualBatch](ctx, batch, dbTx,
			b.store.AddVirtualBatch, b.store.GetVirtualBatchByBatchNumber)
		if err != nil {
			return err
		}
	}
	return nil
}
