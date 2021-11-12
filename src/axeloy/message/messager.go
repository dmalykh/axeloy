package message

import (
	"context"
	"github.com/dmalykh/axeloy/axeloy/message/location"
	"github.com/dmalykh/axeloy/axeloy/message/model"
	"github.com/dmalykh/axeloy/axeloy/message/payload"
	"github.com/google/uuid"
)

//type State string
//
//const (
//	StateNew            State = "new"
//	StateProcessed      State = "processed"
//	StateSent           State = "sent"
//	StateError          State = "error"
//	StateNoDestinations State = "no_destinations"
//)

type Message interface {
	GetId() uuid.UUID
	SetId(u uuid.UUID)
	GetSource() location.Location
	GetDestinations() []location.Location
	GetPayload() payload.Payload
	GetStatus() model.Status
	GetInfo() []string
	//SetInfo(info ...string)
	//AddInfo(info ...string)
}

type Messager interface {
	GetById(ctx context.Context, id uuid.UUID) (Message, error)
	Add(ctx context.Context, msg Message) error
	UpdateStatus(ctx context.Context, msg Message, status model.Status, info ...string) error
}
