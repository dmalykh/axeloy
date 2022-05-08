package driver

import (
	"context"
	"errors"
	"fmt"
	"plugin"
)

var (
	ErrOpenPlugin        = errors.New(`can't open plugin`)
	ErrPluginIsNotDriver = errors.New(`plugin is not driver`)
)

// Open and init driver's plugin
func Open(ctx context.Context, path string, loader ConfigLoader) (Driver, error) {
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf(`%w "%s": %s`, ErrOpenPlugin, path, err.Error())
	}
	drv, err := plug.Lookup("Driver")
	if err != nil {
		return nil, err //@TODO
	}
	d, isDriver := drv.(Driver)
	if !isDriver {
		return nil, fmt.Errorf(`"%s" %w`, path, ErrPluginIsNotDriver)
	}
	if err := d.Init(ctx, loader); err != nil {
		return nil, fmt.Errorf(`can't init plugin %s: %w`, path, err)
	}
	return d, nil
}
