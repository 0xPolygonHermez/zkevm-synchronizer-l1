package internal

import (
	"context"
	"testing"

	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
)

func TestKk(t *testing.T) {
	testData := newTestDataSyncImpl(t)
	testData.sut.Sync(false)

}

type testDataSyncImpl struct {
	mockStorage           *mock_syncinterfaces.StorageInterface
	mockState             *mock_syncinterfaces.StateInterface
	mockEtherman          *mock_syncinterfaces.EthermanFullInterface
	mockStorageChecker    *mock_syncinterfaces.StorageCompatibilityChecker
	mockL1EventProcessors *mock_syncinterfaces.L1EventProcessorManager
	//	mockL1Sync           *mock_syncinterfaces.L1SequentialSync
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
	//mockL1EventProcessors := mock_syncinterfaces.NewL1EventProcessorManager(t)
	//mockL1Sync := mock_syncinterfaces.NewL1SequentialSync(t)
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
		l1EventProcessors:   nil,
		l1Sync:              nil,
		blockRangeProcessor: mockBlockRangeProcessor,
	}
	return &testDataSyncImpl{
		mockStorage:        mockStorage,
		mockState:          mockState,
		mockEtherman:       mockEtherman,
		mockStorageChecker: mockStorageChecker,
		//mockL1EventProcessors: mockL1EventProcessors,
		//mockL1Sync:           mockL1Sync,
		mockBlockRangeProcessor: mockBlockRangeProcessor,
		sut:                     sut,
		ctx:                     ctx,
	}
}
