package router

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
)

type Tracker interface {
	// The DefineTracks method unfolds destination to tracks and saves relation between message and track
	DefineTracks(ctx context.Context, m message.Message, destinations Destination) ([]Track, error)

	// The GetUnsentTracks method returns tracks which were planned, but weren't sent
	GetUnsentTracks(ctx context.Context) ([]Track, error)

	// The Send method sending messages from track by the track
	Send(ctx context.Context, track Track) error

	AddAttempt(ctx context.Context, t Track, status TrackStatus, info ...string) error
}
