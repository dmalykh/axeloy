package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/way"
	dbmodel "github.com/dmalykh/axeloy/repository/db/router/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
)

type RouteRepository struct {
	db *reform.DB
}

var ErrCreateRoute = errors.New(`can't insert route`)
var ErrCreateProfile = errors.New(`profile couldn't be saved`)
var ErrCreateWays = errors.New(`ways couldn't be saved`)

func (r *RouteRepository) Create(ctx context.Context, route *model.Route) error {

	var routeId = uuid.New() //@TODO: возможно, стоит генерировать выше

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err //@TODO
	}

	//Create route @TODO: а нужна ли вообще эта таблица из одной колонки?
	if err := r.createRoute(tx, routeId); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%w %s, rollback error %s`, ErrCreateRoute, routeId.String(), err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrCreateRoute, err.Error())
	}

	//Save source profiles for route
	if err := r.createProfile(tx, routeId, dbmodel.Source, route.GetSource()); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%s %w, rollback error %s`, dbmodel.Source, ErrCreateProfile, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, dbmodel.Source, ErrCreateProfile, err.Error())
	}

	//Save destination profiles for route
	if err := r.createProfile(tx, routeId, dbmodel.Destination, route.GetDestination()); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%s %w, rollback error %s`, dbmodel.Destination, ErrCreateProfile, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, dbmodel.Destination, ErrCreateProfile, err.Error())
	}

	//Save ways
	if err := r.createWays(tx, routeId, route.GetSenders()); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%w, rollback error %s`, ErrCreateWays, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrCreateWays, err.Error())
	}

	return tx.Commit()
}

func (r *RouteRepository) createRoute(tx *reform.TX, routeId uuid.UUID) error {
	return tx.Insert(&dbmodel.Route{
		Id: routeId,
	})
}

func (r *RouteRepository) createProfile(tx *reform.TX, routeId uuid.UUID, kind dbmodel.Type, p profile.Profile) error {
	var profiles = make([]reform.Struct, 0)
	for key, values := range p.GetFields() {
		for _, value := range values {
			profiles = append(profiles, &dbmodel.Profile{
				Id:      uuid.New(),
				RouteId: routeId,
				Type:    kind,
				Key:     key,
				Value:   value,
			})
		}
	}
	return tx.InsertMulti(profiles...)
}

func (r *RouteRepository) createWays(tx *reform.TX, routeId uuid.UUID, senders []way.Sender) error {
	var ways = make([]reform.Struct, len(senders))
	for i, w := range senders {
		ways[i] = &dbmodel.Way{
			RouteId: routeId,
			WayId:   w.GetId(),
		}
	}
	return tx.InsertMulti(ways...)
}

func (r *RouteRepository) GetBySource(ctx context.Context, source message.Payload) ([]*model.Route, error) {
	panic("implement me")
}
