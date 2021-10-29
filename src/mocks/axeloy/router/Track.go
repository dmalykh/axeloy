// Code generated by mockery 2.9.0. DO NOT EDIT.

package mocks

import (
	profile "github.com/dmalykh/axeloy/axeloy/profile"
	model "github.com/dmalykh/axeloy/axeloy/router/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Track is an autogenerated mock type for the Track type
type Track struct {
	mock.Mock
}

// GetAttempts provides a mock function with given fields:
func (_m *Track) GetAttempts() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// GetId provides a mock function with given fields:
func (_m *Track) GetId() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// GetMessageId provides a mock function with given fields:
func (_m *Track) GetMessageId() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// GetProfile provides a mock function with given fields:
func (_m *Track) GetProfile() profile.Profile {
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

// GetSenderId provides a mock function with given fields:
func (_m *Track) GetSenderId() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// GetStatus provides a mock function with given fields:
func (_m *Track) GetStatus() model.TrackStatus {
	ret := _m.Called()

	var r0 model.TrackStatus
	if rf, ok := ret.Get(0).(func() model.TrackStatus); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(model.TrackStatus)
	}

	return r0
}
