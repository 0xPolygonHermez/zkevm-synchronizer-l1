// Code generated by mockery. DO NOT EDIT.

package mock_storage

import (
	context "context"

	entities "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	mock "github.com/stretchr/testify/mock"
)

// BlockStorer is an autogenerated mock type for the BlockStorer type
type BlockStorer struct {
	mock.Mock
}

type BlockStorer_Expecter struct {
	mock *mock.Mock
}

func (_m *BlockStorer) EXPECT() *BlockStorer_Expecter {
	return &BlockStorer_Expecter{mock: &_m.Mock}
}

// AddBlock provides a mock function with given fields: ctx, block, dbTx
func (_m *BlockStorer) AddBlock(ctx context.Context, block *entities.L1Block, dbTx entities.Tx) error {
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

// BlockStorer_AddBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddBlock'
type BlockStorer_AddBlock_Call struct {
	*mock.Call
}

// AddBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - block *entities.L1Block
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) AddBlock(ctx interface{}, block interface{}, dbTx interface{}) *BlockStorer_AddBlock_Call {
	return &BlockStorer_AddBlock_Call{Call: _e.mock.On("AddBlock", ctx, block, dbTx)}
}

func (_c *BlockStorer_AddBlock_Call) Run(run func(ctx context.Context, block *entities.L1Block, dbTx entities.Tx)) *BlockStorer_AddBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.L1Block), args[2].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_AddBlock_Call) Return(_a0 error) *BlockStorer_AddBlock_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BlockStorer_AddBlock_Call) RunAndReturn(run func(context.Context, *entities.L1Block, entities.Tx) error) *BlockStorer_AddBlock_Call {
	_c.Call.Return(run)
	return _c
}

// GetBlockByNumber provides a mock function with given fields: ctx, blockNumber, dbTx
func (_m *BlockStorer) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx entities.Tx) (*entities.L1Block, error) {
	ret := _m.Called(ctx, blockNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetBlockByNumber")
	}

	var r0 *entities.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)); ok {
		return rf(ctx, blockNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) *entities.L1Block); ok {
		r0 = rf(ctx, blockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, blockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockStorer_GetBlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBlockByNumber'
type BlockStorer_GetBlockByNumber_Call struct {
	*mock.Call
}

// GetBlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - blockNumber uint64
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) GetBlockByNumber(ctx interface{}, blockNumber interface{}, dbTx interface{}) *BlockStorer_GetBlockByNumber_Call {
	return &BlockStorer_GetBlockByNumber_Call{Call: _e.mock.On("GetBlockByNumber", ctx, blockNumber, dbTx)}
}

func (_c *BlockStorer_GetBlockByNumber_Call) Run(run func(ctx context.Context, blockNumber uint64, dbTx entities.Tx)) *BlockStorer_GetBlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_GetBlockByNumber_Call) Return(_a0 *entities.L1Block, _a1 error) *BlockStorer_GetBlockByNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BlockStorer_GetBlockByNumber_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)) *BlockStorer_GetBlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// GetFirstUncheckedBlock provides a mock function with given fields: ctx, fromBlockNumber, dbTx
func (_m *BlockStorer) GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx entities.Tx) (*entities.L1Block, error) {
	ret := _m.Called(ctx, fromBlockNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetFirstUncheckedBlock")
	}

	var r0 *entities.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)); ok {
		return rf(ctx, fromBlockNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) *entities.L1Block); ok {
		r0 = rf(ctx, fromBlockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, fromBlockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockStorer_GetFirstUncheckedBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFirstUncheckedBlock'
type BlockStorer_GetFirstUncheckedBlock_Call struct {
	*mock.Call
}

// GetFirstUncheckedBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - fromBlockNumber uint64
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) GetFirstUncheckedBlock(ctx interface{}, fromBlockNumber interface{}, dbTx interface{}) *BlockStorer_GetFirstUncheckedBlock_Call {
	return &BlockStorer_GetFirstUncheckedBlock_Call{Call: _e.mock.On("GetFirstUncheckedBlock", ctx, fromBlockNumber, dbTx)}
}

func (_c *BlockStorer_GetFirstUncheckedBlock_Call) Run(run func(ctx context.Context, fromBlockNumber uint64, dbTx entities.Tx)) *BlockStorer_GetFirstUncheckedBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_GetFirstUncheckedBlock_Call) Return(_a0 *entities.L1Block, _a1 error) *BlockStorer_GetFirstUncheckedBlock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BlockStorer_GetFirstUncheckedBlock_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)) *BlockStorer_GetFirstUncheckedBlock_Call {
	_c.Call.Return(run)
	return _c
}

// GetLastBlock provides a mock function with given fields: ctx, dbTx
func (_m *BlockStorer) GetLastBlock(ctx context.Context, dbTx entities.Tx) (*entities.L1Block, error) {
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

// BlockStorer_GetLastBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLastBlock'
type BlockStorer_GetLastBlock_Call struct {
	*mock.Call
}

// GetLastBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) GetLastBlock(ctx interface{}, dbTx interface{}) *BlockStorer_GetLastBlock_Call {
	return &BlockStorer_GetLastBlock_Call{Call: _e.mock.On("GetLastBlock", ctx, dbTx)}
}

func (_c *BlockStorer_GetLastBlock_Call) Run(run func(ctx context.Context, dbTx entities.Tx)) *BlockStorer_GetLastBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_GetLastBlock_Call) Return(_a0 *entities.L1Block, _a1 error) *BlockStorer_GetLastBlock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BlockStorer_GetLastBlock_Call) RunAndReturn(run func(context.Context, entities.Tx) (*entities.L1Block, error)) *BlockStorer_GetLastBlock_Call {
	_c.Call.Return(run)
	return _c
}

// GetPreviousBlock provides a mock function with given fields: ctx, offset, dbTx
func (_m *BlockStorer) GetPreviousBlock(ctx context.Context, offset uint64, dbTx entities.Tx) (*entities.L1Block, error) {
	ret := _m.Called(ctx, offset, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetPreviousBlock")
	}

	var r0 *entities.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)); ok {
		return rf(ctx, offset, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, entities.Tx) *entities.L1Block); ok {
		r0 = rf(ctx, offset, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, offset, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockStorer_GetPreviousBlock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPreviousBlock'
type BlockStorer_GetPreviousBlock_Call struct {
	*mock.Call
}

// GetPreviousBlock is a helper method to define mock.On call
//   - ctx context.Context
//   - offset uint64
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) GetPreviousBlock(ctx interface{}, offset interface{}, dbTx interface{}) *BlockStorer_GetPreviousBlock_Call {
	return &BlockStorer_GetPreviousBlock_Call{Call: _e.mock.On("GetPreviousBlock", ctx, offset, dbTx)}
}

func (_c *BlockStorer_GetPreviousBlock_Call) Run(run func(ctx context.Context, offset uint64, dbTx entities.Tx)) *BlockStorer_GetPreviousBlock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_GetPreviousBlock_Call) Return(_a0 *entities.L1Block, _a1 error) *BlockStorer_GetPreviousBlock_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BlockStorer_GetPreviousBlock_Call) RunAndReturn(run func(context.Context, uint64, entities.Tx) (*entities.L1Block, error)) *BlockStorer_GetPreviousBlock_Call {
	_c.Call.Return(run)
	return _c
}

// GetUncheckedBlocks provides a mock function with given fields: ctx, fromBlockNumber, toBlockNumber, dbTx
func (_m *BlockStorer) GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx entities.Tx) (*[]entities.L1Block, error) {
	ret := _m.Called(ctx, fromBlockNumber, toBlockNumber, dbTx)

	if len(ret) == 0 {
		panic("no return value specified for GetUncheckedBlocks")
	}

	var r0 *[]entities.L1Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, entities.Tx) (*[]entities.L1Block, error)); ok {
		return rf(ctx, fromBlockNumber, toBlockNumber, dbTx)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, entities.Tx) *[]entities.L1Block); ok {
		r0 = rf(ctx, fromBlockNumber, toBlockNumber, dbTx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]entities.L1Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64, entities.Tx) error); ok {
		r1 = rf(ctx, fromBlockNumber, toBlockNumber, dbTx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BlockStorer_GetUncheckedBlocks_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUncheckedBlocks'
type BlockStorer_GetUncheckedBlocks_Call struct {
	*mock.Call
}

// GetUncheckedBlocks is a helper method to define mock.On call
//   - ctx context.Context
//   - fromBlockNumber uint64
//   - toBlockNumber uint64
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) GetUncheckedBlocks(ctx interface{}, fromBlockNumber interface{}, toBlockNumber interface{}, dbTx interface{}) *BlockStorer_GetUncheckedBlocks_Call {
	return &BlockStorer_GetUncheckedBlocks_Call{Call: _e.mock.On("GetUncheckedBlocks", ctx, fromBlockNumber, toBlockNumber, dbTx)}
}

func (_c *BlockStorer_GetUncheckedBlocks_Call) Run(run func(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx entities.Tx)) *BlockStorer_GetUncheckedBlocks_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(uint64), args[3].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_GetUncheckedBlocks_Call) Return(_a0 *[]entities.L1Block, _a1 error) *BlockStorer_GetUncheckedBlocks_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *BlockStorer_GetUncheckedBlocks_Call) RunAndReturn(run func(context.Context, uint64, uint64, entities.Tx) (*[]entities.L1Block, error)) *BlockStorer_GetUncheckedBlocks_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCheckedBlockByNumber provides a mock function with given fields: ctx, blockNumber, newCheckedStatus, dbTx
func (_m *BlockStorer) UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx entities.Tx) error {
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

// BlockStorer_UpdateCheckedBlockByNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCheckedBlockByNumber'
type BlockStorer_UpdateCheckedBlockByNumber_Call struct {
	*mock.Call
}

// UpdateCheckedBlockByNumber is a helper method to define mock.On call
//   - ctx context.Context
//   - blockNumber uint64
//   - newCheckedStatus bool
//   - dbTx entities.Tx
func (_e *BlockStorer_Expecter) UpdateCheckedBlockByNumber(ctx interface{}, blockNumber interface{}, newCheckedStatus interface{}, dbTx interface{}) *BlockStorer_UpdateCheckedBlockByNumber_Call {
	return &BlockStorer_UpdateCheckedBlockByNumber_Call{Call: _e.mock.On("UpdateCheckedBlockByNumber", ctx, blockNumber, newCheckedStatus, dbTx)}
}

func (_c *BlockStorer_UpdateCheckedBlockByNumber_Call) Run(run func(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx entities.Tx)) *BlockStorer_UpdateCheckedBlockByNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint64), args[2].(bool), args[3].(entities.Tx))
	})
	return _c
}

func (_c *BlockStorer_UpdateCheckedBlockByNumber_Call) Return(_a0 error) *BlockStorer_UpdateCheckedBlockByNumber_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *BlockStorer_UpdateCheckedBlockByNumber_Call) RunAndReturn(run func(context.Context, uint64, bool, entities.Tx) error) *BlockStorer_UpdateCheckedBlockByNumber_Call {
	_c.Call.Return(run)
	return _c
}

// NewBlockStorer creates a new instance of BlockStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBlockStorer(t interface {
	mock.TestingT
	Cleanup(func())
}) *BlockStorer {
	mock := &BlockStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
