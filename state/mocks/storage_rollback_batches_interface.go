// Code generated by mockery. DO NOT EDIT.

package mock_state

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"
)

// StorageRollbackBatchesInterface is an autogenerated mock type for the StorageRollbackBatchesInterface type
type StorageRollbackBatchesInterface struct {
	mock.Mock
}

type StorageRollbackBatchesInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *StorageRollbackBatchesInterface) EXPECT() *StorageRollbackBatchesInterface_Expecter {
	return &StorageRollbackBatchesInterface_Expecter{mock: &_m.Mock}
}

// AddRollbackBatchesLogEntry provides a mock function with given fields: ctx, entry, dbTx
func (_m *StorageRollbackBatchesInterface) AddRollbackBatchesLogEntry(ctx context.Context, entry *entities.RollbackBatchesLogEntry, dbTx entities.Tx) error {
	ret := _m.Called(ctx, entry, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for AddRollbackBatchesLogEntry")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.RollbackBatchesLogEntry, entities.Tx) error); ok {
		r0 = rf(ctx, entry, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddRollbackBatchesLogEntry'
type StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call struct {
	*mock.Call
}

// AddRollbackBatchesLogEntry is a helper method to define mock.On call
//   - ctx context.Context
//   - entry *entities.RollbackBatchesLogEntry
//   - dbTx entities.Tx
func (_e *StorageRollbackBatchesInterface_Expecter) AddRollbackBatchesLogEntry(ctx interface{}, entry interface{}, dbTx interface{}) *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call {
	return &StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call{Call: _e.mock.On("AddRollbackBatchesLogEntry", ctx, entry, dbTx)}
}

func (_c *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call) Run(run func(ctx context.Context, entry *entities.RollbackBatchesLogEntry, dbTx entities.Tx)) *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.RollbackBatchesLogEntry), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call) Return(_a0 error) *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call) RunAndReturn(run func(context.Context, *entities.RollbackBatchesLogEntry, entities.Tx) error) *StorageRollbackBatchesInterface_AddRollbackBatchesLogEntry_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteSequencesGreatestOrEqualBatchNumber provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StorageRollbackBatchesInterface) DeleteSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx entities.Tx) error {
	ret := _m.Called(ctx, batchNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSequencesGreatestOrEqualBatchNumber")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) error); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteSequencesGreatestOrEqualBatchNumber'
type StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call struct {
	*mock.Call
}

// DeleteSequencesGreatestOrEqualBatchNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - batchNumber uint64
//   - dbTx entities.Tx
func (_e *StorageRollbackBatchesInterface_Expecter) DeleteSequencesGreatestOrEqualBatchNumber(ctx interface{}, batchNumber interface{}, dbTx interface{}) *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call {
	return &StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call{Call: _e.mock.On("DeleteSequencesGreatestOrEqualBatchNumber", ctx, batchNumber, dbTx)}
}

func (_c *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call) Run(run func(ctx context.Context, batchNumber uint64, dbTx entities.Tx)) *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call) Return(_a0 error) *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) error) *StorageRollbackBatchesInterface_DeleteSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber provides a mock function with given fields: ctx, l1BlockNumber, dbTx
func (_m *StorageRollbackBatchesInterface) GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx context.Context, l1BlockNumber uint64, dbTx entities.Tx) ([]entities.RollbackBatchesLogEntry, error) {
	ret := _m.Called(ctx, l1BlockNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber")
	}

	var r0 []entities.RollbackBatchesLogEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) ([]entities.RollbackBatchesLogEntry, error)); ok {
		return rf(ctx, l1BlockNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) []entities.RollbackBatchesLogEntry); ok {
		r0 = rf(ctx, l1BlockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entities.RollbackBatchesLogEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, l1BlockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber'
type StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call struct {
	*mock.Call
}

// GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - l1BlockNumber uint64
//   - dbTx entities.Tx
func (_e *StorageRollbackBatchesInterface_Expecter) GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber(ctx interface{}, l1BlockNumber interface{}, dbTx interface{}) *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call {
	return &StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call{Call: _e.mock.On("GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber", ctx, l1BlockNumber, dbTx)}
}

func (_c *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call) Run(run func(ctx context.Context, l1BlockNumber uint64, dbTx entities.Tx)) *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call) Return(_a0 []entities.RollbackBatchesLogEntry, _a1 error) *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) ([]entities.RollbackBatchesLogEntry, error)) *StorageRollbackBatchesInterface_GetRollbackBatchesLogEntryGreaterOrEqualL1BlockNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetSequencesGreatestOrEqualBatchNumber provides a mock function with given fields: ctx, batchNumber, dbTx
func (_m *StorageRollbackBatchesInterface) GetSequencesGreatestOrEqualBatchNumber(ctx context.Context, batchNumber uint64, dbTx entities.Tx) (*entities.SequencesBatchesSlice, error) {
	ret := _m.Called(ctx, batchNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetSequencesGreatestOrEqualBatchNumber")
	}

	var r0 *entities.SequencesBatchesSlice
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) (*entities.SequencesBatchesSlice, error)); ok {
		return rf(ctx, batchNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) *entities.SequencesBatchesSlice); ok {
		r0 = rf(ctx, batchNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.SequencesBatchesSlice)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, batchNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSequencesGreatestOrEqualBatchNumber'
type StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call struct {
	*mock.Call
}

// GetSequencesGreatestOrEqualBatchNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - batchNumber uint64
//   - dbTx entities.Tx
func (_e *StorageRollbackBatchesInterface_Expecter) GetSequencesGreatestOrEqualBatchNumber(ctx interface{}, batchNumber interface{}, dbTx interface{}) *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call {
	return &StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call{Call: _e.mock.On("GetSequencesGreatestOrEqualBatchNumber", ctx, batchNumber, dbTx)}
}

func (_c *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call) Run(run func(ctx context.Context, batchNumber uint64, dbTx entities.Tx)) *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call) Return(_a0 *entities.SequencesBatchesSlice, _a1 error) *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) (*entities.SequencesBatchesSlice, error)) *StorageRollbackBatchesInterface_GetSequencesGreatestOrEqualBatchNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewStorageRollbackBatchesInterface creates a new instance of StorageRollbackBatchesInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageRollbackBatchesInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageRollbackBatchesInterface {
	mock := &StorageRollbackBatchesInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
