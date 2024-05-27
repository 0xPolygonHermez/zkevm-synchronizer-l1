package pgstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	mock_pgstorage "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage/mocks"
	"github.com/stretchr/testify/require"
)

type txTestData struct {
	ctx       context.Context
	sut       entities.Tx
	mockDbTx  *mock_pgstorage.DbTxMock
	testError error
}

func newTxTestData(t *testing.T) *txTestData {
	mockDbTx := mock_pgstorage.NewDbTxMock(t)
	ctx := context.TODO()
	sut := pgstorage.NewTxImpl(mockDbTx)
	return &txTestData{
		ctx:       ctx,
		sut:       sut,
		mockDbTx:  mockDbTx,
		testError: fmt.Errorf("test error"),
	}
}

func TestTxCommitCbNoError(t *testing.T) {
	testData := newTxTestData(t)
	testData.mockDbTx.EXPECT().Commit(testData.ctx).Return(nil).Once()

	calledCommitCb := false
	testData.sut.AddCommitCallback(func(tx entities.Tx, err error) {
		require.NoError(t, err)
		calledCommitCb = true
	})
	err := testData.sut.Commit(testData.ctx)
	require.NoError(t, err)
	require.True(t, calledCommitCb)
}

func TestTxCommitCbError(t *testing.T) {
	testData := newTxTestData(t)
	testData.mockDbTx.EXPECT().Commit(testData.ctx).Return(testData.testError).Once()

	calledCommitCb := false
	testData.sut.AddCommitCallback(func(tx entities.Tx, err error) {
		require.ErrorIs(t, err, testData.testError)
		calledCommitCb = true
	})
	err := testData.sut.Commit(testData.ctx)
	require.ErrorIs(t, err, testData.testError)
	require.True(t, calledCommitCb)
}

func TestTxRollbackCbNoError(t *testing.T) {
	testData := newTxTestData(t)
	testData.mockDbTx.EXPECT().Rollback(testData.ctx).Return(nil).Once()

	calledCb := false
	testData.sut.AddRollbackCallback(func(tx entities.Tx, err error) {
		require.NoError(t, err)
		calledCb = true
	})
	err := testData.sut.Rollback(testData.ctx)
	require.NoError(t, err)
	require.True(t, calledCb)
}

func TestTxRollbackCbError(t *testing.T) {
	testData := newTxTestData(t)
	testData.mockDbTx.EXPECT().Rollback(testData.ctx).Return(testData.testError).Once()

	calledCb := false
	testData.sut.AddRollbackCallback(func(tx entities.Tx, err error) {
		require.ErrorIs(t, err, testData.testError)
		calledCb = true
	})
	err := testData.sut.Rollback(testData.ctx)
	require.ErrorIs(t, err, testData.testError)
	require.True(t, calledCb)
}
