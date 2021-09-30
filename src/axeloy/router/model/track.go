package model

import (
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/google/uuid"
)

type TrackStatus string

const (
	New TrackStatus = `new`
)

type Track struct {
	Id      uuid.UUID
	Sender  way.Sender
	Message message.Message
	Profile profile.Profile
	Status  TrackStatus
}

func (t *Track) GetId() uuid.UUID {
	panic("implement me")
}

func (t *Track) GetSender() way.Sender {
	panic("implement me")
}

func (t *Track) GetMessage() message.Message {
	panic("implement me")
}

func (t *Track) GetProfile() profile.Profile {
	panic("implement me")
}

func (t *Track) GetStatus() TrackStatus {
	panic("implement me")
}
