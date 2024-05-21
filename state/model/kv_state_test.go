package model_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	mock_model "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model/mocks"
	mock_storage "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetKVHelperOk(t *testing.T) {
	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	dbTx := mock_entities.NewTx(t)

	// Create a mock implementation of StateKVInterface
	mockGetter := mock_model.NewStateKVInterface(t)

	// Set up the mock to return the value when GetKV is called
	mockGetter.On("GetKV", ctx, key, mock.AnythingOfType("*string"), dbTx).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*string)
		*arg = value
	})

	// Call the GetKVHelper function
	result, err := model.GetKVHelper[string](ctx, key, mockGetter, dbTx)

	// Assert that the result is equal to the expected value
	assert.Equal(t, &value, result)

	// Assert that there is no error
	require.NoError(t, err)
}

func TestGetKVHelperErrorNotFoundReturnsNil(t *testing.T) {
	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	dbTx := mock_entities.NewTx(t)

	// Create a mock implementation of StateKVInterface
	mockGetter := mock_model.NewStateKVInterface(t)

	// Set up the mock to return the value when GetKV is called
	mockGetter.On("GetKV", ctx, key, mock.AnythingOfType("*string"), dbTx).Return(entities.ErrNotFound).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*string)
		*arg = value
	})

	// Call the GetKVHelper function
	result, err := model.GetKVHelper[string](ctx, key, mockGetter, dbTx)

	// Assert that the result is equal to the expected value
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestGetKVHelperIfErorIsNotNotFoundReturnsTheError(t *testing.T) {
	ctx := context.Background()
	key := "testKey"
	value := "testValue"
	dbTx := mock_entities.NewTx(t)

	// Create a mock implementation of StateKVInterface
	mockGetter := mock_model.NewStateKVInterface(t)
	returnedError := entities.ErrStorageNotRegister
	// Set up the mock to return the value when GetKV is called
	mockGetter.On("GetKV", ctx, key, mock.AnythingOfType("*string"), dbTx).Return(returnedError).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*string)
		*arg = value
	})

	// Call the GetKVHelper function
	result, err := model.GetKVHelper[string](ctx, key, mockGetter, dbTx)

	// Assert that the result is equal to the expected value
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorIs(t, err, returnedError)
}

func TestKVState(t *testing.T) {
	ctx := context.TODO()
	key := "testKey"
	value := "testValue"
	mockStorage := mock_storage.NewKvStorer(t)
	dbTx := mock_entities.NewTx(t)

	sut := model.NewKVState(mockStorage)
	var metadata *entities.KVMetadataEntry
	mockStorage.EXPECT().KVSetJson(ctx, key, value, metadata, dbTx).Return(nil).Once()

	err := sut.SetKV(ctx, key, value, dbTx)
	require.NoError(t, err)

	returnedError := entities.ErrStorageNotRegister
	mockStorage.EXPECT().KVSetJson(ctx, key, value, metadata, dbTx).Return(returnedError).Once()

	err = sut.SetKV(ctx, key, value, dbTx)
	require.ErrorIs(t, err, returnedError)
}
