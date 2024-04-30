// Code generated by mockery. DO NOT EDIT.

package mock_synchronizer

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// CriticalErrorHandler is an autogenerated mock type for the CriticalErrorHandler type
type CriticalErrorHandler struct {
	mock.Mock
}

type CriticalErrorHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *CriticalErrorHandler) EXPECT() *CriticalErrorHandler_Expecter {
	return &CriticalErrorHandler_Expecter{mock: &_m.Mock}
}

// CriticalError provides a mock function with given fields: ctx, err
func (_m *CriticalErrorHandler) CriticalError(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

// CriticalErrorHandler_CriticalError_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CriticalError'
type CriticalErrorHandler_CriticalError_Call struct {
	*mock.Call
}

// CriticalError is a helper method to define mock.On call
//   - ctx context.Context
//   - err error
func (_e *CriticalErrorHandler_Expecter) CriticalError(ctx interface{}, err interface{}) *CriticalErrorHandler_CriticalError_Call {
	return &CriticalErrorHandler_CriticalError_Call{Call: _e.mock.On("CriticalError", ctx, err)}
}

func (_c *CriticalErrorHandler_CriticalError_Call) Run(run func(ctx context.Context, err error)) *CriticalErrorHandler_CriticalError_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(error))
	})
	return _c
}

func (_c *CriticalErrorHandler_CriticalError_Call) Return() *CriticalErrorHandler_CriticalError_Call {
	_c.Call.Return()
	return _c
}

func (_c *CriticalErrorHandler_CriticalError_Call) RunAndReturn(run func(context.Context, error)) *CriticalErrorHandler_CriticalError_Call {
	_c.Call.Return(run)
	return _c
}

// NewCriticalErrorHandler creates a new instance of CriticalErrorHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCriticalErrorHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *CriticalErrorHandler {
	mock := &CriticalErrorHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}