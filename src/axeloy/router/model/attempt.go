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
	Id         uuid.UUID
	TrackId    uuid.UUID
	StartedAt  time.Time
	FinishedAt time.Time
	Status     AttemptStatus
	Info       string
}
