package axeloy

import (
	"context"
	"errors"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/way"
	mocks4 "github.com/dmalykh/axeloy/mocks/axeloy/message"
	mocks2 "github.com/dmalykh/axeloy/mocks/axeloy/profile"
	mocks "github.com/dmalykh/axeloy/mocks/axeloy/router"
	mocks3 "github.com/dmalykh/axeloy/mocks/axeloy/way"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAxeloy_Subscribe(t *testing.T) {
	var ctx = context.Background()
	var routeServiceMock = new(mocks.Router)
	routeServiceMock.On("ApplyRoute", ctx)
	var ax = &Axeloy{
		routerService: routeServiceMock,
	}

	ax.Subscribe(context.Background(), &mocks2.Profile{}, &mocks2.Profile{}, &mocks3.Sender{})
}

func TestAxeloy_Handle(t *testing.T) {
	type fields struct {
		routerService  router.Router
		messageService message.Messager
		waysService    way.Wayer
		log            Logger
	}

	tests := []struct {
		name                                 string
		fields                               fields
		messageServiceSaveReturn             error
		messageServiceSaveStateReturn        error
		routerServiceGetDestinationsReturn   func() ([]router.Destination, error)
		routerServiceApplyDestinationsReturn error
		err                                  error
	}{
		{
			name:                     `Error to save message`,
			messageServiceSaveReturn: errors.New(`any error`),
			err:                      ErrSaveMessage,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return nil, nil
			},
		},
		{
			name:                          `Error to save messages state when destinations not found`,
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return nil, errors.New(`any error`)
			},
		},
		{
			name:                          `Error to save messages state when destinations couldn't be saved`,
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return nil, nil
			},
			routerServiceApplyDestinationsReturn: errors.New(`any error`),
		},
		{
			name:                     `No destinations`,
			messageServiceSaveReturn: nil,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{}, nil
			},
			err: ErrNoDestinations,
		},
		{
			name:                     `Get destinations error`,
			messageServiceSaveReturn: nil,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, errors.New(`any error`)
			},
			err: ErrFetchDestinations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()

			var messageMock = new(mocks4.Message)

			var messageServiceMock = new(mocks4.Messager)
			messageServiceMock.On("Save", ctx, messageMock).Return(tt.messageServiceSaveReturn)
			messageServiceMock.On("SaveState", ctx, messageMock, mock.Anything, mock.Anything).
				Return(tt.messageServiceSaveStateReturn)

			var routerServiceMock = new(mocks.Router)
			routerServiceMock.On("GetDestinations", ctx, messageMock).Return(tt.routerServiceGetDestinationsReturn())
			routerServiceMock.On("ApplyDestinations", ctx, messageMock, mock.Anything).Return(tt.routerServiceApplyDestinationsReturn)

			a := &Axeloy{
				messageService: messageServiceMock,
				routerService:  routerServiceMock,
			}
			err := a.Handle(ctx, messageMock)

			if !assert.ErrorIs(t, err, tt.err) {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
