// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	subscribers "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	mock "github.com/stretchr/testify/mock"
)

// UnconfirmedAdded is an autogenerated mock type for the UnconfirmedAdded type
type UnconfirmedAdded struct {
	mock.Mock
}

// AddHandlers provides a mock function with given fields: handle, handlers
func (_m *UnconfirmedAdded) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...subscribers.UnconfirmedAddedHandler) error {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...subscribers.UnconfirmedAddedHandler) error); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHandlers provides a mock function with given fields: handle
func (_m *UnconfirmedAdded) GetHandlers(handle *sdk.CompoundChannelHandle) []*subscribers.UnconfirmedAddedHandler {
	ret := _m.Called(handle)

	var r0 []*subscribers.UnconfirmedAddedHandler
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle) []*subscribers.UnconfirmedAddedHandler); ok {
		r0 = rf(handle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*subscribers.UnconfirmedAddedHandler)
		}
	}

	return r0
}

// GetHandles provides a mock function with given fields:
func (_m *UnconfirmedAdded) GetHandles() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// HasHandlers provides a mock function with given fields: handle
func (_m *UnconfirmedAdded) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
	ret := _m.Called(handle)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle) bool); ok {
		r0 = rf(handle)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemoveHandlers provides a mock function with given fields: handle, handlers
func (_m *UnconfirmedAdded) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*subscribers.UnconfirmedAddedHandler) bool {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...*subscribers.UnconfirmedAddedHandler) bool); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewUnconfirmedAdded interface {
	mock.TestingT
	Cleanup(func())
}

// NewUnconfirmedAdded creates a new instance of UnconfirmedAdded. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUnconfirmedAdded(t mockConstructorTestingTNewUnconfirmedAdded) *UnconfirmedAdded {
	mock := &UnconfirmedAdded{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
