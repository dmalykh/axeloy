// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	profile "github.com/dmalykh/axeloy/axeloy/profile"
	mock "github.com/stretchr/testify/mock"
)

// Payload is an autogenerated mock type for the Payload type
type Payload struct {
	mock.Mock
}

// GetProfile provides a mock function with given fields:
func (_m *Payload) GetProfile() profile.Profile {
	ret := _m.Called()

	var r0 profile.Profile
	if rf, ok := ret.Get(0).(func() profile.Profile); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(profile.Profile)
		}
	}

	return r0
}

// GetWays provides a mock function with given fields:
func (_m *Payload) GetWays() []string {
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
