package model_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
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

	mockGetter := mock_storage.NewKvStorer(t)
	var metadata *entities.KVMetadataEntry
	// Set up the mock to return the value when GetKV is called
	mockGetter.On("KVGetJson", ctx, key, mock.AnythingOfType("*string"), metadata, dbTx).Return(nil).Run(func(args mock.Arguments) {
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
	mockGetter := mock_storage.NewKvStorer(t)
	var metadata *entities.KVMetadataEntry
	// Set up the mock to return the value when GetKV is called
	mockGetter.On("KVGetJson", ctx, key, mock.AnythingOfType("*string"), metadata, dbTx).Return(entities.ErrNotFound).Run(func(args mock.Arguments) {
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
	mockGetter := mock_storage.NewKvStorer(t)

	returnedError := entities.ErrStorageNotRegister
	var metadata *entities.KVMetadataEntry
	// Set up the mock to return the value when GetKV is called
	mockGetter.On("KVGetJson", ctx, key, mock.AnythingOfType("*string"), metadata, dbTx).Return(returnedError).Run(func(args mock.Arguments) {
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
