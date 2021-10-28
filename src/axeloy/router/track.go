package router

import (
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/router/model"
	"github.com/google/uuid"
)

type TrackStatus string

const (
	Planned TrackStatus = `planned`
	Error   TrackStatus = `error`
)

type Track interface {
	GetId() uuid.UUID
	GetSenderId() uuid.UUID
	GetMessageId() uuid.UUID
	GetProfile() profile.Profile
	GetAttempts() int
	GetStatus() model.TrackStatus
}
