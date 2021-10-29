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
	"time"
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
		name                               string
		fields                             fields
		messageServiceSaveReturn           error
		messageServiceSaveStateReturn      error
		routerServiceGetDestinationsReturn func() ([]router.Destination, error)
		trackServiceDefineTracksReturn     func() ([]router.Track, error)
		sentTracks                         int
		err                                error
	}{
		{
			name:                     `Error to save message`,
			messageServiceSaveReturn: errors.New(`any error`),
			err:                      ErrSaveMessage,
		},
		{
			name: `Error to got destinations, saving state error`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return nil, errors.New(`any error`)
			},
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
		},
		{
			name: `Error to got destinations, saving state ok`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return nil, errors.New(`any error`)
			},
			messageServiceSaveStateReturn: nil,
			err:                           ErrFetchDestinations,
		},
		{
			name: `Error when getting destinations, saving state ok`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{}, nil
			},
			messageServiceSaveReturn: nil,
			err:                      ErrNoDestinations,
		},
		{
			name: `Error when getting destinations, saving state error`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{}, nil
			},
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
		},
		{
			name:                     `No destinations`,
			messageServiceSaveReturn: nil,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, errors.New(`any error`)
			},
			err: ErrFetchDestinations,
		},
		{
			name: `Error when getting tracks, saving state ok`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, nil
			},
			trackServiceDefineTracksReturn: func() ([]router.Track, error) {
				return nil, errors.New(`any error`)
			},
			messageServiceSaveReturn: nil,
			err:                      ErrDefineTracks,
		},
		{
			name: `Error when getting tracks, saving state error`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, nil
			},
			trackServiceDefineTracksReturn: func() ([]router.Track, error) {
				return nil, errors.New(`any error`)
			},
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
		},
		{
			name: `Range tracks`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, nil
			},
			trackServiceDefineTracksReturn: func() ([]router.Track, error) {
				return []router.Track{
					&mocks.Track{},
					&mocks.Track{},
				}, nil
			},
			sentTracks: 2,
		},
		{
			name: `Error save Processed message's state`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, nil
			},
			trackServiceDefineTracksReturn: func() ([]router.Track, error) {
				return []router.Track{
					&mocks.Track{},
				}, nil
			},
			sentTracks:                    1,
			messageServiceSaveStateReturn: errors.New(`any error`),
			err:                           ErrInternalError,
		},
		{
			name: `No errors`,
			routerServiceGetDestinationsReturn: func() ([]router.Destination, error) {
				return []router.Destination{
					&mocks.Destination{},
				}, nil
			},
			trackServiceDefineTracksReturn: func() ([]router.Track, error) {
				return []router.Track{
					&mocks.Track{},
				}, nil
			},
			sentTracks:                    1,
			messageServiceSaveStateReturn: nil,
			err:                           nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.Background()

			if tt.routerServiceGetDestinationsReturn == nil {
				tt.routerServiceGetDestinationsReturn = func() ([]router.Destination, error) {
					return nil, nil
				}
			}

			if tt.trackServiceDefineTracksReturn == nil {
				tt.trackServiceDefineTracksReturn = func() ([]router.Track, error) {
					return nil, nil
				}
			}

			var messageMock = new(mocks4.Message)

			var messageServiceMock = new(mocks4.Messager)
			messageServiceMock.On("Save", ctx, messageMock).Return(tt.messageServiceSaveReturn)
			messageServiceMock.On("SaveState", ctx, messageMock, mock.Anything, mock.Anything).
				Return(tt.messageServiceSaveStateReturn)

			var routerServiceMock = new(mocks.Router)
			routerServiceMock.On("GetDestinations", ctx, messageMock).Return(tt.routerServiceGetDestinationsReturn())

			var trackServiceMock = new(mocks.Tracker)
			trackServiceMock.On("DefineTracks", ctx, messageMock, mock.Anything).Return(tt.trackServiceDefineTracksReturn())

			a := &Axeloy{
				messageService: messageServiceMock,
				routerService:  routerServiceMock,
				trackService:   trackServiceMock,
				senderChan:     make(chan router.Track),
			}
			defer func() {
				close(a.senderChan)
			}()

			var sent = 0
			go func() {
				for range a.senderChan {
					sent++
				}
			}()
			err := a.Handle(ctx, messageMock)

			if !assert.ErrorIs(t, err, tt.err) {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.err)
			}
			time.Sleep(time.Microsecond * time.Duration(tt.sentTracks))
			assert.Equal(t, tt.sentTracks, sent)
		})
	}
}

func TestAxeloy_Subscribe1(t *testing.T) {

	tests := []struct {
		name                          string
		routerServiceApplyRouteReturn error
		wantErr                       error
	}{
		{
			name:                          `Error`,
			routerServiceApplyRouteReturn: errors.New(`any error`),
			wantErr:                       ErrSubscription,
		},
		{
			name:                          `No error`,
			routerServiceApplyRouteReturn: nil,
			wantErr:                       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ctx = context.TODO()

			var routerServiceMock = new(mocks.Router)
			routerServiceMock.On("ApplyRoute", ctx, mock.Anything, mock.Anything, mock.Anything).Return(tt.routerServiceApplyRouteReturn)

			a := &Axeloy{
				routerService: routerServiceMock,
			}
			err := a.Subscribe(ctx, &mocks2.Profile{}, &mocks2.Profile{}, &mocks3.Sender{})
			if !assert.ErrorIs(t, err, tt.wantErr) {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
