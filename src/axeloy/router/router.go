package router

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
)

type Router interface {
	// The GetDestinations method returns destinations for message
	GetDestinations(ctx context.Context, m message.Message) ([]Destination, error)
	//
	//// The CreateDestination method create destination
	//CreateDestination(ctx context.Context, p profile.Profile, ways ...way.Sender) (Payload, error)

	// The ApplyRoute method adds new route for messages received from source profile for sending them for destination profile by ways.
	// If a route exists, the method does nothing or updates ways list if it is necessary.
	ApplyRoute(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error

	// The ApplyDestinations method saves relation between message and destination
	ApplyDestinations(ctx context.Context, m message.Message, destinations []Destination) error
}
