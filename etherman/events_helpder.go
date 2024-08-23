package etherman

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func addNewOrder(order *Order, blockHash common.Hash, blocksOrder *map[common.Hash][]Order) {
	(*blocksOrder)[blockHash] = append((*blocksOrder)[blockHash], *order)
}

type BlockRetriever interface {
	RetrieveFullBlockForEvent(ctx context.Context, vLog types.Log) (*Block, error)
}

// addNewEvent adds a new event to the blocks array and order array.
// it returns the block that must be filled with event data
func addNewBlockToResult(ctx context.Context, blockRetriever BlockRetriever, vLog types.Log, blocks *[]Block, blocksOrder *map[common.Hash][]Order) (*Block, error) {
	var block *Block
	var err error
	if !isheadBlockInArray(blocks, vLog.BlockHash, vLog.BlockNumber) {
		// Need to add the block, doesnt mind if inside the blocks because I have to respect the order so insert at end
		//TODO: Check if the block is already in the blocks array and copy it instead of retrieve it again
		block, err = blockRetriever.RetrieveFullBlockForEvent(ctx, vLog)
		if err != nil {
			return nil, err
		}
		*blocks = append(*blocks, *block)
	}
	block = &(*blocks)[len(*blocks)-1]
	return block, nil
}
