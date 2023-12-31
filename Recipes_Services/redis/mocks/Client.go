// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"
	redis "khanhanhtr/sample/redis"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *Client) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type Client_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *Client_Expecter) Close() *Client_Close_Call {
	return &Client_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *Client_Close_Call) Run(run func()) *Client_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_Close_Call) Return(_a0 error) *Client_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Close_Call) RunAndReturn(run func() error) *Client_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Del provides a mock function with given fields: ctx, keys
func (_m *Client) Del(ctx context.Context, keys ...string) redis.IntegerResult {
	_va := make([]interface{}, len(keys))
	for _i := range keys {
		_va[_i] = keys[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 redis.IntegerResult
	if rf, ok := ret.Get(0).(func(context.Context, ...string) redis.IntegerResult); ok {
		r0 = rf(ctx, keys...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.IntegerResult)
		}
	}

	return r0
}

// Client_Del_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Del'
type Client_Del_Call struct {
	*mock.Call
}

// Del is a helper method to define mock.On call
//   - ctx context.Context
//   - keys ...string
func (_e *Client_Expecter) Del(ctx interface{}, keys ...interface{}) *Client_Del_Call {
	return &Client_Del_Call{Call: _e.mock.On("Del",
		append([]interface{}{ctx}, keys...)...)}
}

func (_c *Client_Del_Call) Run(run func(ctx context.Context, keys ...string)) *Client_Del_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *Client_Del_Call) Return(_a0 redis.IntegerResult) *Client_Del_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Del_Call) RunAndReturn(run func(context.Context, ...string) redis.IntegerResult) *Client_Del_Call {
	_c.Call.Return(run)
	return _c
}

// Expire provides a mock function with given fields: ctx, key, expiration
func (_m *Client) Expire(ctx context.Context, key string, expiration time.Duration) redis.BoolResult {
	ret := _m.Called(ctx, key, expiration)

	var r0 redis.BoolResult
	if rf, ok := ret.Get(0).(func(context.Context, string, time.Duration) redis.BoolResult); ok {
		r0 = rf(ctx, key, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.BoolResult)
		}
	}

	return r0
}

// Client_Expire_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Expire'
type Client_Expire_Call struct {
	*mock.Call
}

// Expire is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - expiration time.Duration
func (_e *Client_Expecter) Expire(ctx interface{}, key interface{}, expiration interface{}) *Client_Expire_Call {
	return &Client_Expire_Call{Call: _e.mock.On("Expire", ctx, key, expiration)}
}

func (_c *Client_Expire_Call) Run(run func(ctx context.Context, key string, expiration time.Duration)) *Client_Expire_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(time.Duration))
	})
	return _c
}

func (_c *Client_Expire_Call) Return(_a0 redis.BoolResult) *Client_Expire_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Expire_Call) RunAndReturn(run func(context.Context, string, time.Duration) redis.BoolResult) *Client_Expire_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key
func (_m *Client) Get(ctx context.Context, key string) redis.StringResult {
	ret := _m.Called(ctx, key)

	var r0 redis.StringResult
	if rf, ok := ret.Get(0).(func(context.Context, string) redis.StringResult); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.StringResult)
		}
	}

	return r0
}

// Client_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type Client_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *Client_Expecter) Get(ctx interface{}, key interface{}) *Client_Get_Call {
	return &Client_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *Client_Get_Call) Run(run func(ctx context.Context, key string)) *Client_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Client_Get_Call) Return(_a0 redis.StringResult) *Client_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Get_Call) RunAndReturn(run func(context.Context, string) redis.StringResult) *Client_Get_Call {
	_c.Call.Return(run)
	return _c
}

// IncrBy provides a mock function with given fields: ctx, key, value
func (_m *Client) IncrBy(ctx context.Context, key string, value int64) redis.IntegerResult {
	ret := _m.Called(ctx, key, value)

	var r0 redis.IntegerResult
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) redis.IntegerResult); ok {
		r0 = rf(ctx, key, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.IntegerResult)
		}
	}

	return r0
}

// Client_IncrBy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IncrBy'
type Client_IncrBy_Call struct {
	*mock.Call
}

// IncrBy is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value int64
func (_e *Client_Expecter) IncrBy(ctx interface{}, key interface{}, value interface{}) *Client_IncrBy_Call {
	return &Client_IncrBy_Call{Call: _e.mock.On("IncrBy", ctx, key, value)}
}

func (_c *Client_IncrBy_Call) Run(run func(ctx context.Context, key string, value int64)) *Client_IncrBy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int64))
	})
	return _c
}

func (_c *Client_IncrBy_Call) Return(_a0 redis.IntegerResult) *Client_IncrBy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_IncrBy_Call) RunAndReturn(run func(context.Context, string, int64) redis.IntegerResult) *Client_IncrBy_Call {
	_c.Call.Return(run)
	return _c
}

// LRange provides a mock function with given fields: ctx, key, start, stop
func (_m *Client) LRange(ctx context.Context, key string, start int64, stop int64) redis.SliceStringResult {
	ret := _m.Called(ctx, key, start, stop)

	var r0 redis.SliceStringResult
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64) redis.SliceStringResult); ok {
		r0 = rf(ctx, key, start, stop)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.SliceStringResult)
		}
	}

	return r0
}

// Client_LRange_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LRange'
type Client_LRange_Call struct {
	*mock.Call
}

// LRange is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - start int64
//   - stop int64
func (_e *Client_Expecter) LRange(ctx interface{}, key interface{}, start interface{}, stop interface{}) *Client_LRange_Call {
	return &Client_LRange_Call{Call: _e.mock.On("LRange", ctx, key, start, stop)}
}

func (_c *Client_LRange_Call) Run(run func(ctx context.Context, key string, start int64, stop int64)) *Client_LRange_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(int64), args[3].(int64))
	})
	return _c
}

func (_c *Client_LRange_Call) Return(_a0 redis.SliceStringResult) *Client_LRange_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_LRange_Call) RunAndReturn(run func(context.Context, string, int64, int64) redis.SliceStringResult) *Client_LRange_Call {
	_c.Call.Return(run)
	return _c
}

// Llen provides a mock function with given fields: ctx, key
func (_m *Client) Llen(ctx context.Context, key string) redis.IntegerResult {
	ret := _m.Called(ctx, key)

	var r0 redis.IntegerResult
	if rf, ok := ret.Get(0).(func(context.Context, string) redis.IntegerResult); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.IntegerResult)
		}
	}

	return r0
}

// Client_Llen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Llen'
type Client_Llen_Call struct {
	*mock.Call
}

// Llen is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
func (_e *Client_Expecter) Llen(ctx interface{}, key interface{}) *Client_Llen_Call {
	return &Client_Llen_Call{Call: _e.mock.On("Llen", ctx, key)}
}

func (_c *Client_Llen_Call) Run(run func(ctx context.Context, key string)) *Client_Llen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Client_Llen_Call) Return(_a0 redis.IntegerResult) *Client_Llen_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Llen_Call) RunAndReturn(run func(context.Context, string) redis.IntegerResult) *Client_Llen_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: ctx
func (_m *Client) Ping(ctx context.Context) redis.StatusResult {
	ret := _m.Called(ctx)

	var r0 redis.StatusResult
	if rf, ok := ret.Get(0).(func(context.Context) redis.StatusResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.StatusResult)
		}
	}

	return r0
}

// Client_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type Client_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Client_Expecter) Ping(ctx interface{}) *Client_Ping_Call {
	return &Client_Ping_Call{Call: _e.mock.On("Ping", ctx)}
}

func (_c *Client_Ping_Call) Run(run func(ctx context.Context)) *Client_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Client_Ping_Call) Return(_a0 redis.StatusResult) *Client_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Ping_Call) RunAndReturn(run func(context.Context) redis.StatusResult) *Client_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// Rpush provides a mock function with given fields: ctx, key, values
func (_m *Client) Rpush(ctx context.Context, key string, values ...interface{}) redis.IntegerResult {
	var _ca []interface{}
	_ca = append(_ca, ctx, key)
	_ca = append(_ca, values...)
	ret := _m.Called(_ca...)

	var r0 redis.IntegerResult
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) redis.IntegerResult); ok {
		r0 = rf(ctx, key, values...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.IntegerResult)
		}
	}

	return r0
}

// Client_Rpush_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rpush'
type Client_Rpush_Call struct {
	*mock.Call
}

// Rpush is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - values ...interface{}
func (_e *Client_Expecter) Rpush(ctx interface{}, key interface{}, values ...interface{}) *Client_Rpush_Call {
	return &Client_Rpush_Call{Call: _e.mock.On("Rpush",
		append([]interface{}{ctx, key}, values...)...)}
}

func (_c *Client_Rpush_Call) Run(run func(ctx context.Context, key string, values ...interface{})) *Client_Rpush_Call {
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

func (_c *Client_Rpush_Call) Return(_a0 redis.IntegerResult) *Client_Rpush_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Rpush_Call) RunAndReturn(run func(context.Context, string, ...interface{}) redis.IntegerResult) *Client_Rpush_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value, expiration
func (_m *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redis.StatusResult {
	ret := _m.Called(ctx, key, value, expiration)

	var r0 redis.StatusResult
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) redis.StatusResult); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(redis.StatusResult)
		}
	}

	return r0
}

// Client_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type Client_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
//   - expiration time.Duration
func (_e *Client_Expecter) Set(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *Client_Set_Call {
	return &Client_Set_Call{Call: _e.mock.On("Set", ctx, key, value, expiration)}
}

func (_c *Client_Set_Call) Run(run func(ctx context.Context, key string, value interface{}, expiration time.Duration)) *Client_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(time.Duration))
	})
	return _c
}

func (_c *Client_Set_Call) Return(_a0 redis.StatusResult) *Client_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Set_Call) RunAndReturn(run func(context.Context, string, interface{}, time.Duration) redis.StatusResult) *Client_Set_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
