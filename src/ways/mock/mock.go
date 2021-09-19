package main

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/google/uuid"
)

type MockWay struct {
}

func (m *MockWay) SetParams(params model.Params) {
	panic("implement me")
}

func (m *MockWay) GetId() uuid.UUID {
	panic("implement me")
}

func (m *MockWay) ValidateProfile(ctx context.Context, profile profile.Profile) error {
	panic("implement me")
}

func (m *MockWay) Listen(ctx context.Context, f func(ctx context.Context, message message.Message) error) error {
	panic("implement me")
}

var Mock MockWay
