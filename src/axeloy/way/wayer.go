package way

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/core"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/model"
)

type Way interface {
	core.Identitable
	//ValidateProfile(ctx context.Context, profile profile.Profile) error
	SetParams(params model.Params)
	//GetRequiredFields() []string
}

type Sender interface {
	Way
	Send(ctx context.Context, profile profile.Profile, message message.Message) (message.State, error)
}

type Listener interface {
	Way
	Listen(context.Context, func(ctx context.Context, message message.Message) error) error
}

type Wayer interface {
	GetAvailableListeners(ctx context.Context) ([]Listener, error)
	GetSenderByName(ctx context.Context, name string) (Sender, error)
}
