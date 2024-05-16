// Code generated by mockery. DO NOT EDIT.

package mock_synchronizer

import (
	context "context"

	synchronizer "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
	mock "github.com/stretchr/testify/mock"
)

// SynchronizerBlockQuerier is an autogenerated mock type for the SynchronizerBlockQuerier type
type SynchronizerBlockQuerier struct {
	mock.Mock
}

type SynchronizerBlockQuerier_Expecter struct {
	mock *mock.Mock
}

func (_m *SynchronizerBlockQuerier) EXPECT() *SynchronizerBlockQuerier_Expecter {
	return &SynchronizerBlockQuerier_Expecter{mock: &_m.Mock}
}

// GetL1BlockByNumber provides a mock function with given fields: ctx, blockNumber
func (_m *SynchronizerBlockQuerier) GetL1BlockByNumber(ctx context.Context, blockNumber uint64) (*synchronizer.L1Block, error) {
	ret := _m.Called(ctx, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for GetL1BlockByNumber")
	}

	var r0 *synchronizer.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) (*synchronizer.L1Block, error)); ok {
		return rf(ctx, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *synchronizer.L1Block); ok {
		r0 = rf(ctx, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*synchronizer.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SynchronizerBlockQuerier_GetL1BlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetL1BlockByNumber'
type SynchronizerBlockQuerier_GetL1BlockByNumber_Call struct {
	*mock.Call
}

// GetL1BlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - blockNumber uint64
func (_e *SynchronizerBlockQuerier_Expecter) GetL1BlockByNumber(ctx interface{}, blockNumber interface{}) *SynchronizerBlockQuerier_GetL1BlockByNumber_Call {
	return &SynchronizerBlockQuerier_GetL1BlockByNumber_Call{Call: _e.mock.On("GetL1BlockByNumber", ctx, blockNumber)}
}

func (_c *SynchronizerBlockQuerier_GetL1BlockByNumber_Call) Run(run func(ctx context.Context, blockNumber uint64)) *SynchronizerBlockQuerier_GetL1BlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64))
	})
	return _c
}

func (_c *SynchronizerBlockQuerier_GetL1BlockByNumber_Call) Return(_a0 *synchronizer.L1Block, _a1 error) *SynchronizerBlockQuerier_GetL1BlockByNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SynchronizerBlockQuerier_GetL1BlockByNumber_Call) RunAndReturn(run func(context.Context, uint64) (*synchronizer.L1Block, error)) *SynchronizerBlockQuerier_GetL1BlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewSynchronizerBlockQuerier creates a new instance of SynchronizerBlockQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSynchronizerBlockQuerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *SynchronizerBlockQuerier {
	mock := &SynchronizerBlockQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
