package l1event_orders

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
)

type Sequence struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
}

// GetSequenceFromL1EventOrder returns the sequence of batches of  given event
// There are event that are Batch based or not, if not it returns a nil
func GetSequenceFromL1EventOrder(event etherman.EventOrder, l1Block *etherman.Block, position int) *Sequence {
	switch event {
	case etherman.SequenceBatchesOrder:
		return getSequence(l1Block.SequencedBatches[position],
			func(batch etherman.SequencedBatch) uint64 { return batch.BatchNumber })

	}
	return nil
}

func getSequence[T any](batches []T, getBatchNumber func(T) uint64) *Sequence {
	if len(batches) == 0 {
		return nil
	}
	res := Sequence{FromBatchNumber: getBatchNumber(batches[0]),
		ToBatchNumber: getBatchNumber(batches[0])}
	for _, batch := range batches {
		if getBatchNumber(batch) < res.FromBatchNumber {
			res.FromBatchNumber = getBatchNumber(batch)
		}
		if getBatchNumber(batch) > res.ToBatchNumber {
			res.ToBatchNumber = getBatchNumber(batch)
		}
	}
	return &res
}
