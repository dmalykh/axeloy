package model

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/google/uuid"
)

type Route struct {
	Id          uuid.UUID
	Source      profile.Profile
	Destination profile.Profile
	Senders     []way.Sender
}

func (r *Route) GetId() uuid.UUID {
	return r.Id
}

func (r *Route) GetSenders() []way.Sender {
	return r.Senders
}

func (r Route) GetDestination() profile.Profile {
	return r.Destination
}
func (r Route) GetSource() profile.Profile {
	return r.Source
}
