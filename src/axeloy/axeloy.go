package axeloy

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	messagemodel "github.com/dmalykh/axeloy/axeloy/message/model"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var (
	ErrSaveMessage       = errors.New(`saving message error`)
	ErrNoDestinations    = errors.New(`not found destinations for message`)
	ErrFetchDestinations = errors.New(`error finding destinations for message`)
	ErrInternalError     = errors.New(`internal error`)
	ErrSubscription      = errors.New(`can't subscribe`)
	ErrDefineTracks      = errors.New(`can't define tracks from destinations`)
	ErrRunListeners      = errors.New(`can't run listener`)
	ErrReceiveUnsent     = errors.New(`can't receive unsent`)
)

type Axeloy struct {
	routerService  router.Router
	trackService   router.Tracker
	messageService message.Messager
	waysService    way.Wayer

	senderChan chan router.Track
}

type Config struct {
	Router   router.Router
	Tracker  router.Tracker
	Messager message.Messager
	Wayer    way.Wayer
}

func New(config *Config) *Axeloy {
	return &Axeloy{
		routerService:  config.Router,
		trackService:   config.Tracker,
		messageService: config.Messager,
		waysService:    config.Wayer,
	}
}

// Run Axeloy.
func (a *Axeloy) Run(ctx context.Context) error {
	a.senderChan = make(chan router.Track)
	a.runSender(ctx, a.senderChan)
	if err := a.sendUnsent(ctx); err != nil {
		return fmt.Errorf(`%w: %s`, ErrReceiveUnsent, err.Error())
	}
	// Start listeners
	if err := a.waysService.RunListeners(ctx, a.Handle); err != nil {
		return fmt.Errorf(`%w: %s`, ErrRunListeners, err.Error())
	}
	return nil
}

func (a *Axeloy) sendUnsent(ctx context.Context) error {
	tracks, err := a.trackService.GetUnsentTracks(ctx)
	if err != nil {
		return err
	}
	go a.send(tracks...)
	return nil
}

func (a *Axeloy) runSender(ctx context.Context, sender chan router.Track) {
	go func(ctx context.Context, sender chan router.Track) {
		for track := range sender {
			if err := a.trackService.Send(ctx, track); err != nil {
				//@TODO
			}
		}
	}(ctx, sender)
}

//The GetWaysListeners returns way for listen
func (a *Axeloy) GetWaysListeners(ctx context.Context) ([]way.Listener, error) {
	listeners, err := a.waysService.GetAvailableListeners(ctx)
	if err != nil {
		return nil, err
	}
	return listeners, nil
}

// Subscribe profile. If a message will be received from a source, it should be sent to a destination by a ways.
func (a *Axeloy) Subscribe(ctx context.Context, source profile.Profile, destination profile.Profile, senders ...way.Sender) error {
	if err := a.routerService.ApplyRoute(ctx, source, destination, senders...); err != nil {
		return fmt.Errorf(`%w:%s`, ErrSubscription, err.Error())
	}
	return nil
}

// The Handle function receives message and saves it
func (a *Axeloy) Handle(ctx context.Context, m message.Message) error {
	//Save message
	if err := a.messageService.Add(ctx, m); err != nil {
		return fmt.Errorf(`%w:%s`, ErrSaveMessage, err.Error())
	}

	//Get destinations for message
	destinations, err := a.routerService.GetDestinations(ctx, m)
	if err != nil {
		if err := a.messageService.UpdateStatus(ctx, m, messagemodel.StatusError, err.Error()); err != nil {
			return fmt.Errorf(`%w: %s`, ErrInternalError, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrFetchDestinations, err.Error())
	}

	//Return error when destinations absent
	if len(destinations) == 0 {
		if err := a.messageService.UpdateStatus(ctx, m, messagemodel.StatusNoDestinations); err != nil {
			return fmt.Errorf(`%w:%s`, ErrInternalError, err.Error())
		}
		return ErrNoDestinations
	}

	//Define tracks for message and send them to sender channel
	err = a.Send(ctx, m, destinations...)
	if err != nil {
		if err := a.messageService.UpdateStatus(ctx, m, messagemodel.StatusError, err.Error()); err != nil {
			return fmt.Errorf(`%w: %s`, ErrInternalError, err.Error())
		}
		return fmt.Errorf(`%w %s`, ErrDefineTracks, err.Error())
	}

	//Mark message as processed
	if err := a.messageService.UpdateStatus(ctx, m, messagemodel.StatusProcessed); err != nil {
		return fmt.Errorf(`%w:%s`, ErrInternalError, err.Error())
	}
	return nil
}

//Send message by destinations. Define tracks for messages destinations and sent them.
func (a *Axeloy) Send(ctx context.Context, m message.Message, destinations ...router.Destination) error {
	for _, destination := range destinations {
		//Define tracks from destinations for the messages history
		tracks, err := a.trackService.DefineTracks(ctx, m, destination)
		if err != nil {
			return err
		}
		a.send(tracks...)
	}
	return nil
}

// The send method sends message related with track by track
func (a *Axeloy) send(tracks ...router.Track) {
	for _, track := range tracks {
		a.senderChan <- track
	}
}
