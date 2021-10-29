package way

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/core"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/google/uuid"
)

type Way interface {
	core.Identitable
	ValidateProfile(ctx context.Context, profile profile.Profile) (map[string]string, error)
	SetParams(params model.Params)
	//GetRequiredFields() []string
	Stop() error
}

type Sender interface {
	Way
	Send(ctx context.Context, profile profile.Profile, message message.Message) ([]string, error)
}

type Listener interface {
	Way
	Listen(context.Context, func(ctx context.Context, message message.Message) error) error
}

type Wayer interface {
	GetAvailableListeners(ctx context.Context) ([]Listener, error)
	GetSenderByName(ctx context.Context, name string) (Sender, error)
	GetSenderById(ctx context.Context, id uuid.UUID) (Sender, error)

	//The RunListeners starts listen all available listeners.
	RunListeners(ctx context.Context, handler func(ctx context.Context, message message.Message) error) error

	//The RunListener method starts listen listener. If error happens, restarts listening.
	RunListener(ctx context.Context, listener Listener, handler func(ctx context.Context, message message.Message) error) error

	StopListener(ctx context.Context, listener Listener) error

	// Stop all listeners and senders
	Stop() error
}
