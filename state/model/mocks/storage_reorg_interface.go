// Code generated by mockery. DO NOT EDIT.

package mock_model

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"
)

// StorageReorgInterface is an autogenerated mock type for the StorageReorgInterface type
type StorageReorgInterface struct {
	mock.Mock
}

type StorageReorgInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageReorgInterface) EXPECT() *StorageReorgInterface_Expecter {
	return &StorageReorgInterface_Expecter{mock: &_m.Mock}
}

// GetLastBlock provides a mock function with given fields: ctx, dbTx
func (_m *StorageReorgInterface) GetLastBlock(ctx context.Context, dbTx entities.Tx) (*entities.L1Block, error) {
	ret := _m.Called(ctx, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetLastBlock")
	}

	var r0 *entities.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Tx) (*entities.L1Block, error)); ok {
		return rf(ctx, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entities.Tx) *entities.L1Block); ok {
		r0 = rf(ctx, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, entities.Tx) error); ok {
		r1 = rf(ctx, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageReorgInterface_GetLastBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastBlock'
type StorageReorgInterface_GetLastBlock_Call struct {
	*mock.Call
}

// GetLastBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - dbTx entities.Tx
func (_e *StorageReorgInterface_Expecter) GetLastBlock(ctx interface{}, dbTx interface{}) *StorageReorgInterface_GetLastBlock_Call {
	return &StorageReorgInterface_GetLastBlock_Call{Call: _e.mock.On("GetLastBlock", ctx, dbTx)}
}

func (_c *StorageReorgInterface_GetLastBlock_Call) Run(run func(ctx context.Context, dbTx entities.Tx)) *StorageReorgInterface_GetLastBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entities.Tx))
	})
	return _c
}

func (_c *StorageReorgInterface_GetLastBlock_Call) Return(_a0 *entities.L1Block, _a1 error) *StorageReorgInterface_GetLastBlock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageReorgInterface_GetLastBlock_Call) RunAndReturn(run func(context.Context, entities.Tx) (*entities.L1Block, error)) *StorageReorgInterface_GetLastBlock_Call {
	_c.Call.Return(run)
	return _c
}

// ResetToL1BlockNumber provides a mock function with given fields: ctx, firstBlockNumberToKeep, dbTx
func (_m *StorageReorgInterface) ResetToL1BlockNumber(ctx context.Context, firstBlockNumberToKeep uint64, dbTx entities.Tx) error {
	ret := _m.Called(ctx, firstBlockNumberToKeep, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for ResetToL1BlockNumber")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) error); ok {
		r0 = rf(ctx, firstBlockNumberToKeep, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageReorgInterface_ResetToL1BlockNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetToL1BlockNumber'
type StorageReorgInterface_ResetToL1BlockNumber_Call struct {
	*mock.Call
}

// ResetToL1BlockNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - firstBlockNumberToKeep uint64
//   - dbTx entities.Tx
func (_e *StorageReorgInterface_Expecter) ResetToL1BlockNumber(ctx interface{}, firstBlockNumberToKeep interface{}, dbTx interface{}) *StorageReorgInterface_ResetToL1BlockNumber_Call {
	return &StorageReorgInterface_ResetToL1BlockNumber_Call{Call: _e.mock.On("ResetToL1BlockNumber", ctx, firstBlockNumberToKeep, dbTx)}
}

func (_c *StorageReorgInterface_ResetToL1BlockNumber_Call) Run(run func(ctx context.Context, firstBlockNumberToKeep uint64, dbTx entities.Tx)) *StorageReorgInterface_ResetToL1BlockNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageReorgInterface_ResetToL1BlockNumber_Call) Return(_a0 error) *StorageReorgInterface_ResetToL1BlockNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageReorgInterface_ResetToL1BlockNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) error) *StorageReorgInterface_ResetToL1BlockNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorageReorgInterface creates a new instance of StorageReorgInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageReorgInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageReorgInterface {
	mock := &StorageReorgInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
