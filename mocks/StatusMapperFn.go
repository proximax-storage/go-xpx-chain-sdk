// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// StatusMapperFn is an autogenerated mock type for the StatusMapperFn type
type StatusMapperFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: m
func (_m *StatusMapperFn) Execute(m []byte) (*sdk.StatusInfo, error) {
	ret := _m.Called(m)

	var r0 *sdk.StatusInfo
	if rf, ok := ret.Get(0).(func([]byte) *sdk.StatusInfo); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sdk.StatusInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(m)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewStatusMapperFn interface {
	mock.TestingT
	Cleanup(func())
}

// NewStatusMapperFn creates a new instance of StatusMapperFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStatusMapperFn(t mockConstructorTestingTNewStatusMapperFn) *StatusMapperFn {
	mock := &StatusMapperFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}