// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	profile "github.com/dmalykh/axeloy/axeloy/profile"
	mock "github.com/stretchr/testify/mock"
)

// Profile is an autogenerated mock type for the Profile type
type Profile struct {
	mock.Mock
}

// GetFields provides a mock function with given fields:
func (_m *Profile) GetFields() profile.Fields {
	ret := _m.Called()

	var r0 profile.Fields
	if rf, ok := ret.Get(0).(func() profile.Fields); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(profile.Fields)
		}
	}

	return r0
}
