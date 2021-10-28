package model

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/google/uuid"
)

type Route struct {
	Id          uuid.UUID
	Source      profile.Profile
	Destination profile.Profile
	WaysIds     []uuid.UUID
}

func (r *Route) GetId() uuid.UUID {
	return r.Id
}

func (r *Route) GetWaysIds() []uuid.UUID {
	return r.WaysIds
}

func (r Route) GetDestination() profile.Profile {
	return r.Destination
}
func (r Route) GetSource() profile.Profile {
	return r.Source
}
