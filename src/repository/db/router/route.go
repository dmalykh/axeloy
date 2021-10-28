package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/profile"
	profilemodel "github.com/dmalykh/axeloy/axeloy/profile/model"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	dbmodel "github.com/dmalykh/axeloy/repository/db/router/model"
	"github.com/google/uuid"
	"gopkg.in/reform.v1"
	"log"
	"strings"
)

var ErrCreateRoute = errors.New(`can't insert route`)
var ErrCreateProfile = errors.New(`profile couldn't be saved`)
var ErrCreateWays = errors.New(`ways couldn't be saved`)

type RouteRepository struct {
	db *reform.DB
}

func (r *RouteRepository) GetBySourceProfile(ctx context.Context, p profile.Profile) ([]*model.Route, error) {
	var placeholders = make([]interface{}, 0)

	//Create conditions for placeholders
	for key, values := range p.GetFields() {
		for _, value := range values {
			placeholders = append(placeholders, key, value)
		}
	}

	// Database agnostic hack. PostgreSQL placeholder parameter uses indexes, ex "$1", other uses "?".
	// So right sequence of placeholders needed.
	var conditions = make([]string, len(placeholders)/2)
	var subconditions = make([]string, cap(conditions))
	var plen = len(placeholders)
	for i := 1; i < plen; i = i + 2 {
		conditions[i/2] = fmt.Sprintf(`(%s = %s AND %s = %s)`,
			r.db.QuoteIdentifier(`key`),
			r.db.Placeholder(i-1),
			r.db.QuoteIdentifier(`value`),
			r.db.Placeholder(i),
		)
		subconditions[i/2] = fmt.Sprintf(`(%s = %s AND %s = %s)`,
			r.db.QuoteIdentifier(`key`),
			r.db.Placeholder(plen+i-1),
			r.db.QuoteIdentifier(`value`),
			r.db.Placeholder(plen+i),
		)
	}

	//Generate subquery conditions
	var subtail = fmt.Sprintf(`(%s) AND %s NOT IN (SELECT %s FROM %s WHERE NOT(%s))`,
		strings.Join(conditions, ` OR `),
		r.db.QuoteIdentifier(`route_id`),
		r.db.QuoteIdentifier(`route_id`),
		dbmodel.ProfileTable.Name(),
		strings.Join(subconditions, ` OR `),
	)

	//Merge placeholders and generate main query
	placeholders = append(placeholders, placeholders...)
	placeholders = append(placeholders, dbmodel.Destination)
	var tail = fmt.Sprintf(`%s IN(SELECT %s FROM %s WHERE (%s)) AND %s=%s`,
		r.db.QuoteIdentifier(`route_id`),
		r.db.QuoteIdentifier(`route_id`),
		dbmodel.ProfileTable.Name(),
		subtail,
		r.db.QuoteIdentifier(`type`),
		r.db.Placeholder(len(placeholders)),
	)

	//Final query
	var query = fmt.Sprintf(`SELECT %s FROM %s INNER JOIN %s ON (%s.%s = %s.%s) WHERE %s `,
		strings.Join([]string{
			dbmodel.WayTable.Name() + "." + r.db.QuoteIdentifier(`route_id`),
			dbmodel.WayTable.Name() + "." + r.db.QuoteIdentifier(`way_id`),
			dbmodel.ProfileTable.Name() + "." + r.db.QuoteIdentifier(`key`),
			dbmodel.ProfileTable.Name() + "." + r.db.QuoteIdentifier(`value`),
		}, ", "),
		dbmodel.ProfileTable.Name(),
		dbmodel.WayTable.Name(),
		dbmodel.ProfileTable.Name(),
		r.db.QuoteIdentifier(`route_id`),
		dbmodel.WayTable.Name(),
		r.db.QuoteIdentifier(`route_id`),
		tail,
	)

	//Query
	rows, err := r.db.WithContext(ctx).Query(query, placeholders)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrDbQuery, err.Error())
	}
	defer rows.Close()

	type tempData struct {
		RouteId uuid.UUID
		WayId   uuid.UUID
		Key     string
		Value   string
	}
	//Fetch and collect profile fields
	var ways = make(map[uuid.UUID][]uuid.UUID)
	var profileFields = make(map[uuid.UUID]map[string][]string)
	for rows.Next() {
		var temp tempData
		if err = rows.Scan(&temp.RouteId, &temp.WayId, &temp.Key, &temp.Value); err != nil {
			log.Print(err) //@TODO
		}
		//Create slices if route doesn't exists
		if _, exists := profileFields[temp.RouteId]; !exists {
			profileFields[temp.RouteId] = make(map[string][]string)
			ways[temp.RouteId] = make([]uuid.UUID, 0)
		}
		//Collect ways
		ways[temp.RouteId] = append(ways[temp.RouteId], temp.WayId)

		//Collect fields for profile
		if _, exists := profileFields[temp.RouteId][temp.Key]; !exists {
			profileFields[temp.RouteId][temp.Key] = make([]string, 0)
		}
		profileFields[temp.RouteId][temp.Key] = append(profileFields[temp.RouteId][temp.Key], temp.Value)
	}

	//Make routes
	var routes = make([]*model.Route, 0)
	for id, fields := range profileFields {
		routes = append(routes, &model.Route{
			Id:     id,
			Source: p,
			Destination: &profilemodel.Profile{
				Fields: fields,
			},
			WaysIds: ways[id],
		})
	}
	return routes, nil
}

func (r *RouteRepository) CreateRoute(ctx context.Context, route *model.Route) error {

	var routeId = uuid.New()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err //@TODO
	}

	//CreateRoute route @TODO: а нужна ли вообще эта таблица из одной колонки?
	if err := r.createRoute(tx, routeId); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%w %s, rollback error %s`, ErrCreateRoute, routeId.String(), err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrCreateRoute, err.Error())
	}

	//Save source profiles for route
	if err := r.createProfile(tx, routeId, dbmodel.Source, route.GetSource()); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%s %w, rollback error %s`, dbmodel.Source, ErrDbError, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, dbmodel.Source, ErrCreateProfile, err.Error())
	}

	//Save destination profiles for route
	if err := r.createProfile(tx, routeId, dbmodel.Destination, route.GetDestination()); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf(`%s %w, rollback error %s`, dbmodel.Destination, ErrDbError, err.Error())
		}
		return fmt.Errorf(`%s %w %s`, dbmodel.Destination, ErrCreateProfile, err.Error())
	}

	//Save ways
	if err := r.createWays(tx, routeId, route.GetWaysIds()); err != nil {
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

func (r *RouteRepository) createWays(tx *reform.TX, routeId uuid.UUID, senders []uuid.UUID) error {
	var ways = make([]reform.Struct, len(senders))
	for i, wid := range senders {
		ways[i] = &dbmodel.Way{
			RouteId: routeId,
			WayId:   wid,
		}
	}
	return tx.InsertMulti(ways...)
}
