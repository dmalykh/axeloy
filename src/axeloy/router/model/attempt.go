package model

import (
	"github.com/google/uuid"
	"time"
)

type AttemptStatus string

const (
	AttemptStatusError      AttemptStatus = `error`
	AttemptStatusDone       AttemptStatus = `done`
	AttemptStatusInProgress AttemptStatus = `in_progress`
)

type Attempt struct {
	TrackId  uuid.UUID
	Started  time.Time
	Finished time.Time
	Status   AttemptStatus
	Info     string
}
