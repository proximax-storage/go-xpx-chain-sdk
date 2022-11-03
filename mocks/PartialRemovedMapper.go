// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// PartialRemovedMapper is an autogenerated mock type for the PartialRemovedMapper type
type PartialRemovedMapper struct {
	mock.Mock
}

// MapPartialRemoved provides a mock function with given fields: m
func (_m *PartialRemovedMapper) MapPartialRemoved(m []byte) (*sdk.PartialRemovedInfo, error) {
	ret := _m.Called(m)

	var r0 *sdk.PartialRemovedInfo
	if rf, ok := ret.Get(0).(func([]byte) *sdk.PartialRemovedInfo); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sdk.PartialRemovedInfo)
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

type mockConstructorTestingTNewPartialRemovedMapper interface {
	mock.TestingT
	Cleanup(func())
}

// NewPartialRemovedMapper creates a new instance of PartialRemovedMapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPartialRemovedMapper(t mockConstructorTestingTNewPartialRemovedMapper) *PartialRemovedMapper {
	mock := &PartialRemovedMapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
