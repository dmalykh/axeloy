// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Way is an autogenerated mock type for the Way type
type Way struct {
	mock.Mock
}

// GetName provides a mock function with given fields:
func (_m *Way) GetName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetTitle provides a mock function with given fields:
func (_m *Way) GetTitle() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
