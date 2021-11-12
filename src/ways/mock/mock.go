package main

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

type MockDriver struct {
}

func (m *MockDriver) ValidateProfile(ctx context.Context, p profile.Profile) (map[string]string, error) {
	panic("implement me")
}

func (m *MockDriver) SetWayParams(params driver.Params) {
	panic("implement me")
}

func (m *MockDriver) SetConfig(config driver.DriverConfig) {
	panic("implement me")
}

func (m *MockDriver) Stop() error {
	panic("implement me")
}

var Driver MockDriver
