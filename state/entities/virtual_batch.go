package entities

import (
	"fmt"
	"time"

	ethtypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/ethereum/go-ethereum/common"
)

type VirtualBatch struct {
	BatchNumber             uint64
	ForkID                  uint64
	BatchL2Data             []byte
	VlogTxHash              common.Hash // Hash of tx inside L1Block that emit this log
	Coinbase                common.Address
	SequencerAddr           common.Address
	SequenceFromBatchNumber uint64 // Linked to sync.sequenced_batches table
	BlockNumber             uint64 // Linked to sync.block table
	L1InfoRoot              *common.Hash
	ReceivedAt              time.Time
	BatchTimestamp          *time.Time // This is optional depend on ForkID
	ExtraInfo               *string
}

type BatchExtraInfo struct {
	Description string
}

func (s *VirtualBatch) IsEqual(o interface{}) bool {
	other, ok := o.(*VirtualBatch)
	if !ok {
		return false
	}
	if s == other {
		return true
	}
	return s.String() == other.String()
}

func (b *VirtualBatch) Key() uint64 {
	return b.BatchNumber
}

func (b *VirtualBatch) String() string {
	if b == nil {
		return "nil"
	}
	res := fmt.Sprintf("BatchNumber: %d, ForkID: %d, BatchL2Data: %s, TxHash: %s, Coinbase: %s, SequencerAddr: %s, BlockNumber: %d, L1InfoRoot: %s, ReceivedAt: %s, BatchTimestamp: %s,",
		b.BatchNumber, b.ForkID, string(b.BatchL2Data), b.VlogTxHash.String(), b.Coinbase.String(), b.SequencerAddr.String(), b.BlockNumber, b.L1InfoRoot.String(), b.ReceivedAt.String(), b.BatchTimestamp.String())
	if b.ExtraInfo != nil {
		res += fmt.Sprintf(", ExtraInfo: %s", *b.ExtraInfo)
	}
	return res
}

func NewVirtualBatchFromL1(l1BlockNumber, seqFromBatchNumber, forkID uint64, ethSeqBatch ethtypes.SequencedBatch) *VirtualBatch {
	res := &VirtualBatch{
		BatchNumber:             ethSeqBatch.BatchNumber,
		ForkID:                  forkID,
		BatchL2Data:             ethSeqBatch.BatchL2Data(),
		VlogTxHash:              ethSeqBatch.TxHash,
		Coinbase:                ethSeqBatch.Coinbase,
		SequencerAddr:           ethSeqBatch.SequencerAddr,
		SequenceFromBatchNumber: seqFromBatchNumber,
		BlockNumber:             l1BlockNumber,
		L1InfoRoot:              ethSeqBatch.L1InfoRoot,
		ReceivedAt:              time.Now(),
	}
	if ethSeqBatch.SequencedBatchElderberryData != nil {
		tstamp := time.Unix(int64(ethSeqBatch.SequencedBatchElderberryData.MaxSequenceTimestamp), 0)
		res.BatchTimestamp = &tstamp
	}
	return res
}

// VirtualBatchConstraints is a struct that contains the constraints to filter the virtual batches.
// is ready to add constraints to the query.
type VirtualBatchConstraints struct {
	batchNumberEqual *uint64
	batchNumberGt    *uint64
	batchNumberLt    *uint64
}

func (c *VirtualBatchConstraints) BatchNumberEqual(batchNumber uint64) {
	c.batchNumberEqual = &batchNumber
}

func (c *VirtualBatchConstraints) BatchNumberGt(batchNumber uint64) {
	c.batchNumberGt = &batchNumber
}

func (c *VirtualBatchConstraints) BatchNumberLt(batchNumber uint64) {
	c.batchNumberLt = &batchNumber
}

func (c *VirtualBatchConstraints) WhereClause() string {
	res := ""
	if c.batchNumberEqual != nil {
		res += fmt.Sprintf("batch_num = %d ", *c.batchNumberEqual)
	}
	if c.batchNumberGt != nil {
		res += fmt.Sprintf("batch_num>%d ", *c.batchNumberEqual)
	}
	if c.batchNumberLt != nil {
		res += fmt.Sprintf("batch_num<%d ", *c.batchNumberEqual)
	}
	return res
}
