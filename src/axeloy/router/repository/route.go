package repository

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/router/model"
)

type RouteRepository interface {
	CreateRoute(ctx context.Context, route ...*model.Route) error
	CreateTrack(ctx context.Context, track ...*model.Track) error
	GetBySource(ctx context.Context, source message.Payload) ([]*model.Route, error)
	GetTracks(ctx context.Context, m message.Message) ([]*model.Track, error)
}
