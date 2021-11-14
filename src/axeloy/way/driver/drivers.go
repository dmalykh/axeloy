package driver

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
type Drivers struct {
	listeners map[string]Listener
	senders   map[string]Sender
}

func (d *Drivers) GetListener(name string) (Listener, error) {
	listener, ok := d.listeners[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownListenerDriver, name)
	}
	return listener, nil
}

func (d *Drivers) GetSender(name string) (Sender, error) {
	sender, ok := d.senders[name]
	if !ok {
		return nil, fmt.Errorf(`%w %s`, ErrUnknownSenderDriver, name)
	}
	return sender, nil
}

func (d *Drivers) RegistryListener(name string, listener Listener) error {
	if d.listeners == nil {
		d.listeners = make(map[string]Listener)
	}
	if _, exists := d.listeners[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.listeners[name] = listener
	return nil
}

func (d *Drivers) RegistrySender(name string, sender Sender) error {
	if d.senders == nil {
		d.senders = make(map[string]Sender)
	}
	if _, exists := d.senders[name]; exists {
		return fmt.Errorf(`%s %w`, name, ErrDriverNameExists)
	}
	d.senders[name] = sender
	return nil
}
