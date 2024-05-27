package internal

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestFirstExecutionIsSynced
// This test is that receive 1 block and is synced
// is called to return on sync to it returns with no error
func TestSyncImplFirstExecutionIsSynced(t *testing.T) {
	testData := newTestDataSyncImpl(t)
	testData.mockStorage.EXPECT().GetLastBlock(mock.Anything, mock.Anything).Return(nil, entities.ErrNotFound)
	block := entities.L1Block{
		BlockNumber: 123,
	}
	testData.mockL1Syncer.EXPECT().SyncBlocks(testData.ctx, mock.Anything).Return(&block, true, nil)
	err := testData.sut.Sync(FlagReturnBeforeReorg | FlagReturnOnSync)
	require.NoError(t, err)
}

// TestSyncImplReturnsBeforeReorg
// SyncBlocks detects a reorg in block 123
// the Sync() call returns this error
func TestSyncImplReturnsBeforeReorg(t *testing.T) {
	testData := newTestDataSyncImpl(t)
	testData.mockStorage.EXPECT().GetLastBlock(mock.Anything, mock.Anything).Return(nil, entities.ErrNotFound)
	block := entities.L1Block{
		BlockNumber: 123,
	}
	errReorg := common.NewReorgError(123, fmt.Errorf("reorg"))
	testData.mockL1Syncer.EXPECT().SyncBlocks(testData.ctx, mock.Anything).Return(&block, false, errReorg)
	err := testData.sut.Sync(FlagReturnBeforeReorg | FlagReturnOnSync)
	require.ErrorIs(t, err, errReorg)
}

func TestSyncImplReturnsAfterReorg(t *testing.T) {
	testData := newTestDataSyncImpl(t)
	testData.mockStorage.EXPECT().GetLastBlock(mock.Anything, mock.Anything).Return(nil, entities.ErrNotFound)
	block := entities.L1Block{
		BlockNumber: 123,
	}
	errReorg := common.NewReorgError(123, fmt.Errorf("reorg"))
	testData.mockL1Syncer.EXPECT().SyncBlocks(testData.ctx, mock.Anything).Return(&block, false, errReorg)
	// It execute the reorg
	testData.mockState.EXPECT().BeginTransaction(testData.ctx).Return(testData.mockTx, nil)
	reorgReq := model.ReorgRequest{
		FirstL1BlockNumberToKeep: 122,
		ReasonError:              errReorg,
	}
	reorgResult := model.ReorgExecutionResult{
		Request:           reorgReq,
		ExecutionCounter:  1,
		ExecutionError:    nil,
		ExecutionTime:     time.Now(),
		ExecutionDuration: time.Second,
	}
	testData.mockState.EXPECT().ExecuteReorg(testData.ctx, reorgReq, testData.mockTx).Return(reorgResult)
	testData.mockTx.EXPECT().Commit(testData.ctx).Return(nil)
	err := testData.sut.Sync(FlagReturnAfterReorg | FlagReturnOnSync)
	require.NoError(t, err)
}

// --- HELPER FUNCTIONS ----------------------------------------------
type testDataSyncImpl struct {
	mockStorage             *mock_syncinterfaces.StorageInterface
	mockState               *mock_syncinterfaces.StateInterface
	mockEtherman            *mock_syncinterfaces.EthermanFullInterface
	mockStorageChecker      *mock_syncinterfaces.StorageCompatibilityChecker
	mockL1Syncer            *mock_syncinterfaces.L1Syncer
	mockBlockRangeProcessor *mock_syncinterfaces.BlockRangeProcessor
	mockTx                  *mock_entities.Tx
	sut                     *SynchronizerImpl
	ctx                     context.Context
}

func newTestDataSyncImpl(t *testing.T) *testDataSyncImpl {
	ctx, cancel := context.WithCancel(context.TODO())
	mockStorage := mock_syncinterfaces.NewStorageInterface(t)
	mockState := mock_syncinterfaces.NewStateInterface(t)
	mockEtherman := mock_syncinterfaces.NewEthermanFullInterface(t)
	mockStorageChecker := mock_syncinterfaces.NewStorageCompatibilityChecker(t)
	mockL1Syncer := mock_syncinterfaces.NewL1Syncer(t)
	mockBlockRangeProcessor := mock_syncinterfaces.NewBlockRangeProcessor(t)
	mockTx := mock_entities.NewTx(t)
	cfg := syncconfig.Config{
		GenesisBlockNumber: 123,
	}

	sut := &SynchronizerImpl{
		storage:             mockStorage,
		state:               mockState,
		etherMan:            mockEtherman,
		ctx:                 ctx,
		cancelCtx:           cancel,
		genBlockNumber:      cfg.GenesisBlockNumber,
		cfg:                 cfg,
		networkID:           0,
		storageChecker:      mockStorageChecker,
		l1Sync:              mockL1Syncer,
		blockRangeProcessor: mockBlockRangeProcessor,
	}
	return &testDataSyncImpl{
		mockStorage:             mockStorage,
		mockState:               mockState,
		mockEtherman:            mockEtherman,
		mockStorageChecker:      mockStorageChecker,
		mockL1Syncer:            mockL1Syncer,
		mockBlockRangeProcessor: mockBlockRangeProcessor,
		mockTx:                  mockTx,
		sut:                     sut,
		ctx:                     ctx,
	}
}
