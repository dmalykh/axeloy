package service

import (
	"errors"
	"fmt"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

var (
	ErrUnknownListenerDriver           = errors.New(`unknown listener`)
	ErrUnknownSenderDriver             = errors.New(`unknown sender`)
	ErrDriverNameExists                = errors.New(`driver name already exists`)
	ErrNotImplementAnyOfWaysInterfaces = errors.New(`doesn't implements any of ways interfaces`)
)

func NewDriverService() way.Driverer {
	return &DriverService{}
}

// Load register initialized drivers as listeners and senders.
func (d *DriverService) Load(drivers map[string]driver.Driver) error {
	for name, drv := range drivers {
		//If driver is listener, register it
		listener, isListener := drv.(driver.Listener)
		if isListener {
			if err := d.RegistryListener(name, listener); err != nil {
				return err //@TODO
			}
		}
		//If driver is sender, register it
		sender, isSender := drv.(driver.Sender)
		if isSender {
			if err := d.RegistrySender(name, sender); err != nil {
				return err //@TODO
			}
		}
		//If driver isn't sender and isn't listener, what is it?!
		if !isSender && !isListener {
			return fmt.Errorf(`%s %w`, name, ErrNotImplementAnyOfWaysInterfaces)
		}
	}
	return nil
}

//Repository for drivers
type DriverService struct {
	listeners map[string]driver.Listener
	senders   map[string]driver.Sender
	//Loader    driver.ConfigLoader
}

func (d *DriverService) GetListener(name string) (driver.Listener, error) {
	listener, ok := d.listeners[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownListenerDriver, name)
	}
	return listener, nil
}

func (d *DriverService) GetSender(name string) (driver.Sender, error) {
	sender, ok := d.senders[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownSenderDriver, name)
	}
	return sender, nil
}

func (d *DriverService) RegistryListener(name string, listener driver.Listener) error {
	if d.listeners == nil {
		d.listeners = make(map[string]driver.Listener)
	}
	if _, exists := d.listeners[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.listeners[name] = listener
	return nil
}

func (d *DriverService) RegistrySender(name string, sender driver.Sender) error {
	if d.senders == nil {
		d.senders = make(map[string]driver.Sender)
	}
	if _, exists := d.senders[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.senders[name] = sender
	return nil
}
