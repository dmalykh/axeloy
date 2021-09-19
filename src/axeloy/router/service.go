package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/core"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/router/repository"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var ErrUnknownSender = errors.New(`unknown sender`)
var ErrGettingDestinationError = errors.New(`getting destination with error`)

type RouteService struct {
	routeRepository repository.RouteRepository
	wayService      way.Wayer
	messafeService  message.Messager
}

func NewService() *RouteService {
	return &RouteService{}
}

func (r *RouteService) MakeDestination(p profile.Profile, ways ...way.Sender) Destination {
	return &model.RouteDestination{
		Ways:    ways,
		Profile: p,
	}
}

//
func (r *RouteService) ApplyDestinations(ctx context.Context, m message.Message, destinations []Destination) error {

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
func (r *RouteService) GetDestinationsByRoute(ctx context.Context, payload message.Payload) ([]Destination, error) {
	routes, err := r.routeRepository.GetBySource(ctx, payload)
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
			destinations = append(destinations, r.MakeDestination(route.GetDestination(), route.GetSenders()...))
		}
		return destinations
	}(routes)
	return destinations, nil
}

// The GetDestinationsForDirectMessage returns destinations for messages which has own destinations
func (r *RouteService) GetDestinationsForDirectMessage(ctx context.Context, payloads []message.Payload) ([]Destination, error) {
	var destinations = make([]Destination, 0)
	//Range message's destinations
	for _, payload := range payloads {
		//Get ways by names and convert them to route's destinations
		for _, wayName := range payload.GetWays() {
			sender, err := r.wayService.GetSenderByName(ctx, wayName)
			if err != nil {
				return destinations, fmt.Errorf(`%w %s`, ErrUnknownSender, wayName)
			}
			destinations = append(destinations, r.MakeDestination(payload.GetProfile(), sender))
		}
	}
	return destinations, nil
}

func (r *RouteService) ApplyRoute(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error {
	//@TODO: Check for exists route
	for _, s := range senders {
		if err := s.ValidateProfile(ctx, destination); err != nil {
			return err //@TODO
		}
	}
	return r.routeRepository.Create(ctx, &model.Route{
		Source:      source,
		Destination: destination,
		Senders:     senders,
	})
}
