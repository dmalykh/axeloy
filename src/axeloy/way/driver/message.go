package driver

import (
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/dmalykh/axeloy/axeloy/profile"
)

type Message interface {
	GetPayload() payload.Payload
	GetPublisher() profile.Fields
	GetDestinations() []Destination
}

type Destination interface {
	GetWay() string
	GetConsumer() profile.Fields
}
