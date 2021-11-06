package model

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/google/uuid"
)

type Route struct {
	Id          uuid.UUID
	Source      profile.Profile
	Destination profile.Profile
	Ways        []string
}

func (r *Route) GetId() uuid.UUID {
	return r.Id
}

func (r *Route) GetWays() []string {
	return r.Ways
}

func (r Route) GetDestination() profile.Profile {
	return r.Destination
}
func (r Route) GetSource() profile.Profile {
	return r.Source
}
