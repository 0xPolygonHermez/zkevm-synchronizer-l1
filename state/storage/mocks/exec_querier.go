// Code generated by mockery. DO NOT EDIT.

package mock_storage

import (
	context "context"

	pgconn "github.com/jackc/pgconn"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v4"
)

// execQuerier is an autogenerated mock type for the execQuerier type
type execQuerier struct {
	mock.Mock
}

type execQuerier_Expecter struct {
	mock *mock.Mock
}

func (_m *execQuerier) EXPECT() *execQuerier_Expecter {
	return &execQuerier_Expecter{mock: &_m.Mock}
}

// CopyFrom provides a mock function with given fields: ctx, tableName, columnNames, rowSrc
func (_m *execQuerier) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	ret := _m.Called(ctx, tableName, columnNames, rowSrc)

	if len(ret) == 0 {
		panic("no return value specified for CopyFrom")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)); ok {
		return rf(ctx, tableName, columnNames, rowSrc)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) int64); ok {
		r0 = rf(ctx, tableName, columnNames, rowSrc)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) error); ok {
		r1 = rf(ctx, tableName, columnNames, rowSrc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// execQuerier_CopyFrom_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CopyFrom'
type execQuerier_CopyFrom_Call struct {
	*mock.Call
}

// CopyFrom is a helper method to define mock.On call
//   - ctx context.Context
//   - tableName pgx.Identifier
//   - columnNames []string
//   - rowSrc pgx.CopyFromSource
func (_e *execQuerier_Expecter) CopyFrom(ctx interface{}, tableName interface{}, columnNames interface{}, rowSrc interface{}) *execQuerier_CopyFrom_Call {
	return &execQuerier_CopyFrom_Call{Call: _e.mock.On("CopyFrom", ctx, tableName, columnNames, rowSrc)}
}

func (_c *execQuerier_CopyFrom_Call) Run(run func(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource)) *execQuerier_CopyFrom_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pgx.Identifier), args[2].([]string), args[3].(pgx.CopyFromSource))
	})
	return _c
}

func (_c *execQuerier_CopyFrom_Call) Return(_a0 int64, _a1 error) *execQuerier_CopyFrom_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *execQuerier_CopyFrom_Call) RunAndReturn(run func(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)) *execQuerier_CopyFrom_Call {
	_c.Call.Return(run)
	return _c
}

// Exec provides a mock function with given fields: ctx, sql, arguments
func (_m *execQuerier) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, arguments...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 pgconn.CommandTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)); ok {
		return rf(ctx, sql, arguments...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgconn.CommandTag); ok {
		r0 = rf(ctx, sql, arguments...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgconn.CommandTag)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, arguments...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// execQuerier_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type execQuerier_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - arguments ...interface{}
func (_e *execQuerier_Expecter) Exec(ctx interface{}, sql interface{}, arguments ...interface{}) *execQuerier_Exec_Call {
	return &execQuerier_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{ctx, sql}, arguments...)...)}
}

func (_c *execQuerier_Exec_Call) Run(run func(ctx context.Context, sql string, arguments ...interface{})) *execQuerier_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *execQuerier_Exec_Call) Return(commandTag pgconn.CommandTag, err error) *execQuerier_Exec_Call {
	_c.Call.Return(commandTag, err)
	return _c
}

func (_c *execQuerier_Exec_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)) *execQuerier_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, sql, args
func (_m *execQuerier) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 pgx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgx.Rows, error)); ok {
		return rf(ctx, sql, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Rows); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// execQuerier_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type execQuerier_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *execQuerier_Expecter) Query(ctx interface{}, sql interface{}, args ...interface{}) *execQuerier_Query_Call {
	return &execQuerier_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *execQuerier_Query_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *execQuerier_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *execQuerier_Query_Call) Return(_a0 pgx.Rows, _a1 error) *execQuerier_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *execQuerier_Query_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgx.Rows, error)) *execQuerier_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: ctx, sql, args
func (_m *execQuerier) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRow")
	}

	var r0 pgx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Row); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Row)
		}
	}

	return r0
}

// execQuerier_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type execQuerier_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - args ...interface{}
func (_e *execQuerier_Expecter) QueryRow(ctx interface{}, sql interface{}, args ...interface{}) *execQuerier_QueryRow_Call {
	return &execQuerier_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *execQuerier_QueryRow_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *execQuerier_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *execQuerier_QueryRow_Call) Return(_a0 pgx.Row) *execQuerier_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *execQuerier_QueryRow_Call) RunAndReturn(run func(context.Context, string, ...interface{}) pgx.Row) *execQuerier_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

// newExecQuerier creates a new instance of execQuerier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newExecQuerier(t interface {
	mock.TestingT
	Cleanup(func())
}) *execQuerier {
	mock := &execQuerier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
