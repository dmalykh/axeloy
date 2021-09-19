package axeloy

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router"
	"github.com/dmalykh/axeloy/axeloy/way"
)

var ErrSaveMessage = errors.New(`saving message error`)
var ErrNoDestinations = errors.New(`not found destinations for message`)
var ErrFetchDestinations = errors.New(`error finding destinations for message`)
var ErrSaveDestinations = errors.New(`error save destinations for message`)
var ErrInternalError = errors.New(`internal error`)
var ErrSubscription = errors.New(`can't subscribe`)

type Axeloy struct {
	routerService  router.Router
	messageService message.Messager
	waysService    way.Wayer
	log            Logger
}

type Config struct {
	router router.Router
}

// Run Axeloy.
// Starts listening ways.
func (a *Axeloy) Run(ctx context.Context, config *Config) error {
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

// The ApplyDestinations method fetch destinations for message and save relation between message and route
func (a *Axeloy) ApplyDestinations(ctx context.Context, m message.Message) ([]router.Destination, error) {
	//Get destinations for message
	destinations, err := a.routerService.GetDestinations(ctx, m)
	if err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrFetchDestinations, err.Error())
	}
	//Add destinations for the messages history
	if err := a.routerService.ApplyDestinations(ctx, m, destinations); err != nil {
		return nil, fmt.Errorf(`%w %s`, ErrSaveDestinations, err.Error())
	}
	return nil, err
}

// The Handle function receives message and saves it
func (a *Axeloy) Handle(ctx context.Context, m message.Message) error {
	//Save message
	if err := a.messageService.Save(ctx, m); err != nil {
		return fmt.Errorf(`%w:%s`, ErrSaveMessage, err.Error())
	}

	//Get destinations for message
	destinations, err := a.ApplyDestinations(ctx, m)
	if err != nil {
		if err := a.messageService.SaveState(ctx, m, message.Error, err.Error()); err != nil {
			return fmt.Errorf(`%w:%s`, ErrInternalError, err.Error())
		}
		return err
	}

	//Range destinations and send messages
	if len(destinations) == 0 {
		return ErrNoDestinations
	}
	for _, destination := range destinations {
		//@TODO ОБРАБОТКА ДУБЛЕЙ! Может быть найдено два одинаковых маршрута для одного сообщения
		go func(ctx context.Context, m message.Message, destination router.Destination) {
			for _, w := range destination.GetWays(ctx) {
				err := w.ValidateProfile(ctx, destination.GetProfile(ctx))
				if err != nil {
					if err := a.messageService.SaveState(ctx, m, message.NotValidProfile, err.Error()); err != nil { //@TODO: Add info about profile
						//@todo
					}
				}
				if err := a.Send(ctx, m, destination, w); err != nil {
					//@todo
				}
			}
		}(ctx, m, destination)
	}
	return nil
}

func (a *Axeloy) Send(ctx context.Context, m message.Message, desination router.Destination, w way.Sender) error {
	state, err := w.Send(ctx, desination.GetProfile(ctx), m)
	// Save error
	if err != nil {
		if err := a.messageService.SaveState(ctx, m, state, err.Error()); err != nil { //@TODO: Human readable error
			//@TODO err
		}
	}
	if err := a.messageService.SaveState(ctx, m, message.Sent); err != nil {
		//@TODO err
	}
	return err

}
