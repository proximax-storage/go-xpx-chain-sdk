// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	subscribers "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	mock "github.com/stretchr/testify/mock"
)

// ConfirmedAdded is an autogenerated mock type for the ConfirmedAdded type
type ConfirmedAdded struct {
	mock.Mock
}

// AddHandlers provides a mock function with given fields: handle, handlers
func (_m *ConfirmedAdded) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...subscribers.ConfirmedAddedHandler) error {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...subscribers.ConfirmedAddedHandler) error); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHandlers provides a mock function with given fields: handle
func (_m *ConfirmedAdded) GetHandlers(handle *sdk.CompoundChannelHandle) []*subscribers.ConfirmedAddedHandler {
	ret := _m.Called(handle)

	var r0 []*subscribers.ConfirmedAddedHandler
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle) []*subscribers.ConfirmedAddedHandler); ok {
		r0 = rf(handle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*subscribers.ConfirmedAddedHandler)
		}
	}

	return r0
}

// GetHandles provides a mock function with given fields:
func (_m *ConfirmedAdded) GetHandles() []string {
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
func (_m *ConfirmedAdded) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
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
func (_m *ConfirmedAdded) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*subscribers.ConfirmedAddedHandler) bool {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...*subscribers.ConfirmedAddedHandler) bool); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewConfirmedAdded interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfirmedAdded creates a new instance of ConfirmedAdded. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfirmedAdded(t mockConstructorTestingTNewConfirmedAdded) *ConfirmedAdded {
	mock := &ConfirmedAdded{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
