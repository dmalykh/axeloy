package model

import (
	"github.com/dmalykh/axeloy/axeloy/way/driver"
	"github.com/google/uuid"
)

type WayType string

const WayTypeListener WayType = "listener"
const WayTypeSender WayType = "sender"

type Way struct {
	Id         uuid.UUID
	Name       string
	Title      string
	Type       WayType
	DriverName string
	Params     driver.Params
}

func (w *Way) GetId() uuid.UUID {
	return w.Id
}

func (w *Way) GetType() WayType {
	return w.Type
}

func (w *Way) GetDriverName() string {
	return w.DriverName
}

func (w *Way) GetParams() driver.Params {
	return w.Params
}

func (w *Way) GetName() string {
	return w.Name
}

func (w *Way) GetTitle() string {
	return w.Title
}
