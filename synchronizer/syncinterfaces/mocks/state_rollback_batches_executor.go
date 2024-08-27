// Code generated by mockery. DO NOT EDIT.

package mock_syncinterfaces

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"

	model "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
)

// StateRollbackBatchesExecutor is an autogenerated mock type for the StateRollbackBatchesExecutor type
type StateRollbackBatchesExecutor struct {
	mock.Mock
}

type StateRollbackBatchesExecutor_Expecter struct {
	mock *mock.Mock
}

func (_m *StateRollbackBatchesExecutor) EXPECT() *StateRollbackBatchesExecutor_Expecter {
	return &StateRollbackBatchesExecutor_Expecter{mock: &_m.Mock}
}

// ExecuteRollbackBatches provides a mock function with given fields: ctx, rollbackBatchesRequest, dbTx
func (_m *StateRollbackBatchesExecutor) ExecuteRollbackBatches(ctx context.Context, rollbackBatchesRequest model.RollbackBatchesRequest, dbTx entities.Tx) (*model.RollbackBatchesExecutionResult, error) {
	ret := _m.Called(ctx, rollbackBatchesRequest, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteRollbackBatches")
	}

	var r0 *model.RollbackBatchesExecutionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.RollbackBatchesRequest, entities.Tx) (*model.RollbackBatchesExecutionResult, error)); ok {
		return rf(ctx, rollbackBatchesRequest, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.RollbackBatchesRequest, entities.Tx) *model.RollbackBatchesExecutionResult); ok {
		r0 = rf(ctx, rollbackBatchesRequest, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RollbackBatchesExecutionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.RollbackBatchesRequest, entities.Tx) error); ok {
		r1 = rf(ctx, rollbackBatchesRequest, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteRollbackBatches'
type StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call struct {
	*mock.Call
}

// ExecuteRollbackBatches is a helper method to define mock.On call
//   - ctx context.Context
//   - rollbackBatchesRequest model.RollbackBatchesRequest
//   - dbTx entities.Tx
func (_e *StateRollbackBatchesExecutor_Expecter) ExecuteRollbackBatches(ctx interface{}, rollbackBatchesRequest interface{}, dbTx interface{}) *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call {
	return &StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call{Call: _e.mock.On("ExecuteRollbackBatches", ctx, rollbackBatchesRequest, dbTx)}
}

func (_c *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call) Run(run func(ctx context.Context, rollbackBatchesRequest model.RollbackBatchesRequest, dbTx entities.Tx)) *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.RollbackBatchesRequest), args[2].(entities.Tx))
	})
	return _c
}

func (_c *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call) Return(_a0 *model.RollbackBatchesExecutionResult, _a1 error) *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call) RunAndReturn(run func(context.Context, model.RollbackBatchesRequest, entities.Tx) (*model.RollbackBatchesExecutionResult, error)) *StateRollbackBatchesExecutor_ExecuteRollbackBatches_Call {
	_c.Call.Return(run)
	return _c
}

// NewStateRollbackBatchesExecutor creates a new instance of StateRollbackBatchesExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStateRollbackBatchesExecutor(t interface {
	mock.TestingT
	Cleanup(func())
}) *StateRollbackBatchesExecutor {
	mock := &StateRollbackBatchesExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}