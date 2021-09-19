package message

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/google/uuid"
)

type State string

const (
	New             State = "new"
	NotValidProfile State = "not_valid_profile"
	Sent            State = "sent"
	Error           State = "error"
	NoDestinations  State = "no_destinations"
)

type Message interface {
	GetUUID() uuid.UUID
	GetSource() Payload
	GetDestinations() []Payload
}

type Messager interface {
	Save(ctx context.Context, m Message) error
	SaveState(ctx context.Context, m Message, state State, info ...string) error
	//SaveDestinations(ctx context.Context, destinations []router.Destination) error
}

type Payload interface {
	GetWays() []string
	GetProfile() profile.Profile
}
