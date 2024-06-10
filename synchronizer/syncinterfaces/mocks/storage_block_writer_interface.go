// Code generated by mockery. DO NOT EDIT.

package mock_syncinterfaces

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"
)

// StorageBlockWriterInterface is an autogenerated mock type for the StorageBlockWriterInterface type
type StorageBlockWriterInterface struct {
	mock.Mock
}

type StorageBlockWriterInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageBlockWriterInterface) EXPECT() *StorageBlockWriterInterface_Expecter {
	return &StorageBlockWriterInterface_Expecter{mock: &_m.Mock}
}

// AddBlock provides a mock function with given fields: ctx, block, dbTx
func (_m *StorageBlockWriterInterface) AddBlock(ctx context.Context, block *entities.L1Block, dbTx entities.Tx) error {
	ret := _m.Called(ctx, block, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for AddBlock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.L1Block, entities.Tx) error); ok {
		r0 = rf(ctx, block, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageBlockWriterInterface_AddBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddBlock'
type StorageBlockWriterInterface_AddBlock_Call struct {
	*mock.Call
}

// AddBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - block *entities.L1Block
//   - dbTx entities.Tx
func (_e *StorageBlockWriterInterface_Expecter) AddBlock(ctx interface{}, block interface{}, dbTx interface{}) *StorageBlockWriterInterface_AddBlock_Call {
	return &StorageBlockWriterInterface_AddBlock_Call{Call: _e.mock.On("AddBlock", ctx, block, dbTx)}
}

func (_c *StorageBlockWriterInterface_AddBlock_Call) Run(run func(ctx context.Context, block *entities.L1Block, dbTx entities.Tx)) *StorageBlockWriterInterface_AddBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.L1Block), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageBlockWriterInterface_AddBlock_Call) Return(_a0 error) *StorageBlockWriterInterface_AddBlock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageBlockWriterInterface_AddBlock_Call) RunAndReturn(run func(context.Context, *entities.L1Block, entities.Tx) error) *StorageBlockWriterInterface_AddBlock_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCheckedBlockByNumber provides a mock function with given fields: ctx, blockNumber, newCheckedStatus, dbTx
func (_m *StorageBlockWriterInterface) UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx entities.Tx) error {
	ret := _m.Called(ctx, blockNumber, newCheckedStatus, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCheckedBlockByNumber")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, bool, entities.Tx) error); ok {
		r0 = rf(ctx, blockNumber, newCheckedStatus, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCheckedBlockByNumber'
type StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call struct {
	*mock.Call
}

// UpdateCheckedBlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - blockNumber uint64
//   - newCheckedStatus bool
//   - dbTx entities.Tx
func (_e *StorageBlockWriterInterface_Expecter) UpdateCheckedBlockByNumber(ctx interface{}, blockNumber interface{}, newCheckedStatus interface{}, dbTx interface{}) *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call {
	return &StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call{Call: _e.mock.On("UpdateCheckedBlockByNumber", ctx, blockNumber, newCheckedStatus, dbTx)}
}

func (_c *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call) Run(run func(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx entities.Tx)) *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(bool), args[3].(entities.Tx))
	})
	return _c
}

func (_c *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call) Return(_a0 error) *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call) RunAndReturn(run func(context.Context, uint64, bool, entities.Tx) error) *StorageBlockWriterInterface_UpdateCheckedBlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorageBlockWriterInterface creates a new instance of StorageBlockWriterInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageBlockWriterInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageBlockWriterInterface {
	mock := &StorageBlockWriterInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}