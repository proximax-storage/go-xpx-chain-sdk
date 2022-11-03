// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// ConfirmedAddedMapper is an autogenerated mock type for the ConfirmedAddedMapper type
type ConfirmedAddedMapper struct {
	mock.Mock
}

// MapConfirmedAdded provides a mock function with given fields: m
func (_m *ConfirmedAddedMapper) MapConfirmedAdded(m []byte) (sdk.Transaction, error) {
	ret := _m.Called(m)

	var r0 sdk.Transaction
	if rf, ok := ret.Get(0).(func([]byte) sdk.Transaction); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sdk.Transaction)
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

type mockConstructorTestingTNewConfirmedAddedMapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfirmedAddedMapper creates a new instance of ConfirmedAddedMapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfirmedAddedMapper(t mockConstructorTestingTNewConfirmedAddedMapper) *ConfirmedAddedMapper {
	mock := &ConfirmedAddedMapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
