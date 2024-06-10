// Code generated by mockery. DO NOT EDIT.

package mock_syncinterfaces

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"

	pgstorage "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
)

// StorageVirtualBatchInterface is an autogenerated mock type for the StorageVirtualBatchInterface type
type StorageVirtualBatchInterface struct {
	mock.Mock
}

type StorageVirtualBatchInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageVirtualBatchInterface) EXPECT() *StorageVirtualBatchInterface_Expecter {
	return &StorageVirtualBatchInterface_Expecter{mock: &_m.Mock}
}

// GetLastestVirtualBatchNumber provides a mock function with given fields: ctx, constrains, dbTx
func (_m *StorageVirtualBatchInterface) GetLastestVirtualBatchNumber(ctx context.Context, constrains *pgstorage.VirtualBatchConstraints, dbTx entities.Tx) (uint64, error) {
	ret := _m.Called(ctx, constrains, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetLastestVirtualBatchNumber")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pgstorage.VirtualBatchConstraints, entities.Tx) (uint64, error)); ok {
		return rf(ctx, constrains, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pgstorage.VirtualBatchConstraints, entities.Tx) uint64); ok {
		r0 = rf(ctx, constrains, dbTx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pgstorage.VirtualBatchConstraints, entities.Tx) error); ok {
		r1 = rf(ctx, constrains, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastestVirtualBatchNumber'
type StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call struct {
	*mock.Call
}

// GetLastestVirtualBatchNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - constrains *pgstorage.VirtualBatchConstraints
//   - dbTx entities.Tx
func (_e *StorageVirtualBatchInterface_Expecter) GetLastestVirtualBatchNumber(ctx interface{}, constrains interface{}, dbTx interface{}) *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call {
	return &StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call{Call: _e.mock.On("GetLastestVirtualBatchNumber", ctx, constrains, dbTx)}
}

func (_c *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call) Run(run func(ctx context.Context, constrains *pgstorage.VirtualBatchConstraints, dbTx entities.Tx)) *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*pgstorage.VirtualBatchConstraints), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call) Return(_a0 uint64, _a1 error) *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call) RunAndReturn(run func(context.Context, *pgstorage.VirtualBatchConstraints, entities.Tx) (uint64, error)) *StorageVirtualBatchInterface_GetLastestVirtualBatchNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetVirtualBatchByBatchNumber provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StorageVirtualBatchInterface) GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx entities.Tx) (*entities.VirtualBatch, error) {
	ret := _m.Called(ctx, batchNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetVirtualBatchByBatchNumber")
	}

	var r0 *entities.VirtualBatch
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) (*entities.VirtualBatch, error)); ok {
		return rf(ctx, batchNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) *entities.VirtualBatch); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.VirtualBatch)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, batchNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVirtualBatchByBatchNumber'
type StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call struct {
	*mock.Call
}

// GetVirtualBatchByBatchNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - batchNumber uint64
//   - dbTx entities.Tx
func (_e *StorageVirtualBatchInterface_Expecter) GetVirtualBatchByBatchNumber(ctx interface{}, batchNumber interface{}, dbTx interface{}) *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call {
	return &StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call{Call: _e.mock.On("GetVirtualBatchByBatchNumber", ctx, batchNumber, dbTx)}
}

func (_c *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call) Run(run func(ctx context.Context, batchNumber uint64, dbTx entities.Tx)) *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call) Return(_a0 *entities.VirtualBatch, _a1 error) *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) (*entities.VirtualBatch, error)) *StorageVirtualBatchInterface_GetVirtualBatchByBatchNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorageVirtualBatchInterface creates a new instance of StorageVirtualBatchInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageVirtualBatchInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageVirtualBatchInterface {
	mock := &StorageVirtualBatchInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}