package synchronizer_test

import (
	"context"
	"testing"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
	mock_synchronizer "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type testProcessBlockRangeData struct {
	mockState            *mock_synchronizer.StorageInterface
	mockForkId           *mock_synchronizer.StateForkIdQuerier
	mockL1EventProcessor *mock_synchronizer.L1EventProcessorManager
	DbTx                 *mock_synchronizer.DbTxMock
	sut                  *synchronizer.BlockRangeProcess
	ctx                  context.Context
}

func newTestProcessBlockRangeData(t *testing.T) *testProcessBlockRangeData {
	mockState := mock_synchronizer.NewStorageInterface(t)
	mockForkId := mock_synchronizer.NewStateForkIdQuerier(t)
	mockL1EventProcessor := mock_synchronizer.NewL1EventProcessorManager(t)
	DbTx := mock_synchronizer.NewDbTxMock(t)
	sut := synchronizer.NewBlockRangeProcessLegacy(mockState, mockForkId, mockL1EventProcessor, nil)
	ctx := context.TODO()
	return &testProcessBlockRangeData{mockState, mockForkId, mockL1EventProcessor, DbTx, sut, ctx}
}

func TestProcessBlockWithNoOrderJustWriteItOnDBAfterFinalizedIsStoreAsCheckedFalse(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
		},
	}
	order := map[common.Hash][]etherman.Order{}
	finalizedBlockNumber := uint64(0)
	data.mockState.EXPECT().BeginStateTransaction(data.ctx).Return(data.DbTx, nil)
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: false}, data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil)
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}

// If the stored block is <= finalized -> Checked = true
func TestProcessBlockWithNoOrderJustWriteItOnDBEqualToFinalizedIsStoreAsCheckedTrue(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
		},
	}
	order := map[common.Hash][]etherman.Order{}
	finalizedBlockNumber := uint64(1)
	data.mockState.EXPECT().BeginStateTransaction(data.ctx).Return(data.DbTx, nil)
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true}, data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil)
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}

// Each block is a new transaction
func TestProcessBlockMultiplesBLocksMultiplesDBTransactions(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	DbTx2 := mock_synchronizer.NewDbTxMock(t)
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
		},
		{
			BlockNumber: 2,
		},
	}
	order := map[common.Hash][]etherman.Order{}
	finalizedBlockNumber := uint64(1)
	// First BeginStateTransaction returns data.DbTx
	data.mockState.EXPECT().BeginStateTransaction(data.ctx).Return(data.DbTx, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true}, data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil).Once()

	// Second iteration returns DbTx2
	data.mockState.EXPECT().BeginStateTransaction(data.ctx).Return(DbTx2, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 2, SyncVersion: zkevm_synchronizer_l1.Version, Checked: false}, DbTx2).Return(nil)
	DbTx2.EXPECT().Commit(data.ctx).Return(nil).Once()
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}
