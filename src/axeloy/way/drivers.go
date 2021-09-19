package way

import (
	"errors"
	"fmt"
)

var ErrUnknownListener = errors.New(`unknown listener`)
var ErrUnknownSender = errors.New(`unknown sender`)

//Repository for drivers
type drivers struct {
	listeners map[string]Listener
	senders   map[string]Sender
}

func (d *drivers) GetListener(name string) (Listener, error) {
	driver, ok := d.listeners[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownListener, name)
	}
	return driver, nil
}

func (d *drivers) GetSender(name string) (Sender, error) {
	driver, ok := d.senders[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownListener, name)
	}
	return driver, nil
}

var ErrDriverNameExists = errors.New(`driver name already exists`)

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
