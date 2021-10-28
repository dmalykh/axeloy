package model

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/google/uuid"
)

type TrackStatus string

const (
	Planned TrackStatus = `planned`
	Process TrackStatus = `process`
	Error   TrackStatus = `error`
	Done    TrackStatus = `done`
)

type Track struct {
	Id        uuid.UUID
	SenderId  uuid.UUID
	MessageId uuid.UUID
	Profile   profile.Profile
	Attempts  int
	Info      string
	Status    TrackStatus
}

func (t *Track) GetId() uuid.UUID {
	return t.Id
}

func (t *Track) GetSenderId() uuid.UUID {
	return t.SenderId
}

func (t *Track) GetMessageId() uuid.UUID {
	return t.MessageId
}

func (t *Track) GetProfile() profile.Profile {
	return t.Profile
}

func (t *Track) GetAttempts() int {
	return t.Attempts
}

func (t *Track) GetStatus() TrackStatus {
	return t.Status
}
