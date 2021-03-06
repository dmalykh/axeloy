package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/core"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/message/location"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/router/repository"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var (
	ErrUnknownSender           = errors.New(`unknown sender`)
	ErrGettingDestinationError = errors.New(`getting destination with error`)
	ErrNotValidProfile         = errors.New(`profile is not valid`)
	ErrCreateTrack             = errors.New(`can't create track`)
	ErrGetTrack                = errors.New(`couldn't got track`)
)

func NewRouter(routeRepository repository.RouteRepository, wayer way.Wayer) Router {
	return &RouteService{
		routeRepository: routeRepository,
		wayService:      wayer,
	}
}

type RouteService struct {
	routeRepository repository.RouteRepository
	wayService      way.Wayer
}

func (r *RouteService) MakeDestination(p profile.Profile, ways ...way.Sender) Destination {
	return &model.RouteDestination{
		Ways:    ways,
		Profile: p,
	}
}

func (r *RouteService) GetDestinations(ctx context.Context, m message.Message) ([]Destination, error) {
	//If message has own destinations, use them
	if payloads := m.GetDestinations(); len(payloads) > 0 {
		return r.GetDestinationsForDirectMessage(ctx, payloads)
	}
	//Find destination for message by route
	return r.GetDestinationsByRoute(ctx, m.GetSource())
}

//
func (r *RouteService) GetDestinationsByRoute(ctx context.Context, loc location.Location) ([]Destination, error) {
	routes, err := r.routeRepository.GetBySourceProfile(ctx, loc.GetProfile())
	if err != nil {
		if errors.Is(err, core.ErrRepositoryFetchError) {
			return nil, fmt.Errorf(`%w: %s`, ErrGettingDestinationError, err.Error())
		}
		return nil, nil
	}
	//Get destinations from routes
	var destinations = func(routes []*model.Route) []Destination {
		var destinations = make([]Destination, 0)
		for _, route := range routes {
			// Get only available senders @TODO: Is it possible to get thousands senders?
			var senders = func(ways []string) []way.Sender {
				var s = make([]way.Sender, 0)
				for _, name := range ways {
					sender, err := r.wayService.GetSenderByName(ctx, name)
					if err != nil {
						continue //@TODO log it
					}
					if err := sender.ValidateProfile(ctx, route.GetDestination()); err != nil {
						continue //@TODO log error
					}
					s = append(s, sender)
				}
				return s
			}(route.GetWays())
			//Add to destinations
			destinations = append(destinations, r.MakeDestination(route.GetDestination(), senders...))
		}
		return destinations
	}(routes)
	return destinations, nil
}

// The GetDestinationsForDirectMessage returns destinations for messages which has own destinations
func (r *RouteService) GetDestinationsForDirectMessage(ctx context.Context, locs []location.Location) ([]Destination, error) {
	var destinations = make([]Destination, 0)
	//Range message's destinations
	for _, loc := range locs {
		//Get way by name and convert it to route's destinations
		sender, err := r.wayService.GetSenderByName(ctx, loc.GetWay())
		if err != nil {
			return destinations, fmt.Errorf(`%w %s`, ErrUnknownSender, loc.GetWay())
		}
		destinations = append(destinations, r.MakeDestination(loc.GetProfile(), sender))
	}
	return destinations, nil
}

func (r *RouteService) ApplyRoute(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error {
	//@TODO: Check for exists route and validate ways with profile
	for _, s := range senders {
		if err := s.ValidateProfile(ctx, destination); err != nil {
			return fmt.Errorf(`%w %s`, ErrNotValidProfile, err.Error())
		}
	}
	return r.routeRepository.CreateRoute(ctx, &model.Route{
		Source:      source,
		Destination: destination,
		Ways: func(senders []way.Sender) []string {
			var ways = make([]string, len(senders))
			for i, sender := range senders {
				ways[i] = sender.GetName()
			}
			return ways
		}(senders),
	})
}
