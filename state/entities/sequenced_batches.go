package entities

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

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

func (s *SequencedBatches) Key() uint64 {
	return s.FromBatchNumber
}
