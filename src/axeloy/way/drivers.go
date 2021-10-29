package way

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownListenerDriver = errors.New(`unknown listener`)
	ErrUnknownSenderDriver   = errors.New(`unknown sender`)
	ErrDriverNameExists      = errors.New(`driver name already exists`)
)

//Repository for drivers
type drivers struct {
	listeners map[string]Listener
	senders   map[string]Sender
}

func (d *drivers) GetListener(name string) (Listener, error) {
	driver, ok := d.listeners[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownListenerDriver, name)
	}
	return driver, nil
}

func (d *drivers) GetSender(name string) (Sender, error) {
	driver, ok := d.senders[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownSenderDriver, name)
	}
	return driver, nil
}

func (d *drivers) RegistryListener(name string, listener Listener) error {
	if _, exists := d.listeners[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.listeners[name] = listener
	return nil
}

func (d *drivers) RegistrySender(name string, sender Sender) error {
	if _, exists := d.senders[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.senders[name] = sender
	return nil
}
