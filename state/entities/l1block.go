package entities

import (
	"fmt"
	"time"

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
