package repository

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/router/model"
)

type RouteRepository interface {
	Create(ctx context.Context, route *model.Route) error
	GetBySource(ctx context.Context, source message.Payload) ([]*model.Route, error)
}
