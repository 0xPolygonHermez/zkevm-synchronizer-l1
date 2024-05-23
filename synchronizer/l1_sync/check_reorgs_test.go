package l1sync_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	mock_l1sync "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/l1_sync/mocks"
)

// TestReorgBadParams bad params in the call to CheckReorg
func TestReorgBadParams(t *testing.T) {
	testData := newReorgTestData(t)

	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(nil, nil)
	require.Error(t, err)
	require.Equal(t, uint64(0), lastBadBlockNumber)
	require.Nil(t, firstBlockOk)
}

// TestReorgNoReorgDetectedReturnsNoBlock no reorg detected
// If no reorg returns all values as nil
func TestReorgNoReorgDetectedReturnsNoBlock(t *testing.T) {
	testData := newReorgTestData(t)
	remoteBlock := newEthBlock(123, common.HexToHash("0x1234"))
	localBlock := newStateBlock(remoteBlock, true, true)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock.BlockNumber).Return(remoteBlock, nil)
	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(localBlock, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(0), lastBadBlockNumber)
	require.Nil(t, firstBlockOk)
}

// TestReorgReorgDetectedOnDepth1
// - block 123: wrong
// - block 120: ok
// So it must return:
//   - firstBlockOk: 120
//   - lastBadBlockNumber: 123
func TestReorgReorgDetectedOnDepth1(t *testing.T) {
	testData := newReorgTestData(t)
	remoteBlock123 := newEthBlock(123, common.HexToHash("0x1234"))
	localBlock123 := newStateBlock(remoteBlock123, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlock123 = newEthBlock(123, common.HexToHash("0x123456"))

	remoteBlock120 := newEthBlock(120, common.HexToHash("0x1111"))
	localBlock120 := newStateBlock(remoteBlock120, true, true)

	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock123.BlockNumber).Return(remoteBlock123, nil)
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(1), nil).Return(localBlock120, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock120.BlockNumber).Return(remoteBlock120, nil)

	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(localBlock123, nil)
	require.NoError(t, err)
	require.Equal(t, uint64(123), lastBadBlockNumber)
	require.Equal(t, localBlock120.BlockNumber, firstBlockOk.BlockNumber)
}

// TestReorgReorgDetectedOnDepth1UsingRollupDataInfoReturned
// In that case first block instead of getting from L1 used the one returned by RollupData
//
// - block 123: wrong (comming from Rollup Data)
// - block 120: ok
// So it must return:
//   - firstBlockOk: 120
//   - lastBadBlockNumber: 123
func TestReorgReorgDetectedOnDepth1UsingRollupDataInfoReturned(t *testing.T) {
	testData := newReorgTestData(t)
	localBlock123 := newStateBlock(newEthBlock(123, common.HexToHash("0x1234")), true, true)
	rollupBlock123 := &etherman.Block{
		BlockNumber: 123,
		BlockHash:   common.HexToHash("0x123467"), //Different from localBlock123
		ParentHash:  common.HexToHash("0x4566"),
	}
	remoteBlock120 := newEthBlock(120, common.HexToHash("0x1111"))
	localBlock120 := newStateBlock(remoteBlock120, true, true)

	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(1), nil).Return(localBlock120, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock120.BlockNumber).Return(remoteBlock120, nil)

	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(localBlock123, rollupBlock123)
	require.NoError(t, err)
	require.Equal(t, uint64(123), lastBadBlockNumber)
	require.Equal(t, localBlock120.BlockNumber, firstBlockOk.BlockNumber)
}

// TestReorgReorgDetectedOnDepth2
// - block 123: wrong
// - block 120: wrong
// - block 115: ok
// So it must return:
//   - firstBlockOk: 115
//   - lastBadBlockNumber: 120
func TestReorgReorgDetectedOnDepth2(t *testing.T) {
	testData := newReorgTestData(t)
	remoteBlock123 := newEthBlock(123, common.HexToHash("0x1234"))
	localBlock123 := newStateBlock(remoteBlock123, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlock123 = newEthBlock(123, common.HexToHash("0x123456"))

	remoteBlock120 := newEthBlock(120, common.HexToHash("0x1111"))
	localBlock120 := newStateBlock(remoteBlock120, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlock120 = newEthBlock(120, common.HexToHash("0x12345678"))

	remoteBlock115 := newEthBlock(115, common.HexToHash("0x1112"))
	localBlock115 := newStateBlock(remoteBlock115, true, true)

	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock123.BlockNumber).Return(remoteBlock123, nil)
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(1), nil).Return(localBlock120, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock120.BlockNumber).Return(remoteBlock120, nil)
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(2), nil).Return(localBlock115, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock115.BlockNumber).Return(remoteBlock115, nil)

	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(localBlock123, nil)
	require.NoError(t, err)
	require.Equal(t, localBlock115.BlockNumber, firstBlockOk.BlockNumber, "firstBlockOk")
	require.Equal(t, uint64(120), lastBadBlockNumber, "lastBadBlockNumber")
}

// TestReorgReorgDetectedOnDepth2
// - block 123: wrong
// - block 120: wrong
// - block genesis: wrong
// - No more blocks
// So it must return:
//   - firstBlockOk: nil (no block ok on DB!)
//   - lastBadBlockNumber: genesis
//   - err: ErrReorgAllBlocksOnDBAreBad
func TestReorgReorgDetectedOnGenesis(t *testing.T) {
	testData := newReorgTestData(t)
	remoteBlock123 := newEthBlock(123, common.HexToHash("0x1234"))
	localBlock123 := newStateBlock(remoteBlock123, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlock123 = newEthBlock(123, common.HexToHash("0x123456"))

	remoteBlock120 := newEthBlock(120, common.HexToHash("0x1111"))
	localBlock120 := newStateBlock(remoteBlock120, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlock120 = newEthBlock(120, common.HexToHash("0x12345678"))

	remoteBlockGenesis := newEthBlock(testData.genesisBlockNumber, common.HexToHash("0x1112"))
	localBlockGenesis := newStateBlock(remoteBlockGenesis, true, true)
	// I change parentHash to produce a discrepance between the local and remote block
	remoteBlockGenesis = newEthBlock(testData.genesisBlockNumber, common.HexToHash("0x11124"))

	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock123.BlockNumber).Return(remoteBlock123, nil)
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(1), nil).Return(localBlock120, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlock120.BlockNumber).Return(remoteBlock120, nil)
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(2), nil).Return(localBlockGenesis, nil)
	testData.mockEtherman.EXPECT().EthBlockByNumber(testData.ctx, localBlockGenesis.BlockNumber).Return(remoteBlockGenesis, nil)
	// No more blocks!, returns ErrNotFound
	testData.mockState.EXPECT().GetPreviousBlock(testData.ctx, uint64(3), nil).Return(nil, entities.ErrNotFound)
	firstBlockOk, lastBadBlockNumber, err := testData.sut.CheckReorg(localBlock123, nil)
	require.ErrorIs(t, err, l1sync.ErrReorgAllBlocksOnDBAreBad)
	require.Nil(t, firstBlockOk, "firstBlockOk")
	require.Equal(t, uint64(testData.genesisBlockNumber), lastBadBlockNumber, "lastBadBlockNumber")
}

// --- HELPER FUNCTIONS -----------------------------------------------
type reorgTestData struct {
	mockEtherman       *mock_l1sync.EthermanReorgManager
	mockState          *mock_l1sync.StateReorgInterface
	genesisBlockNumber uint64
	ctx                context.Context
	sut                *l1sync.CheckReorgManager
}

func newReorgTestData(t *testing.T) *reorgTestData {
	mockEtherman := mock_l1sync.NewEthermanReorgManager(t)
	mockState := mock_l1sync.NewStateReorgInterface(t)
	genesisBlockNumber := uint64(100)
	ctx := context.TODO()
	sut := l1sync.NewCheckReorgManager(ctx, mockEtherman, mockState)
	return &reorgTestData{
		mockEtherman:       mockEtherman,
		mockState:          mockState,
		genesisBlockNumber: genesisBlockNumber,
		ctx:                ctx,
		sut:                sut,
	}
}

func newEthBlock(number uint64, parentHash common.Hash) *ethTypes.Block {
	header := &ethTypes.Header{Number: big.NewInt(int64(number)),
		ParentHash: parentHash}
	return ethTypes.NewBlockWithHeader(header)
}

func newStateBlock(ethBLock *ethTypes.Block, checked, hasEvent bool) *entities.L1Block {
	return &entities.L1Block{
		BlockNumber: ethBLock.Number().Uint64(),
		BlockHash:   ethBLock.Hash(),
		ParentHash:  ethBLock.ParentHash(),
		Checked:     checked,
		HasEvents:   hasEvent,
	}
}
