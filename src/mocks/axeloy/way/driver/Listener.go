// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	driver "github.com/dmalykh/axeloy/axeloy/way/driver"
	mock "github.com/stretchr/testify/mock"

	profile "github.com/dmalykh/axeloy/axeloy/profile"
)

// Listener is an autogenerated mock type for the Listener type
type Listener struct {
	mock.Mock
}

// Listen provides a mock function with given fields: _a0, _a1
func (_m *Listener) Listen(_a0 context.Context, _a1 func(context.Context, driver.Message) error) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context, driver.Message) error) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetConfig provides a mock function with given fields: config
func (_m *Listener) SetDriverConfig(config driver.DriverConfig) {
	_m.Called(config)
}

// SetWayParams provides a mock function with given fields: params
func (_m *Listener) SetWayParams(params driver.Params) {
	_m.Called(params)
}

// Stop provides a mock function with given fields:
func (_m *Listener) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateProfile provides a mock function with given fields: ctx, p
func (_m *Listener) ValidateProfile(ctx context.Context, p profile.Profile) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, profile.Profile) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}