package internal_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/internal"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
	"github.com/stretchr/testify/require"
)

type testDataForSanityStorageChecker struct {
	mockStateCompatibility *mock_syncinterfaces.StateStorageCompatibilityCheckerInterface
	mockEtherman           *mock_syncinterfaces.EthermanChainQuerier
	sut                    *internal.StorageCompatibilityCheckerImpl
	ctx                    context.Context
	overrideStorageCheck   bool
}

func newTestDataForSanityStorageChecker(t *testing.T, overrideStorageCheck bool) *testDataForSanityStorageChecker {
	mockStateCompatibility := mock_syncinterfaces.NewStateStorageCompatibilityCheckerInterface(t)
	mockEtherman := mock_syncinterfaces.NewEthermanChainQuerier(t)
	sut := internal.NewSanityStorageCheckerImpl(mockStateCompatibility, mockEtherman, overrideStorageCheck)
	ctx := context.TODO()
	return &testDataForSanityStorageChecker{mockStateCompatibility, mockEtherman, sut, ctx, overrideStorageCheck}
}

func TestStorageCheckerNoError(t *testing.T) {
	testData := newTestDataForSanityStorageChecker(t, false)
	testData.mockEtherman.On("GetRollupID").Return(uint(1))
	testData.mockEtherman.On("GetL1ChainID").Return(uint64(10))
	currentContetsBoundData := entities.StorageContentsBoundData{
		RollupID:  1,
		L1ChainID: 10,
	}
	testData.mockStateCompatibility.On("CheckAndUpdateStorage", testData.ctx, currentContetsBoundData, testData.overrideStorageCheck, nil).Return(nil)

	err := testData.sut.CheckAndUpdateStorage(testData.ctx)

	require.NoError(t, err)
}

func TestStorageCheckerError(t *testing.T) {
	testData := newTestDataForSanityStorageChecker(t, false)
	testData.mockEtherman.On("GetRollupID").Return(uint(1))
	testData.mockEtherman.On("GetL1ChainID").Return(uint64(10))
	currentContetsBoundData := entities.StorageContentsBoundData{
		RollupID:  1,
		L1ChainID: 10,
	}
	returnedError := fmt.Errorf("test error")
	testData.mockStateCompatibility.On("CheckAndUpdateStorage", testData.ctx, currentContetsBoundData, testData.overrideStorageCheck, nil).Return(returnedError)

	err := testData.sut.CheckAndUpdateStorage(testData.ctx)

	require.Error(t, err)
}

func TestStorageCheckerOverrider(t *testing.T) {
	testData := newTestDataForSanityStorageChecker(t, true)
	testData.mockEtherman.On("GetRollupID").Return(uint(1))
	testData.mockEtherman.On("GetL1ChainID").Return(uint64(10))
	currentContetsBoundData := entities.StorageContentsBoundData{
		RollupID:  1,
		L1ChainID: 10,
	}
	returnedError := fmt.Errorf("test error")
	testData.mockStateCompatibility.On("CheckAndUpdateStorage", testData.ctx, currentContetsBoundData, testData.overrideStorageCheck, nil).Return(returnedError)

	err := testData.sut.CheckAndUpdateStorage(testData.ctx)

	require.Error(t, err)
}
