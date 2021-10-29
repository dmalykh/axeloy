// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	message "github.com/dmalykh/axeloy/axeloy/message"
	mock "github.com/stretchr/testify/mock"

	profile "github.com/dmalykh/axeloy/axeloy/profile"

	router "github.com/dmalykh/axeloy/axeloy/router"

	way "github.com/dmalykh/axeloy/axeloy/way"
)

// Router is an autogenerated mock type for the Router type
type Router struct {
	mock.Mock
}

// ApplyRoute provides a mock function with given fields: ctx, source, destination, senders
func (_m *Router) ApplyRoute(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error {
	_va := make([]interface{}, len(senders))
	for _i := range senders {
		_va[_i] = senders[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, source, destination)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, profile.Profile, profile.Profile, ...way.Sender) error); ok {
		r0 = rf(ctx, source, destination, senders...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDestinations provides a mock function with given fields: ctx, m
func (_m *Router) GetDestinations(ctx context.Context, m message.Message) ([]router.Destination, error) {
	ret := _m.Called(ctx, m)

	var r0 []router.Destination
	if rf, ok := ret.Get(0).(func(context.Context, message.Message) []router.Destination); ok {
		r0 = rf(ctx, m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]router.Destination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, message.Message) error); ok {
		r1 = rf(ctx, m)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
