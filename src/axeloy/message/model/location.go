package model

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
)

type Location struct {
	Way     string
	Profile profile.Profile
}

func (l *Location) GetWay() string {
	return l.Way
}

func (l *Location) GetProfile() profile.Profile {
	return l.Profile
}
