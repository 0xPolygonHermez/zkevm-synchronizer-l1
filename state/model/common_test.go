package model_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock_entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities/mocks"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/stretchr/testify/require"
)

type testValueStruct struct {
	value uint64
}

func (t *testValueStruct) IsEqual(o interface{}) bool {
	return t.value == o.(*testValueStruct).value
}

func (t *testValueStruct) Key() uint64 {
	return t.value
}

func TestSetStorageHelper(t *testing.T) {
	ctx := context.TODO()
	obj := testValueStruct{value: 1}
	tx := mock_entities.NewTx(t)

	err := model.SetStorageHelper[*testValueStruct](ctx, &obj, tx,
		func(ctx context.Context, obj *testValueStruct, tx entities.Tx) error {
			return nil
		},
		func(ctx context.Context, key uint64, tx entities.Tx) (*testValueStruct, error) {
			return &testValueStruct{value: key}, nil
		},
	)
	require.NoError(t, err)
}

func TestSetStorageHelperCheckObject(t *testing.T) {
	ctx := context.TODO()
	obj := testValueStruct{value: 1}
	tx := mock_entities.NewTx(t)

	err := model.SetStorageHelper[*testValueStruct](ctx, &obj, tx,
		func(ctx context.Context, obj *testValueStruct, tx entities.Tx) error {
			return entities.ErrAlreadyExists
		},
		func(ctx context.Context, key uint64, tx entities.Tx) (*testValueStruct, error) {
			return &testValueStruct{value: key}, nil
		},
	)
	require.NoError(t, err)
}
