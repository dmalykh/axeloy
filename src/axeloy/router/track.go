package router

import (
	"github.com/dmalykh/axeloy/axeloy/message"
	"github.com/dmalykh/axeloy/axeloy/profile"
	"github.com/dmalykh/axeloy/axeloy/way"
	"github.com/google/uuid"
)

type TrackStatus string

type Track interface {
	GetId() uuid.UUID
	GetSender() way.Sender
	GetMessage() message.Message
	GetProfile() profile.Profile
	GetStatus() TrackStatus
}
