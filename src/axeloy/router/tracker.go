package router

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
)

type Tracker interface {
	// The DefineTracks method unfolds destination to tracks and saves relation between message and track
	DefineTracks(ctx context.Context, m message.Message, destination Destination) ([]Track, error)

	// The GetUnsentTracks method returns tracks which were planned, but weren't sent
	GetUnsentTracks(ctx context.Context) ([]Track, error)

	// The Send method makes attempts to send messages from track by the track
	Send(ctx context.Context, track Track) error
}
