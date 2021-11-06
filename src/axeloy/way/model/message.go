package model

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way/driver"
)

type Handler func(ctx context.Context, m Message) error

type Message struct {
	Payload      payload.Payload
	Publisher    profile.Fields
	Destinations []driver.Destination
}

func (m *Message) GetPayload() payload.Payload {
	return m.Payload
}

func (m *Message) GetPublisher() profile.Fields {
	return m.Publisher
}

func (m *Message) GetDestinations() []driver.Destination {
	return m.Destinations
}

type Destination struct {
	Way      string
	Consumer profile.Fields
}

func (d *Destination) GetWay() string {
	return d.Way
}

func (d *Destination) GetConsumer() profile.Fields {
	return d.Consumer
}
