package repository

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/google/uuid"
)

type WayRepository interface {
	GetByType(ctx context.Context, t model.WayType) ([]model.Way, error)
	GetByName(ctx context.Context, name string) (model.Way, error)
	GetById(ctx context.Context, id uuid.UUID) (model.Way, error)
	AddListener(ctx context.Context, name string, driverName string, params map[string]interface{}) error
	AddSender(ctx context.Context, name string, driverName string, params map[string]interface{}) error
}
