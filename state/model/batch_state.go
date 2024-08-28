package model

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type VirtualBatch = entities.VirtualBatch
type SequencedBatches = entities.SequencedBatches
type SequencesBatchesSlice = entities.SequencesBatchesSlice
type RollbackBatchesLogEntry = entities.RollbackBatchesLogEntry
type dbTxType = entities.Tx

type SequenceOfBatches struct {
	Sequence SequencedBatches
	Batches  []*VirtualBatch
}

func (s *SequenceOfBatches) String() string {
	return fmt.Sprintf("SequenceOfBatches{Sequence: %s, Batches: %v}", s.Sequence.String(), s.Batches)
}

type StorageVirtualBatchInterface interface {
	AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx dbTxType) error
	GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*VirtualBatch, error)
	AddSequencedBatches(ctx context.Context, sequence *SequencedBatches, dbTx dbTxType) error
	GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*SequencedBatches, error)
	GetLatestSequence(ctx context.Context, dbTx storageTxType) (*SequencedBatches, error)
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
	err := b.sanityCheck(ctx, seq, dbTx)
	if err != nil {
		return err
	}
	err = SetStorageHelper[*SequencedBatches](ctx, &seq.Sequence, dbTx,
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

func (b *BatchState) sanityCheck(ctx context.Context, seq SequenceOfBatches, dbTx dbTxType) error {
	lastSeq, err := b.store.GetLatestSequence(ctx, dbTx)
	if err != nil {
		return err
	}
	if lastSeq == nil {
		// Is the first sequence
		if seq.Sequence.FromBatchNumber != 1 {
			log.Warnf("First sequence batch number is not 1: %v", seq.Sequence.String())
		}
		return nil
	}

	if lastSeq.ToBatchNumber+1 != seq.Sequence.FromBatchNumber {
		return fmt.Errorf("sequence batch is not contiguous: last on DB %v, new sequence %v", lastSeq, seq.Sequence.String())
	}
	return nil
}
