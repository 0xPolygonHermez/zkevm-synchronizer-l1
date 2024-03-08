package synchronizer

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/ethereum/go-ethereum/common"
)

// L1Block struct
type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
}

func convertEthermanBlock(block *etherman.Block) *L1Block {
	if block == nil {
		return nil
	}
	return &L1Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
	}
}

func convertArrayEthermanBlocks(ethermanBlocks []etherman.Block) []L1Block {
	l1Blocks := make([]L1Block, len(ethermanBlocks))
	for i, block := range ethermanBlocks {
		l1Blocks[i] = *convertEthermanBlock(&block)
	}
	return l1Blocks
}
