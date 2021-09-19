// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	way "github.com/dmalykh/axeloy/axeloy/way"
	mock "github.com/stretchr/testify/mock"
)

// Wayer is an autogenerated mock type for the Wayer type
type Wayer struct {
	mock.Mock
}

// GetAvailableListeners provides a mock function with given fields: ctx
func (_m *Wayer) GetAvailableListeners(ctx context.Context) ([]way.Listener, error) {
	ret := _m.Called(ctx)

	var r0 []way.Listener
	if rf, ok := ret.Get(0).(func(context.Context) []way.Listener); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]way.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSenderByName provides a mock function with given fields: ctx, name
func (_m *Wayer) GetSenderByName(ctx context.Context, name string) (way.Sender, error) {
	ret := _m.Called(ctx, name)

	var r0 way.Sender
	if rf, ok := ret.Get(0).(func(context.Context, string) way.Sender); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(way.Sender)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
