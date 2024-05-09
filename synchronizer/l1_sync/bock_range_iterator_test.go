package l1sync_test

import (
	"testing"

	l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync"
	"github.com/stretchr/testify/require"
)

func TestKK(t *testing.T) {
	it := l1sync.NewBlockRangeIterator(100, 10, 120)
	require.False(t, it.IsLastRange())
	it = it.NextRange(100)
	require.True(t, it.IsLastRange())
}
