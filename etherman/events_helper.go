package etherman

import (
	"context"

	ethtypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func addNewOrder(order *ethtypes.Order, blockHash common.Hash, blocksOrder *map[common.Hash][]ethtypes.Order) {
	(*blocksOrder)[blockHash] = append((*blocksOrder)[blockHash], *order)
}

// addNewEvent adds a new event to the blocks array and order array.
// it returns the block that must be filled with event data
func addNewBlockToResult(ctx context.Context, blockRetriever ethtypes.BlockRetriever, vLog types.Log, blocks *[]ethtypes.Block, blocksOrder *map[common.Hash][]ethtypes.Order) (*ethtypes.Block, error) {
	var block *ethtypes.Block
	var err error
	if !isheadBlockInArray(blocks, vLog.BlockHash, vLog.BlockNumber) {
		// Need to add the block, doesnt mind if inside the blocks because I have to respect the order so insert at end
		block, err = blockRetriever.RetrieveFullBlockForEvent(ctx, vLog)
		if err != nil {
			return nil, err
		}
		*blocks = append(*blocks, *block)
	}
	block = &(*blocks)[len(*blocks)-1]
	return block, nil
}
