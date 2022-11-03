// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// UnconfirmedAddedHandler is an autogenerated mock type for the UnconfirmedAddedHandler type
type UnconfirmedAddedHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *UnconfirmedAddedHandler) Execute(_a0 sdk.Transaction) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(sdk.Transaction) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewUnconfirmedAddedHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewUnconfirmedAddedHandler creates a new instance of UnconfirmedAddedHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUnconfirmedAddedHandler(t mockConstructorTestingTNewUnconfirmedAddedHandler) *UnconfirmedAddedHandler {
	mock := &UnconfirmedAddedHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
