package entities

import (
	"fmt"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	ethtypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/ethereum/go-ethereum/common"
)

// L1Block struct
type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool // The block is safe (have past the safe point, e.g. Finalized in L1)
	HasEvents   bool // This block have events from the rollup
	SyncVersion string
}

func (b *L1Block) Key() uint64 {
	return b.BlockNumber
}

func (b *L1Block) String() string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprintf("BlockNumber: %d, BlockHash: %s, ParentHash: %s, ReceivedAt: %s, Checked: %t, SyncVersion: %s",
		b.BlockNumber, b.BlockHash.String(), b.ParentHash.String(), b.ReceivedAt.String(), b.Checked, b.SyncVersion)
}

func NewL1BlockFromEthermanBlock(block *ethtypes.Block, isFinalized bool) *L1Block {
	return &L1Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
		SyncVersion: zkevm_synchronizer_l1.Version,
		Checked:     isFinalized,
		HasEvents:   block.HasEvents(),
	}
}

func IsBlockFinalized(blockNumber uint64, finalizedBlockNumber uint64) bool {
	return blockNumber <= finalizedBlockNumber
}

func (b *L1Block) IsUnsafeAndHaveRollupdata() bool {
	if b == nil {
		return false
	}
	return !b.Checked && b.HasEvents
}
