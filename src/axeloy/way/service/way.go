package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/message/location"
	messagemodel "github.com/dmalykh/axeloy/axeloy/message/model"
	profilemodel "github.com/dmalykh/axeloy/axeloy/profile/model"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/dmalykh/axeloy/axeloy/way/repository"
	"plugin"
)

type WayService struct {
	wayRepository repository.WayRepository

	drivers   driver.Drivers
	listeners map[string]way.Listener
}

func NewService(ctx context.Context, config *way.Config) (*WayService, error) {
	var service = &WayService{
		wayRepository: config.WayRepository,
	}
	return service, service.load(config)
}

var (
	ErrNoWayName                       = errors.New(`name of way doesn't exists`)
	ErrNoWayDriver                     = errors.New(`driver doesn't exists`)
	ErrUnknownListener                 = errors.New(`can't receive listener`)
	ErrStopListener                    = errors.New(`can't stop listener`)
	ErrOpenPlugin                      = errors.New(`can't open plugin`)
	ErrNotImplementAnyOfWaysInterfaces = errors.New(`doesn't implements any of ways interfaces`)
)

func (w *WayService) GetSenderByName(ctx context.Context, name string) (way.Sender, error) {
	way, err := w.wayRepository.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, name, ErrNoWayName)
	}
	sender, err := w.drivers.GetSender(way.GetDriverName())
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, way.GetDriverName(), ErrNoWayDriver)
	}
	return model.MakeSender(way, sender), nil
}

//func (w *WayService) GetSenderById(ctx context.Context, id uuid.UUID) (driver.Sender, error) { @TODO REMOVE
//	way, err := w.wayRepository.GetById(ctx, id)
//	if err != nil {
//		return nil, fmt.Errorf(`%s %w`, id.String(), ErrNoWayName)
//	}
//	sender, err := w.drivers.GetSender(way.GetDriverName())
//	if err != nil {
//		return nil, fmt.Errorf(`%s %w`, way.GetDriverName(), ErrNoWayDriver)
//	}
//	return model.MakeSender(way, sender), nil
//}

// The function loads ways from configuration.
func (w *WayService) load(config *way.Config) error {
	for name, driverConfig := range config.Drivers {
		plug, err := plugin.Open(driverConfig.Path)
		if err != nil {
			return fmt.Errorf(`%w "%s": %s`, ErrOpenPlugin, driverConfig.Path, err.Error())
		}
		d, err := plug.Lookup("Driver")
		if err != nil {
			return err //@TODO
		}
		//If driver is listener, register it
		listener, canListen := d.(driver.Listener)
		if canListen {
			listener.SetDriverConfig(driverConfig.Config)
			if err := w.drivers.RegistryListener(name, listener); err != nil {
				return err //@TODO
			}
		}
		//If driver is sender, register it
		sender, canSend := d.(driver.Sender)
		if canSend {
			sender.SetDriverConfig(driverConfig.Config)
			if err := w.drivers.RegistrySender(name, sender); err != nil {
				return err //@TODO
			}
		}
		//If driver isn't sender and isn't listener, what is it?!
		if !canListen && !canSend {
			return fmt.Errorf(`%s %w`, name, ErrNotImplementAnyOfWaysInterfaces)
		}
	}
	return nil
}

// The GetAvailableListeners methos returns ways registered as Listener
func (w *WayService) GetAvailableListeners(ctx context.Context) ([]way.Listener, error) {
	var listeners = make([]way.Listener, 0)
	//Get listeners from repository
	ways, err := w.wayRepository.GetByType(ctx, model.WayTypeListener)
	if err != nil {
		return listeners, err //@TODO
	}
	for _, way := range ways {
		//Get driver for way
		listener, err := w.drivers.GetListener(way.GetDriverName())
		if err != nil {
			return listeners, err //@TODO
		}
		//Set params from way and append
		listeners = append(listeners, model.MakeListener(way, listener))
	}
	return listeners, nil
}

func (w *WayService) StopListener(ctx context.Context, listener way.Listener) error {
	if _, exists := w.listeners[listener.GetName()]; !exists {
		return fmt.Errorf(`%w %s`, ErrUnknownListener, listener.GetName())
	}
	if err := listener.Stop(ctx); err != nil {
		return fmt.Errorf(`%w %s %s`, ErrStopListener, listener.GetName(), err.Error())
	}
	delete(w.listeners, listener.GetName())
	listener = nil
	return nil
}

func (w *WayService) RunListener(ctx context.Context, listener way.Listener, handler func(ctx context.Context, message message.Message) error) error {
	w.listeners[listener.GetName()] = listener

	//Create handler function
	var handleFunc = func(ctx context.Context, m driver.Message) error {
		//Create message
		var msg = &messagemodel.Message{
			Payload: m.GetPayload(),
			Source: &messagemodel.Location{
				Way: listener.GetName(),
				Profile: &profilemodel.Profile{
					Fields: m.GetPublisher(),
				},
			},
			Destination: make([]location.Location, len(m.GetDestinations())),
		}

		//Add destinations for direct messages
		for i, destination := range m.GetDestinations() {
			msg.Destination[i] = &messagemodel.Location{
				Way: destination.GetWay(),
				Profile: &profilemodel.Profile{
					Fields: destination.GetConsumer(),
				},
			}
		}

		return handler(ctx, msg)
	}

	go func(listener way.Listener) {
		if err := listener.Listen(ctx, handleFunc); err != nil {
			if err := w.StopListener(ctx, listener); err != nil {
				//@TODO
			}
			w.RunListener(ctx, listener, handler)
		}
	}(w.listeners[listener.GetName()])

	return nil
}

func (w *WayService) RunListeners(ctx context.Context, handler func(ctx context.Context, message message.Message) error) error {
	w.listeners = make(map[string]way.Listener)

	//Get listeners
	listeners, err := w.GetAvailableListeners(ctx)
	if err != nil {
		return err
	}

	//Run listeners
	for _, listener := range listeners {
		if err := w.RunListener(ctx, listener, handler); err != nil {
			return w.Stop()
		}
	}
	return nil
}

func (w *WayService) Stop() error {
	panic(`not implemented`) //@TODO
}
