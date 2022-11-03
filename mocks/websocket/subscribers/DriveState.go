// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	subscribers "github.com/proximax-storage/go-xpx-chain-sdk/sdk/websocket/subscribers"
	mock "github.com/stretchr/testify/mock"
)

// DriveState is an autogenerated mock type for the DriveState type
type DriveState struct {
	mock.Mock
}

// AddHandlers provides a mock function with given fields: address, handlers
func (_m *DriveState) AddHandlers(address *sdk.Address, handlers ...subscribers.DriveStateHandler) error {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, address)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*sdk.Address, ...subscribers.DriveStateHandler) error); ok {
		r0 = rf(address, handlers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAddresses provides a mock function with given fields:
func (_m *DriveState) GetAddresses() []string {
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

// GetHandlers provides a mock function with given fields: address
func (_m *DriveState) GetHandlers(address *sdk.Address) []*subscribers.DriveStateHandler {
	ret := _m.Called(address)

	var r0 []*subscribers.DriveStateHandler
	if rf, ok := ret.Get(0).(func(*sdk.Address) []*subscribers.DriveStateHandler); ok {
		r0 = rf(address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*subscribers.DriveStateHandler)
		}
	}

	return r0
}

// HasHandlers provides a mock function with given fields: address
func (_m *DriveState) HasHandlers(address *sdk.Address) bool {
	ret := _m.Called(address)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.Address) bool); ok {
		r0 = rf(address)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// RemoveHandlers provides a mock function with given fields: address, handlers
func (_m *DriveState) RemoveHandlers(address *sdk.Address, handlers ...*subscribers.DriveStateHandler) bool {
	_va := make([]interface{}, len(handlers))
	for _i := range handlers {
		_va[_i] = handlers[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, address)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.Address, ...*subscribers.DriveStateHandler) bool); ok {
		r0 = rf(address, handlers...)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewDriveState interface {
	mock.TestingT
	Cleanup(func())
}

// NewDriveState creates a new instance of DriveState. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDriveState(t mockConstructorTestingTNewDriveState) *DriveState {
	mock := &DriveState{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
