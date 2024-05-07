package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type L1Block = entities.L1Block
type stateTxType = entities.Tx

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
