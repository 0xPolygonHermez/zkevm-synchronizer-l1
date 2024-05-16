package l1sync_test

import (
	"testing"

	l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync"
	"github.com/stretchr/testify/require"
)

func TestBlockRange_WhenFromIsGEThanMaximumIsLastRange(t *testing.T) {
	it := l1sync.NewBlockRangeIterator(100, 10, 120)
	require.False(t, it.IsLastRange())
	it = it.NextRange(100)
	require.True(t, it.IsLastRange())
}

func TestBlockRange_NoOverlappedIsNextToFrom(t *testing.T) {
	fromBlock := uint64(100)
	chunck := uint64(10)
	lastBlock := uint64(200)
	it := l1sync.NewBlockRangeIterator(fromBlock, chunck, lastBlock)
	rangeBlock := it.GetRange(false)
	require.Equal(t, l1sync.BlockRange{
		FromBlock:            fromBlock + 1,
		ToBlock:              fromBlock + chunck,
		OverlappedFirstBlock: false,
	}, rangeBlock)
}

func TestBlockRange_OverlappedIncludeFrom(t *testing.T) {
	fromBlock := uint64(100)
	chunck := uint64(10)
	lastBlock := uint64(200)
	it := l1sync.NewBlockRangeIterator(fromBlock, chunck, lastBlock)
	rangeBlock := it.GetRange(true)
	require.Equal(t, l1sync.BlockRange{
		FromBlock:            fromBlock,
		ToBlock:              fromBlock + chunck,
		OverlappedFirstBlock: true,
	}, rangeBlock)
}

func TestBlockRange_EachNextRangeExtendToBlockByChunckIndependentlyFrom(t *testing.T) {
	fromBlock := uint64(100)
	chunck := uint64(10)
	lastBlock := uint64(200)
	it := l1sync.NewBlockRangeIterator(fromBlock, chunck, lastBlock)
	it = it.NextRange(fromBlock)
	require.NotNil(t, it)
	br := it.GetRange(false)
	require.Equal(t, fromBlock+2*chunck, br.ToBlock)

	it = it.NextRange(fromBlock + 5)
	require.NotNil(t, it)
	br = it.GetRange(false)
	require.Equal(t, fromBlock+3*chunck, br.ToBlock)
}

func TestBlockRange_ToBlockIsCappedToMaximumBlock(t *testing.T) {
	fromBlock := uint64(100)
	chunck := uint64(10)
	lastBlock := uint64(110)
	it := l1sync.NewBlockRangeIterator(fromBlock, chunck, lastBlock)
	it = it.NextRange(fromBlock)
	require.NotNil(t, it)
	br := it.GetRange(false)
	require.Equal(t, fromBlock+1*chunck, br.ToBlock)

	it = it.NextRange(fromBlock + 5)
	require.NotNil(t, it)
	br = it.GetRange(false)
	require.Equal(t, fromBlock+1*chunck, br.ToBlock)
}

func TestBlockRange_SettingWrongFromBlockReturnsNil(t *testing.T) {
	chunck := uint64(10)
	lastBlock := uint64(110)
	fromBlock := lastBlock + 1
	it := l1sync.NewBlockRangeIterator(fromBlock, chunck, lastBlock)
	require.Nil(t, it)
}
