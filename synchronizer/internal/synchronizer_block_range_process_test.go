package internal_test

import (
	"context"
	"fmt"
	"testing"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	etherman "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/internal"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testProcessBlockRangeData struct {
	mockState            *mock_syncinterfaces.StorageInterface
	mockTransactions     *mock_syncinterfaces.StateTxProvider
	mockForkId           *mock_syncinterfaces.StateForkIdQuerier
	mockL1EventProcessor *mock_syncinterfaces.L1EventProcessorManager
	DbTx                 *mock_entities.Tx
	sut                  *internal.BlockRangeProcess
	ctx                  context.Context
	blocksNoOrders       []etherman.Block
}

func newTestProcessBlockRangeData(t *testing.T) *testProcessBlockRangeData {
	mockState := mock_syncinterfaces.NewStorageInterface(t)
	mockForkId := mock_syncinterfaces.NewStateForkIdQuerier(t)
	mockL1EventProcessor := mock_syncinterfaces.NewL1EventProcessorManager(t)
	mockTransactions := mock_syncinterfaces.NewStateTxProvider(t)
	DbTx := mock_entities.NewTx(t)
	sut := internal.NewBlockRangeProcessLegacy(mockState, mockForkId, mockTransactions, mockL1EventProcessor)
	ctx := context.TODO()
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
		},
		{
			BlockNumber: 2,
		},
	}
	return &testProcessBlockRangeData{mockState, mockTransactions, mockForkId, mockL1EventProcessor, DbTx, sut, ctx, blocks}
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
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil)
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
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil)
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true}, data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil)
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}

// Each block is a new transaction
func TestProcessBlockMultiplesBLocksMultiplesDBTransactions(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	DbTx2 := mock_entities.NewTx(t)
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
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true}, data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil).Once()

	// Second iteration returns DbTx2
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(DbTx2, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 2, SyncVersion: zkevm_synchronizer_l1.Version, Checked: false}, DbTx2).Return(nil)
	DbTx2.EXPECT().Commit(data.ctx).Return(nil).Once()
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}

func TestProcessBlockMultiplesBLocksWithElements(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	DbTx2 := mock_entities.NewTx(t)
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
			ForkIDs: []etherman.ForkID{
				{
					BatchNumber: 123,
					ForkID:      10,
				},
			},
		},
		{
			BlockNumber: 2,
			BlockHash:   common.HexToHash("0x123"),
		},
	}
	order := map[common.Hash][]etherman.Order{}
	order[blocks[0].BlockHash] = []etherman.Order{
		{
			Name: etherman.ForkIDsOrder,
			Pos:  0,
		},
	}
	finalizedBlockNumber := uint64(1)
	// First BeginStateTransaction returns data.DbTx
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true, HasEvents: true}, data.DbTx).Return(nil)
	data.mockForkId.EXPECT().GetForkIDByBlockNumber(data.ctx, uint64(1), data.DbTx).Return(uint64(1))
	data.mockL1EventProcessor.EXPECT().Process(data.ctx, actions.ForkIdType(1), mock.Anything, &blocks[0], data.DbTx).Return(nil)
	data.DbTx.EXPECT().Commit(data.ctx).Return(nil).Once()

	// Second iteration returns DbTx2
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(DbTx2, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{
		BlockNumber: 2,
		BlockHash:   blocks[1].BlockHash,
		SyncVersion: zkevm_synchronizer_l1.Version, Checked: false}, DbTx2).Return(nil)
	DbTx2.EXPECT().Commit(data.ctx).Return(nil).Once()
	err := data.sut.ProcessBlockRange(data.ctx, blocks, order, finalizedBlockNumber)
	require.NoError(t, err)
}

func TestProcessBlockMultiplesBLocksWithElementsSingleTx(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	blocks := []etherman.Block{
		{
			BlockNumber: 1,
			ForkIDs: []etherman.ForkID{
				{
					BatchNumber: 123,
					ForkID:      10,
				},
			},
		},
		{
			BlockNumber: 2,
			BlockHash:   common.HexToHash("0x123"),
		},
	}
	order := map[common.Hash][]etherman.Order{}
	order[blocks[0].BlockHash] = []etherman.Order{
		{
			Name: etherman.ForkIDsOrder,
			Pos:  0,
		},
	}
	finalizedBlockNumber := uint64(1)
	// First BeginStateTransaction returns data.DbTx
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true, HasEvents: true}, mock.Anything).Return(nil)
	data.mockForkId.EXPECT().GetForkIDByBlockNumber(data.ctx, uint64(1), mock.Anything).Return(uint64(1))
	data.mockL1EventProcessor.EXPECT().Process(data.ctx, actions.ForkIdType(1), mock.Anything, &blocks[0], mock.Anything).Return(nil)

	// Second iteration reuse same tx
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{
		BlockNumber: 2,
		BlockHash:   blocks[1].BlockHash,
		SyncVersion: zkevm_synchronizer_l1.Version, Checked: false}, mock.Anything).Return(nil)
	err := data.sut.ProcessBlockRangeSingleDbTx(data.ctx, blocks, order, finalizedBlockNumber, syncinterfaces.StoreL1Blocks, data.DbTx)
	require.NoError(t, err)
}

func TestProcessBlockErrorBeginTransaction(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(nil, fmt.Errorf("error"))

	err := data.sut.ProcessBlockRange(data.ctx, data.blocksNoOrders, nil, uint64(123))
	require.Error(t, err)
}
func TestProcessBlockErrorAddBlock(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true, HasEvents: false}, data.DbTx).Return(fmt.Errorf("error"))
	data.DbTx.EXPECT().Rollback(data.ctx).Return(nil).Once()

	err := data.sut.ProcessBlockRange(data.ctx, data.blocksNoOrders, nil, uint64(123))
	require.Error(t, err)
}

func TestProcessBlockErrorCommit(t *testing.T) {
	data := newTestProcessBlockRangeData(t)
	data.mockTransactions.EXPECT().BeginTransaction(data.ctx).Return(data.DbTx, nil).Once()
	data.mockState.EXPECT().AddBlock(data.ctx, &pgstorage.L1Block{BlockNumber: 1, SyncVersion: zkevm_synchronizer_l1.Version, Checked: true, HasEvents: false}, data.DbTx).Return(nil)
	returnErr := fmt.Errorf("error")
	data.DbTx.EXPECT().Commit(data.ctx).Return(returnErr).Once()

	err := data.sut.ProcessBlockRange(data.ctx, data.blocksNoOrders, nil, uint64(123))
	require.Error(t, err)
	require.ErrorIs(t, err, returnErr)
}
