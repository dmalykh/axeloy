package model

import (
	"github.com/google/uuid"
)

type WayType string

const Listener WayType = "listener"
const Sender WayType = "sender"

type Params map[string]interface{}

type Way struct {
	Id         uuid.UUID
	Name       string
	Type       WayType
	DriverName string
	Params     Params
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

func (w *Way) GetParams() Params {
	return w.Params
}

func (w *Way) GetName() string {
	return w.Name
}
