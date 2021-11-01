package way

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
)

type WayRepository struct {
	db *reform.DB
}

func NewWayRepository(db *reform.DB) *WayRepository {
	return &WayRepository{
		db: db,
	}
}

func (w *WayRepository) GetByType(ctx context.Context, t model.WayType) ([]model.Way, error) {
	panic("implement me")
}

func (w *WayRepository) GetByName(ctx context.Context, name string) (model.Way, error) {
	panic("implement me")
}

func (w *WayRepository) GetById(ctx context.Context, id uuid.UUID) (model.Way, error) {
	panic("implement me")
}

func (w *WayRepository) AddListener(ctx context.Context, name string, driverName string, params map[string]interface{}) error {
	panic("implement me")
}

func (w *WayRepository) AddSender(ctx context.Context, name string, driverName string, params map[string]interface{}) error {
	panic("implement me")
}
