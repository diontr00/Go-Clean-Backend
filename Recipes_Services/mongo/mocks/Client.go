// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"
	mongo "khanhanhtr/sample/mongo"

	mock "github.com/stretchr/testify/mock"

	mongo_drivermongo "go.mongodb.org/mongo-driver/mongo"
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

// Connect provides a mock function with given fields: ctx
func (_m *Client) Connect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type Client_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Client_Expecter) Connect(ctx interface{}) *Client_Connect_Call {
	return &Client_Connect_Call{Call: _e.mock.On("Connect", ctx)}
}

func (_c *Client_Connect_Call) Run(run func(ctx context.Context)) *Client_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Client_Connect_Call) Return(_a0 error) *Client_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Connect_Call) RunAndReturn(run func(context.Context) error) *Client_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// Disconnect provides a mock function with given fields: ctx
func (_m *Client) Disconnect(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type Client_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Client_Expecter) Disconnect(ctx interface{}) *Client_Disconnect_Call {
	return &Client_Disconnect_Call{Call: _e.mock.On("Disconnect", ctx)}
}

func (_c *Client_Disconnect_Call) Run(run func(ctx context.Context)) *Client_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Client_Disconnect_Call) Return(_a0 error) *Client_Disconnect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Disconnect_Call) RunAndReturn(run func(context.Context) error) *Client_Disconnect_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields: ctx
func (_m *Client) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
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

func (_c *Client_Ping_Call) Return(_a0 error) *Client_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Ping_Call) RunAndReturn(run func(context.Context) error) *Client_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// StartSession provides a mock function with given fields:
func (_m *Client) StartSession() (mongo_drivermongo.Session, error) {
	ret := _m.Called()

	var r0 mongo_drivermongo.Session
	var r1 error
	if rf, ok := ret.Get(0).(func() (mongo_drivermongo.Session, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() mongo_drivermongo.Session); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongo_drivermongo.Session)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_StartSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartSession'
type Client_StartSession_Call struct {
	*mock.Call
}

// StartSession is a helper method to define mock.On call
func (_e *Client_Expecter) StartSession() *Client_StartSession_Call {
	return &Client_StartSession_Call{Call: _e.mock.On("StartSession")}
}

func (_c *Client_StartSession_Call) Run(run func()) *Client_StartSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_StartSession_Call) Return(_a0 mongo_drivermongo.Session, _a1 error) *Client_StartSession_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_StartSession_Call) RunAndReturn(run func() (mongo_drivermongo.Session, error)) *Client_StartSession_Call {
	_c.Call.Return(run)
	return _c
}

// UseDatabase provides a mock function with given fields: _a0
func (_m *Client) UseDatabase(_a0 string) mongo.Database {
	ret := _m.Called(_a0)

	var r0 mongo.Database
	if rf, ok := ret.Get(0).(func(string) mongo.Database); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongo.Database)
		}
	}

	return r0
}

// Client_UseDatabase_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseDatabase'
type Client_UseDatabase_Call struct {
	*mock.Call
}

// UseDatabase is a helper method to define mock.On call
//   - _a0 string
func (_e *Client_Expecter) UseDatabase(_a0 interface{}) *Client_UseDatabase_Call {
	return &Client_UseDatabase_Call{Call: _e.mock.On("UseDatabase", _a0)}
}

func (_c *Client_UseDatabase_Call) Run(run func(_a0 string)) *Client_UseDatabase_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_UseDatabase_Call) Return(_a0 mongo.Database) *Client_UseDatabase_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_UseDatabase_Call) RunAndReturn(run func(string) mongo.Database) *Client_UseDatabase_Call {
	_c.Call.Return(run)
	return _c
}

// UseSession provides a mock function with given fields: ctx, fn
func (_m *Client) UseSession(ctx context.Context, fn func(mongo_drivermongo.SessionContext) error) error {
	ret := _m.Called(ctx, fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(mongo_drivermongo.SessionContext) error) error); ok {
		r0 = rf(ctx, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Client_UseSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseSession'
type Client_UseSession_Call struct {
	*mock.Call
}

// UseSession is a helper method to define mock.On call
//   - ctx context.Context
//   - fn func(mongo_drivermongo.SessionContext) error
func (_e *Client_Expecter) UseSession(ctx interface{}, fn interface{}) *Client_UseSession_Call {
	return &Client_UseSession_Call{Call: _e.mock.On("UseSession", ctx, fn)}
}

func (_c *Client_UseSession_Call) Run(run func(ctx context.Context, fn func(mongo_drivermongo.SessionContext) error)) *Client_UseSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(mongo_drivermongo.SessionContext) error))
	})
	return _c
}

func (_c *Client_UseSession_Call) Return(_a0 error) *Client_UseSession_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_UseSession_Call) RunAndReturn(run func(context.Context, func(mongo_drivermongo.SessionContext) error) error) *Client_UseSession_Call {
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
