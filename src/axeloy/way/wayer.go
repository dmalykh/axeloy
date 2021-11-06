package way

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

//type WayId string

type Way interface {
	// GetId returns id. Same as GetName() method
	//GetId() WayId

	// GetName returns unique name of way
	GetName() string

	// GetTitle returns title of way
	GetTitle() string
}

type Sender interface {
	Way
	driver.Sender
}

type Listener interface {
	Way
	driver.Listener
}

type Wayer interface {
	GetAvailableListeners(ctx context.Context) ([]Listener, error)
	GetSenderByName(ctx context.Context, name string) (Sender, error)
	//GetSenderById(ctx context.Context, id uuid.UUID) (Sender, error)

	//The RunListeners starts listen all available listeners.
	RunListeners(ctx context.Context, handler func(ctx context.Context, message message.Message) error) error

	//The RunListener method starts listen listener. If error happens, restarts listening.
	RunListener(ctx context.Context, listener Listener, handler func(ctx context.Context, message message.Message) error) error

	StopListener(ctx context.Context, listener Listener) error

	// Stop all listeners and senders
	Stop() error
}
