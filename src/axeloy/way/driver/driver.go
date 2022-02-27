package driver

import (
	"context"
	"errors"
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
)

var ErrStopped = errors.New(`service stopped`)

type Params map[string][]string

type ParamsField struct {
	Name  string
	Title string
}

type ParamsFields []ParamsField

type DriverConfig interface{}

type Config struct {
	Path   string
	Config DriverConfig
}

type Driver interface {
	//The ValidateProfile method validates profile when route is created and before message been sent
	ValidateProfile(ctx context.Context, p profile.Profile) error
	// SetWayParams sets params  for driver from ways storage when driver loads
	SetWayParams(params Params)
	// GetWayParamsFields returns fields that describes params of way
	GetWayParamsFields() ParamsFields
	// SetDriverConfig sets config for driver when app starts
	SetDriverConfig(config DriverConfig)
	Stop(ctx context.Context) error
}

type Sender interface {
	Driver
	Send(ctx context.Context, recipient profile.Profile, message message.Message) ([]string, error)
}

type Listener interface {
	Driver
	Listen(context.Context, func(ctx context.Context, message Message) error) error
}

func UnmarshalParams(p Params, v interface{}) error {
	//@TODO
	return nil
}

func MakeMessage() {

}
