package internal_test

import (
	"context"
	"testing"

	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/internal"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/syncinterfaces/mocks"
	"github.com/stretchr/testify/require"
)

func TestReuseStateTxProvider_CallsToBeginTransactionReturnsSameObject(t *testing.T) {
	txProvider := mock_syncinterfaces.NewStateTxProvider(t)
	tx := mock_entities.NewTx(t)

	ctx := context.TODO()
	sut := internal.NewReuseStateTxProvider(txProvider, tx)
	newTx, err := sut.BeginTransaction(ctx)
	require.NoError(t, err)
	newTx2, err2 := sut.BeginTransaction(ctx)
	require.NoError(t, err2)
	require.Equal(t, newTx, newTx2)
}

func TestReuseStateTxProvider_CanCallCommitIfNoRollback(t *testing.T) {
	txProvider := mock_syncinterfaces.NewStateTxProvider(t)
	tx := mock_entities.NewTx(t)

	ctx := context.TODO()
	sut := internal.NewReuseStateTxProvider(txProvider, tx)

	newTx, err := sut.BeginTransaction(ctx)
	require.NoError(t, err)
	//require.Equal(t, tx, newTx)
	err = newTx.Commit(ctx)
	require.NoError(t, err)
	err = newTx.Commit(ctx)
	require.NoError(t, err, "a real tx will fails calling twice Commit")

	// Calling rollback make it for real
	tx.EXPECT().Rollback(ctx).Return(nil)
	err = newTx.Rollback(ctx)
	require.NoError(t, err)

	// Now Commit must fails because the tx have ben rollbacked
	err = newTx.Commit(ctx)
	require.Error(t, err)

}
