package way

import (
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	"github.com/dmalykh/axeloy/axeloy/way/repository"
)

type Config struct {
	WayRepository repository.WayRepository
	Drivers       map[string]driver.Config
}
