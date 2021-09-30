package axeloy

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var ErrSaveMessage = errors.New(`saving message error`)
var ErrNoDestinations = errors.New(`not found destinations for message`)
var ErrFetchDestinations = errors.New(`error finding destinations for message`)
var ErrInternalError = errors.New(`internal error`)
var ErrSubscription = errors.New(`can't subscribe`)
var ErrDefineTracks = errors.New(`can't define tracks from destinations`)
var ErrSend = errors.New(`can't send`)

type Axeloy struct {
	routerService  router.Router
	messageService message.Messager
	waysService    way.Wayer
	log            Logger
	senderChan     chan *model.Track
	ctx            context.Context
}

type Config struct {
	router router.Router
}

// Run Axeloy.
func (a *Axeloy) Run(ctx context.Context, config *Config) error {
	a.senderChan = make(chan *model.Track)
	// Start listen ways.
	for _, w := range a.GetWaysListeners(ctx) {
		go func(way way.Listener) {
			err := way.Listen(ctx, a.Handle)
			if err != nil {
				//@TODO restart channel, log  restart. Стоппинг каналов
			}
		}(w)
	}
	return nil
}

func (a *Axeloy) GetContext() context.Context {
	return a.ctx
}

func (a *Axeloy) runSender() {
	for track := range a.senderChan {
		//track.GetSender().ValidateProfile(track.GetProfile())
		if err := a.routerService.Send(a.GetContext(), track); err != nil {

		}
	}
}

// Subscribe profile. If a message will be received from a source, it should be sent to a destination by a ways.
func (a *Axeloy) Subscribe(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error {
	if err := a.routerService.ApplyRoute(ctx, source, destination, senders...); err != nil {
		return fmt.Errorf(`%w:%s`, ErrSubscription, err.Error())
	}
	return nil
}

func (a *Axeloy) GetWaysListeners(ctx context.Context) []way.Listener {
	listeners, err := a.waysService.GetAvailableListeners(ctx)
	if err != nil {
		//@TODO
	}
	return listeners
}

// The Handle function receives message and saves it
func (a *Axeloy) Handle(ctx context.Context, m message.Message) error {
	//Save message
	if err := a.messageService.Save(ctx, m); err != nil {
		return fmt.Errorf(`%w:%s`, ErrSaveMessage, err.Error())
	}

	//Get destinations for message
	destinations, err := a.routerService.GetDestinations(ctx, m)
	if err != nil {
		if err := a.messageService.SaveState(ctx, m, message.Error, err.Error()); err != nil {
			return fmt.Errorf(`%w: %s`, ErrInternalError, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrFetchDestinations, err.Error())
	}

	//Return error when destinations absent
	if len(destinations) == 0 {
		if err := a.messageService.SaveState(ctx, m, message.NoDestinations); err != nil {
			return fmt.Errorf(`%w:%s`, ErrInternalError, err.Error())
		}
		return ErrNoDestinations
	}

	//Define tracks for message and send them to sender channel
	tracks, err := a.defineTracks(ctx, m, destinations...)
	if err != nil {
		if err := a.messageService.SaveState(ctx, m, message.Error, err.Error()); err != nil {
			return fmt.Errorf(`%w: %s`, ErrInternalError, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrDefineTracks, err.Error())
	}

	//Send messages by  tracks
	for _, track := range tracks {
		a.senderChan <- track
	}

	//Mark message as processed
	if err := a.messageService.SaveState(ctx, m, message.Processed); err != nil {
		return fmt.Errorf(`%w:%s`, ErrInternalError, err.Error())
	}
	return nil
}

func (a *Axeloy) defineTracks(ctx context.Context, m message.Message, destinations ...router.Destination) ([]*model.Track, error) {
	var tracks = make([]*model.Track, 0)
	for _, destination := range destinations {
		//Define tracks from destinations for the messages history
		t, err := a.routerService.DefineTracks(ctx, m, destination)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, t...)
	}
	return tracks, nil
}
