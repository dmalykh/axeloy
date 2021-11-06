package model

import "github.com/dmalykh/axeloy/axeloy/way/driver"

type Sender struct {
	Way
	driver.Sender
}

func MakeSender(way Way, sender driver.Sender) *Sender {
	sender.SetWayParams(way.GetParams())
	return &Sender{way, sender}
}
