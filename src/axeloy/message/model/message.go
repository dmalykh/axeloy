package model

import (
	"github.com/dmalykh/axeloy/axeloy/message/location"
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/google/uuid"
)

type Status string

var (
	StatusNew            Status = "new"
	StatusProcessed      Status = "processed"
	StatusSent           Status = "sent"
	StatusError          Status = "error"
	StatusNoDestinations Status = "no_destinations"
)

type Message struct {
	Id          uuid.UUID
	Payload     payload.Payload
	Source      location.Location
	Destination []location.Location
	Info        []string
	Status      Status
}

func (m *Message) GetStatus() Status {
	return m.Status
}

func (m *Message) GetInfo() []string {
	return m.Info
}

func (m *Message) SetInfo(info ...string) {
	m.Info = info
}

func (m *Message) AddInfo(info ...string) {
	m.Info = append(m.Info, info...)
}

func (m *Message) GetPayload() payload.Payload {
	return m.Payload
}

func (m *Message) GetUUID() uuid.UUID {
	return m.Id
}

func (m *Message) GetSource() location.Location {
	return m.Source
}

func (m *Message) GetDestinations() []location.Location {
	return m.Destination
}
