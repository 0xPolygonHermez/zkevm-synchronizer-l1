package entities

import (
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
)

type SequencesBatchesSlice []SequencedBatches

type SequencedBatches struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	L1BlockNumber   uint64
	ForkID          uint64
	Timestamp       time.Time
	ReceivedAt      time.Time
	L1InfoRoot      common.Hash
	Source          string
}

func (s *SequencedBatches) String() string {
	return fmt.Sprintf("SequencedBatches{FromBatchNumber: %d, ToBatchNumber: %d, L1BlockNumber: %d, ForkID: %d, Timestamp: %s, ReceivedAt: %s, L1InfoRoot: %s, Source: %s}",
		s.FromBatchNumber, s.ToBatchNumber, s.L1BlockNumber, s.ForkID, s.Timestamp, s.ReceivedAt, s.L1InfoRoot, s.Source)
}

func (s *SequencedBatches) IsEqual(o interface{}) bool {
	other, ok := o.(*SequencedBatches)
	if !ok {
		return false
	}
	if s == other {
		return true
	}
	return *s == *other
}

func (s *SequencedBatches) Key() uint64 {
	return s.FromBatchNumber
}

func NewSequencedBatches(fromBatchNumber, toBatchNumber, l1BlockNumber, forkID uint64,
	timestamp time.Time, receivedAt time.Time, l1InfoRoot common.Hash, source string) *SequencedBatches {
	return &SequencedBatches{
		FromBatchNumber: fromBatchNumber,
		ToBatchNumber:   toBatchNumber,
		L1BlockNumber:   l1BlockNumber,
		ForkID:          forkID,
		Timestamp:       timestamp,
		ReceivedAt:      receivedAt,
		L1InfoRoot:      l1InfoRoot,
		Source:          source,
	}
}

func (seq *SequencesBatchesSlice) GetMinimumBlockNumber(defaultBatchNumber uint64) uint64 {
	if seq == nil {
		return defaultBatchNumber
	}
	minBlockNumber := uint64(0)
	for _, s := range *seq {
		if s.L1BlockNumber < minBlockNumber || minBlockNumber == 0 {
			minBlockNumber = s.L1BlockNumber
		}
	}
	if minBlockNumber >= defaultBatchNumber {
		log.Warnf("GetMinimumBlockNumber: minBlockNumber %d is greater than defaultBatchNumber %d", minBlockNumber, defaultBatchNumber)
		return defaultBatchNumber
	}
	return minBlockNumber
}

func (seq *SequencesBatchesSlice) Len() int {
	if seq == nil {
		return 0
	}
	return len(*seq)
}

func (seq *SequencesBatchesSlice) NumBatchesIncluded() int {
	if seq == nil {
		return 0
	}
	numBatches := 0
	for _, s := range *seq {
		numBatches += int(s.ToBatchNumber - s.FromBatchNumber + 1)
	}
	return numBatches
}
