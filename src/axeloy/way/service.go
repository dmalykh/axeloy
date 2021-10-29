package way

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/dmalykh/axeloy/axeloy/way/repository"
	"github.com/google/uuid"
	"plugin"
)

type WayService struct {
	wayRepository repository.WayRepository
	drivers       drivers
	listeners     map[uuid.UUID]Listener
}

var (
	ErrNoWayName       = errors.New(`name of way doesn't exists`)
	ErrNoWayDriver     = errors.New(`driver doesn't exists`)
	ErrUnknownListener = errors.New(`can't receive listener`)
	ErrStopListener    = errors.New(`can't stop listener`)
)

func (w *WayService) GetSenderByName(ctx context.Context, name string) (Sender, error) {
	way, err := w.wayRepository.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, name, ErrNoWayName)
	}
	sender, err := w.drivers.GetSender(way.GetDriverName())
	if err != nil {
		return nil, fmt.Errorf(`%s %w`, way.GetDriverName(), ErrNoWayDriver)
	}
	return sender, nil
}

type Config struct {
	Ways map[string]struct {
		DriverName string
		DriverPath string
		Params     map[string]interface{}
	}
}

func NewService() *WayService {
	return &WayService{}
}

var ErrrNotImplementAnyOfWaysInterfaces = errors.New(`doesn't implements any of ways interfaces`)

// The function loads ways from configuration.
func (w *WayService) load(ctx context.Context, config Config) error {
	for name, wayConfig := range config.Ways {
		plug, err := plugin.Open(wayConfig.DriverPath)
		if err != nil {
			return err //@TODO
		}
		way, err := plug.Lookup("Way")
		if err != nil {
			return err //@TODO
		}
		//If driver is listener, register it
		listener, canListen := way.(Listener)
		if canListen {
			if err := w.drivers.RegistryListener(wayConfig.DriverName, listener); err != nil {
				return err //@TODO
			}
			//if err := w.wayRepository.AddListener(ctx, name, wayConfig.DriverName, wayConfig.Params); err != nil {
			//	return err //@TODO
			//} Тут регаются системные драйверы, а ways могут приплывать "на лету" из базы, например
		}
		//If driver is sender, register it
		sender, canSend := way.(Sender)
		if canSend {
			if err := w.drivers.RegistrySender(wayConfig.DriverName, sender); err != nil {
				return err //@TODO
			}
			//if err := w.wayRepository.AddSender(ctx, name, wayConfig.DriverName, wayConfig.Params); err != nil {
			//	return err //@TODO
			//} Тут регаются системные драйверы, а ways могут приплывать "на лету" из базы, например
		}
		//If driver isn't sender and isn't listener, what is it?!
		if !canListen && !canSend {
			return fmt.Errorf(`%s %w`, name, ErrrNotImplementAnyOfWaysInterfaces)
		}
	}
	return nil
}

// The GetAvailableListeners methos returns ways registered as Listener
func (w *WayService) GetAvailableListeners(ctx context.Context) ([]Listener, error) {
	var listeners = make([]Listener, 0)
	//Get listeners from repository
	ways, err := w.wayRepository.GetByType(ctx, model.Listener)
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
		listener.SetParams(way.GetParams()) //@TODO: check params fields
		listeners = append(listeners, listener)
	}
	return listeners, nil
}

func (w *WayService) StopListener(ctx context.Context, listener Listener) error {
	if _, exists := w.listeners[listener.GetId()]; !exists {
		return fmt.Errorf(`%w %s`, ErrUnknownListener, listener.GetId())
	}
	if err := listener.Stop(); err != nil {
		return fmt.Errorf(`%w %s %s`, ErrStopListener, listener.GetId(), err.Error())
	}
	delete(w.listeners, listener.GetId())
	listener = nil
	return nil
}

func (w *WayService) RunListener(ctx context.Context, listener Listener, handler func(ctx context.Context, message message.Message) error) error {
	w.listeners[listener.GetId()] = listener

	go func(listener Listener) {
		if err := listener.Listen(ctx, handler); err != nil {
			if err := w.StopListener(ctx, listener); err != nil {
				//@TODO
			}
			w.RunListener(ctx, listener, handler)
		}
	}(w.listeners[listener.GetId()])

	return nil
}

func (w *WayService) RunListeners(ctx context.Context, handler func(ctx context.Context, message message.Message) error) error {
	w.listeners = make(map[uuid.UUID]Listener)

	listeners, err := w.GetAvailableListeners(ctx)
	if err != nil {
		return err
	}

	for _, listener := range listeners {
		if err := w.RunListener(ctx, listener, handler); err != nil {
			return w.Stop()
		}
	}
	return nil
}

func (w *WayService) Stop() error {
	panic(`not implemented`)
}
