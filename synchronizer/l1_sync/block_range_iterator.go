package l1sync

import (
	"fmt"
	"log"
)

type BlockRange struct {
	FromBlock, ToBlock   uint64
	OverlappedFirstBlock bool
}

func (b BlockRange) InsideRange(blockNumber uint64) bool {
	if b.OverlappedFirstBlock {
		return blockNumber >= b.FromBlock && blockNumber <= b.ToBlock
	}
	return blockNumber > b.FromBlock && blockNumber <= b.ToBlock
}

func (b BlockRange) String() string {
	return fmt.Sprintf("FromBlock: %d, ToBlock: %d overlappedFirstBlock:%t", b.FromBlock, b.ToBlock, b.OverlappedFirstBlock)
}

type BlockRangeIterator struct {
	fromBlock, toBlock uint64
	SyncChunkSize      uint64
	MaximumBlock       uint64
}

func NewBlockRangeIterator(fromBlock, syncChunkSize uint64, maximumBlock uint64) *BlockRangeIterator {
	res := &BlockRangeIterator{
		fromBlock:     fromBlock,
		toBlock:       fromBlock,
		SyncChunkSize: syncChunkSize,
		MaximumBlock:  maximumBlock,
	}
	res = res.NextRange(fromBlock)
	return res
}
func (i *BlockRangeIterator) IsLastRange() bool {
	return i.toBlock >= i.MaximumBlock
}

func (i *BlockRangeIterator) NextRange(fromBlock uint64) *BlockRangeIterator {
	// The FromBlock is the new block (can be the previous one if no blocks found in the range)
	if fromBlock < i.fromBlock {
		log.Fatal("FromBlock is less than the current fromBlock")
	}
	i.fromBlock = fromBlock
	// Extend toBlock by sync chunk size
	i.toBlock = i.toBlock + i.SyncChunkSize

	if i.toBlock > i.MaximumBlock {
		i.toBlock = i.MaximumBlock
	}
	if i.fromBlock >= i.toBlock {
		return nil
	}
	return i
}

func (i *BlockRangeIterator) GetRange(overlappedFirst bool) BlockRange {
	if overlappedFirst {
		return BlockRange{
			FromBlock:            i.fromBlock,
			ToBlock:              i.toBlock,
			OverlappedFirstBlock: overlappedFirst,
		}
	}
	return BlockRange{
		FromBlock:            i.fromBlock + 1,
		ToBlock:              i.toBlock,
		OverlappedFirstBlock: overlappedFirst,
	}
}

func (i *BlockRangeIterator) String() string {
	return fmt.Sprintf("FromBlock: %d, ToBlock: %d, MaximumBlock: %d", i.fromBlock, i.toBlock, i.MaximumBlock)
}
