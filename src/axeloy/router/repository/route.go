package repository

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router/model"
)

type RouteRepository interface {
	CreateRoute(ctx context.Context, route ...*model.Route) error
	GetBySourceProfile(ctx context.Context, profile profile.Profile) ([]*model.Route, error)
}
