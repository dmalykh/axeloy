package model

import "github.com/dmalykh/axeloy/axeloy/way/driver"

type Listener struct {
	Way
	driver.Listener
}

func MakeListener(way Way, listener driver.Listener) *Listener {
	listener.SetWayParams(way.GetParams())
	return &Listener{way, listener}
}
