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

// ConfigLoader using for load configs. For example:
//	type Zero struct {
//		ListenAddr string `hcl:"listen_addr"`
//	}
//	var k Zero
//
//	And unmarshal driver's config to k struct:
//
// 	err := loader(&k)
type ConfigLoader func(v interface{}) error

//type Config config.DriverConfig //@TODO: SOLID

// Driver used for work with Axeloy
type Driver interface {
	//The ValidateProfile method validates profile when route is created and before message been sent
	ValidateProfile(ctx context.Context, p profile.Profile) error
	// SetWayParams sets params  for driver from ways storage when driver loads
	SetWayParams(params Params)
	// GetWayParamsFields returns fields that describes params of way
	GetWayParamsFields() ParamsFields
	// Init calls when Axeloy starts. Loader used for driver's config loading, see ConfigLoader
	Init(ctx context.Context, loader ConfigLoader) error
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
