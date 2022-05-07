package way

import (
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

type Driverer interface {
	GetSender(name string) (driver.Sender, error)
	GetListener(name string) (driver.Listener, error)
	RegistryListener(name string, listener driver.Listener) error
	RegistrySender(name string, sender driver.Sender) error
	Load(drivers map[string]driver.Driver) error
}
