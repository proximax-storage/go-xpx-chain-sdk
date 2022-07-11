// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	subscribers "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	mock "github.com/stretchr/testify/mock"
)

// Status is an autogenerated mock type for the Status type
type Status struct {
	mock.Mock
}

// AddHandlers provides a mock function with given fields: handle, handlers
func (_m *Status) AddHandlers(handle *sdk.TransactionChannelHandle, handlers ...subscribers.StatusHandler) error {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sdk.TransactionChannelHandle, ...subscribers.StatusHandler) error); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHandlers provides a mock function with given fields: handle
func (_m *Status) GetHandlers(handle *sdk.TransactionChannelHandle) []*subscribers.StatusHandler {
	ret := _m.Called(handle)

	var r0 []*subscribers.StatusHandler
	if rf, ok := ret.Get(0).(func(*sdk.TransactionChannelHandle) []*subscribers.StatusHandler); ok {
		r0 = rf(handle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*subscribers.StatusHandler)
		}
	}

	return r0
}

// GetHandles provides a mock function with given fields:
func (_m *Status) GetHandles() []string {
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
func (_m *Status) HasHandlers(handle *sdk.TransactionChannelHandle) bool {
	ret := _m.Called(handle)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.TransactionChannelHandle) bool); ok {
		r0 = rf(handle)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemoveHandlers provides a mock function with given fields: handle, handlers
func (_m *Status) RemoveHandlers(handle *sdk.TransactionChannelHandle, handlers ...*subscribers.StatusHandler) bool {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.TransactionChannelHandle, ...*subscribers.StatusHandler) bool); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
