// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// CosignatureHandler is an autogenerated mock type for the CosignatureHandler type
type CosignatureHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *CosignatureHandler) Execute(_a0 *sdk.SignerInfo) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*sdk.SignerInfo) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewCosignatureHandler interface {
	mock.TestingT
	Cleanup(func())
}

// NewCosignatureHandler creates a new instance of CosignatureHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCosignatureHandler(t mockConstructorTestingTNewCosignatureHandler) *CosignatureHandler {
	mock := &CosignatureHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}