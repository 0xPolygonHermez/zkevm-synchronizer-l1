package internal

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFirstExecutionNoDataOnDb(t *testing.T) {
	testData := newTestDataSyncImpl(t)
	testData.mockStorage.EXPECT().GetLastBlock(mock.Anything, mock.Anything).Return(nil, entities.ErrNotFound)
	block := entities.L1Block{
		BlockNumber: 123,
	}
	testData.mockL1Syncer.EXPECT().SyncBlocks(testData.ctx, mock.Anything).Return(&block, true, nil)
	err := testData.sut.Sync(FlagReturnBeforeReorg | FlagReturnOnSync)
	require.NoError(t, err)
}

type testDataSyncImpl struct {
	mockStorage             *mock_syncinterfaces.StorageInterface
	mockState               *mock_syncinterfaces.StateInterface
	mockEtherman            *mock_syncinterfaces.EthermanFullInterface
	mockStorageChecker      *mock_syncinterfaces.StorageCompatibilityChecker
	mockL1Syncer            *mock_syncinterfaces.L1Syncer
	mockBlockRangeProcessor *mock_syncinterfaces.BlockRangeProcessor
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
		sut:                     sut,
		ctx:                     ctx,
	}
}
