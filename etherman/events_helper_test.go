package etherman

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocketh "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/mocks"
	ethtypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
)

func TestAddaddNewBlockToResult(t *testing.T) {
	ctx := context.TODO()
	blockRetrieverMock := mocketh.NewBlockRetriever(t)
	vLog := types.Log{}
	blocks := []ethtypes.Block{}
	blocksOrder := map[common.Hash][]ethtypes.Order{}
	block := &ethtypes.Block{}
	blockRetrieverMock.EXPECT().RetrieveFullBlockForEvent(ctx, mock.Anything).Return(block, nil)
	block, err := addNewBlockToResult(ctx, blockRetrieverMock, vLog, &blocks, &blocksOrder)
	require.NoError(t, err)
	require.NotNil(t, block)
}
