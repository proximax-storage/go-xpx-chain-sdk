// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	sdk "github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	mock "github.com/stretchr/testify/mock"
)

// ReceiptMapperFn is an autogenerated mock type for the ReceiptMapperFn type
type ReceiptMapperFn struct {
	mock.Mock
}

// Execute provides a mock function with given fields: m
func (_m *ReceiptMapperFn) Execute(m []byte) (*sdk.AnonymousReceipt, error) {
	ret := _m.Called(m)

	var r0 *sdk.AnonymousReceipt
	if rf, ok := ret.Get(0).(func([]byte) *sdk.AnonymousReceipt); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sdk.AnonymousReceipt)
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

type mockConstructorTestingTNewReceiptMapperFn interface {
	mock.TestingT
	Cleanup(func())
}

// NewReceiptMapperFn creates a new instance of ReceiptMapperFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReceiptMapperFn(t mockConstructorTestingTNewReceiptMapperFn) *ReceiptMapperFn {
	mock := &ReceiptMapperFn{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
