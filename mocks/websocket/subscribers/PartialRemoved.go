// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	subscribers "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	mock "github.com/stretchr/testify/mock"
)

// PartialRemoved is an autogenerated mock type for the PartialRemoved type
type PartialRemoved struct {
	mock.Mock
}

// AddHandlers provides a mock function with given fields: handle, handlers
func (_m *PartialRemoved) AddHandlers(handle *sdk.CompoundChannelHandle, handlers ...subscribers.PartialRemovedHandler) error {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...subscribers.PartialRemovedHandler) error); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetHandlers provides a mock function with given fields: handle
func (_m *PartialRemoved) GetHandlers(handle *sdk.CompoundChannelHandle) []*subscribers.PartialRemovedHandler {
	ret := _m.Called(handle)

	var r0 []*subscribers.PartialRemovedHandler
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle) []*subscribers.PartialRemovedHandler); ok {
		r0 = rf(handle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*subscribers.PartialRemovedHandler)
		}
	}

	return r0
}

// GetHandles provides a mock function with given fields:
func (_m *PartialRemoved) GetHandles() []string {
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
func (_m *PartialRemoved) HasHandlers(handle *sdk.CompoundChannelHandle) bool {
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
func (_m *PartialRemoved) RemoveHandlers(handle *sdk.CompoundChannelHandle, handlers ...*subscribers.PartialRemovedHandler) bool {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, handle)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.CompoundChannelHandle, ...*subscribers.PartialRemovedHandler) bool); ok {
		r0 = rf(handle, handlers...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewPartialRemoved interface {
	mock.TestingT
	Cleanup(func())
}

// NewPartialRemoved creates a new instance of PartialRemoved. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPartialRemoved(t mockConstructorTestingTNewPartialRemoved) *PartialRemoved {
	mock := &PartialRemoved{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
