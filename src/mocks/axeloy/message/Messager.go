// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	context "context"

	message "github.com/dmalykh/axeloy/axeloy/message"
	mock "github.com/stretchr/testify/mock"

	model "github.com/dmalykh/axeloy/axeloy/message/model"

	uuid "github.com/google/uuid"
)

// Messager is an autogenerated mock type for the Messager type
type Messager struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, msg
func (_m *Messager) Add(ctx context.Context, msg message.Message) error {
	ret := _m.Called(ctx, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, message.Message) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetById provides a mock function with given fields: ctx, id
func (_m *Messager) GetById(ctx context.Context, id uuid.UUID) (message.Message, error) {
	ret := _m.Called(ctx, id)

	var r0 message.Message
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) message.Message); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(message.Message)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, msg, status, info
func (_m *Messager) UpdateStatus(ctx context.Context, msg message.Message, status model.Status, info ...string) error {
	_va := make([]interface{}, len(info))
	for _i := range info {
		_va[_i] = info[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, msg, status)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, message.Message, model.Status, ...string) error); ok {
		r0 = rf(ctx, msg, status, info...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
