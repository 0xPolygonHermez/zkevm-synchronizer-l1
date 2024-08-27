package etherman

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestBananaEventsNoTopics(t *testing.T) {
	sut := Client{}
	ctx := context.TODO()
	vLog := types.Log{}
	blocks := &[]Block{}
	blocksOrder := &map[common.Hash][]Order{}

	processed, err := sut.processBananaEvent(ctx, vLog, blocks, blocksOrder)
	require.NoError(t, err)
	require.False(t, processed)
}

func TestBananaEventsUnknownTopics(t *testing.T) {
	sut := Client{}
	ctx := context.TODO()
	vLog := types.Log{
		Topics: []common.Hash{common.HexToHash("0x12345678")},
	}
	blocks := &[]Block{}
	blocksOrder := &map[common.Hash][]Order{}

	processed, err := sut.processBananaEvent(ctx, vLog, blocks, blocksOrder)
	require.NoError(t, err)
	require.False(t, processed)
}

/*
func TestBananaEventsRollbackBatchesSignatureHash(t *testing.T) {
	bananaZkEVM, err := bananarollupcontract.NewPolygonvalidiumetrog(common.HexToAddress("0x12"), nil)
	require.NoError(t, err)
	sut := Client{
		BananaZkEVM: bananaZkEVM,
	}
	ctx := context.TODO()
	vlogData := generateVLogRollbackBatches(t)
	vLog := types.Log{
		Topics: []common.Hash{rollbackBatchesSignatureHash},
		Data:   vlogData[23:],
	}
	blocks := &[]Block{}
	blocksOrder := &map[common.Hash][]Order{}

	processed, err := sut.processBananaEvent(ctx, vLog, blocks, blocksOrder)
	require.NoError(t, err)
	require.False(t, processed)
}

func generateVLogRollbackBatches(t *testing.T) []byte {
	log.Debugf("rollbackBatchesSignatureHash %v", rollbackBatchesSignatureHash.String())
	abi, err := abi.JSON(strings.NewReader(bananarollupcontract.PolygonvalidiumetrogABI))
	require.NoError(t, err)
	event, err := abi.EventByID(rollbackBatchesSignatureHash)
	require.NoError(t, err)
	accHash := [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	data, err := event.Inputs.Pack(uint64(0xffffffffffffffff), accHash)
	require.NoError(t, err)
	return data
}
*/
