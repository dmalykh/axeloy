package repository

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/google/uuid"
)

type TrackRepository interface {
	CreateTrack(ctx context.Context, track ...*model.Track) error
	GetTracksByMessageId(ctx context.Context, messageId uuid.UUID) ([]*model.Track, error)
	SetStatus(ctx context.Context, id uuid.UUID, status model.TrackStatus, info ...string) error
	StartAttempt(ctx context.Context, trackId uuid.UUID, status model.AttemptStatus) (*model.Attempt, error)
	FinishAttempt(ctx context.Context, attempt *model.Attempt, status model.AttemptStatus, info ...string) error
}
