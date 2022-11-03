// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ReceiptBody is an autogenerated mock type for the ReceiptBody type
type ReceiptBody struct {
	mock.Mock
}

// toStruct provides a mock function with given fields:
func (_m *ReceiptBody) toStruct() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

type mockConstructorTestingTNewReceiptBody interface {
	mock.TestingT
	Cleanup(func())
}

// NewReceiptBody creates a new instance of ReceiptBody. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewReceiptBody(t mockConstructorTestingTNewReceiptBody) *ReceiptBody {
	mock := &ReceiptBody{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
