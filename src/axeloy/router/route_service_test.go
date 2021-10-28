package router

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/core"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/way"
	messageMock "github.com/dmalykh/axeloy/mocks/axeloy/message"
	mocks "github.com/dmalykh/axeloy/mocks/axeloy/profile"
	routerMock "github.com/dmalykh/axeloy/mocks/axeloy/router/repository"
	wayMock "github.com/dmalykh/axeloy/mocks/axeloy/way"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRouteService_GetDestinationsByRoute(t *testing.T) {
	type getBySource struct {
		routes []*model.Route
		err    error
	}
	type testCase struct {
		name          string
		messageSource message.Payload
		GetBySource   getBySource
		err           error
	}

	type routeService struct {
		RouteService
		mock.Mock
	}

	var ctx = context.Background()

	//Set test cases
	var testCases = []testCase{
		{
			name:          `Route exists`,
			messageSource: &messageMock.Payload{},
			GetBySource: getBySource{[]*model.Route{
				{
					Destination: &mocks.Profile{},
					Senders: []way.Sender{
						&wayMock.Sender{},
						&wayMock.Sender{},
						&wayMock.Sender{},
					},
				},
			}, nil},
			err: nil,
		},
		{
			name:          `Repository error when route fetched`,
			messageSource: &messageMock.Payload{},
			GetBySource:   getBySource{nil, core.ErrRepositoryFetchError},
			err:           ErrGettingDestinationError,
		},
		{
			name:          `Repository unknown error`,
			messageSource: &messageMock.Payload{},
			GetBySource:   getBySource{nil, errors.New(`another error`)},
			err:           nil,
		},
	}

	//Run testcases
	for _, tc := range testCases {
		//Prepare for testing
		var s = new(routeService)
		var routeRepositoryMock = new(routerMock.RouteRepository)
		routeRepositoryMock.On("GetBySource", ctx, tc.messageSource).Return(tc.GetBySource.routes, tc.GetBySource.err)
		var wantDestinations = make([]Destination, 0)
		for _, route := range tc.GetBySource.routes {
			var dest = s.MakeDestination(route.GetDestination(), route.GetSenders()...)
			//
			// ATTENTION! Destination in route MUST be unique
			//
			s.On("MakeDestination", route.GetDestination(), mock.AnythingOfType(`way.Sender`)).
				Return(dest)
			wantDestinations = append(wantDestinations, dest)
		}
		s.routeRepository = routeRepositoryMock

		t.Run(fmt.Sprintf(`%s`, tc.name), func(t *testing.T) {
			destinations, err := s.GetDestinationsByRoute(ctx, tc.messageSource)
			assert.ErrorIs(t, err, tc.err)
			if err == nil {
				assert.ElementsMatch(t, destinations, wantDestinations)
			}
		})
	}
}

func TestRouteService_GetDestinationsForDirectMessage(t *testing.T) {
	type testDestination struct {
		ways []string
	}
	type testCase struct {
		name                string
		messageDestinations []testDestination
		err                 error
	}

	var ctx = context.Background()

	//Register ways
	var wayServiceMock = new(wayMock.Wayer)
	wayServiceMock.On("GetSenderByName", ctx, `dummy`).Return(&wayMock.Sender{}, nil)
	wayServiceMock.On("GetSenderByName", ctx, `huyumi`).Return(&wayMock.Sender{}, nil)
	wayServiceMock.On("GetSenderByName", ctx, `bobor`).Return(nil, errors.New(`no bobor`))

	//Set test cases
	var testCases = []testCase{
		{
			name: `All clear, only known ways`,
			messageDestinations: []testDestination{
				{
					ways: []string{`dummy`, `huyumi`},
				},
			},
			err: nil,
		},
		{
			name: `Message with unknown way`,
			messageDestinations: []testDestination{
				{
					ways: []string{`dummy`, `bobor`},
				},
			},
			err: ErrUnknownSender,
		},
	}

	//Run testcases
	for _, tc := range testCases {
		//Prepare for testing, create mocks
		var s = RouteService{
			wayService: wayServiceMock,
		}
		var m = new(messageMock.Message)
		//Set router destinations for message
		var lenDestinations = 0
		if len(tc.messageDestinations) > 0 {
			var destinations = func(td []testDestination) []message.Payload {
				var destinations = make([]message.Payload, len(td))
				for i, d := range td {
					var dest messageMock.Destination
					dest.On("GetWays").Return(d.ways)
					dest.On("GetProfile").Return(&mocks.Profile{})
					destinations[i] = &dest
					lenDestinations = lenDestinations + len(d.ways)
				}
				return destinations
			}(tc.messageDestinations)
			m.On("GetDestinations").Return(destinations)
		}

		t.Run(fmt.Sprintf(`%s`, tc.name), func(t *testing.T) {
			destinations, err := s.GetDestinationsForDirectMessage(ctx, m.GetDestinations())
			assert.ErrorIs(t, err, tc.err)
			if err == nil {
				assert.Len(t, destinations, lenDestinations)
			}
		})
	}
}
