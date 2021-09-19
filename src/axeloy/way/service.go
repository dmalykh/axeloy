package way

import (
	"context"
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/way/model"
	"github.com/dmalykh/axeloy/axeloy/way/repository"
	"plugin"
)

type WayService struct {
	wayRepository repository.WayRepository
	drivers       drivers
}

func (w *WayService) GetSenderByName(ctx context.Context, name string) (Sender, error) {
	way, err := w.wayRepository.GetByName(ctx, name)
	if err != nil {
		return nil, err //@TODO
	}
	sender, err := w.drivers.GetSender(way.GetDriverName())
	if err != nil {
		return nil, err //@TODO
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
