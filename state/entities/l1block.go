package entities

import (
	"fmt"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/ethereum/go-ethereum/common"
)

// L1Block struct
type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool // The block is safe (have past the safe point, e.g. Finalized in L1)
	SyncVersion string
}

func (b *L1Block) String() string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprintf("BlockNumber: %d, BlockHash: %s, ParentHash: %s, ReceivedAt: %s, Checked: %t, SyncVersion: %s",
		b.BlockNumber, b.BlockHash.String(), b.ParentHash.String(), b.ReceivedAt.String(), b.Checked, b.SyncVersion)
}

func NewL1BlockFromEthermanBlock(block *etherman.Block, isFinalized bool) *L1Block {
	return &L1Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
		SyncVersion: zkevm_synchronizer_l1.Version,
		Checked:     isFinalized,
	}
}
