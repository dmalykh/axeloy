//go:generate reform
package model

import (
	"github.com/google/uuid"
	"time"
)

//reform:route_track_attempt
type Attempt struct {
	Id         uuid.UUID  `reform:"id,pk"`
	TrackId    uuid.UUID  `reform:"track_id"`
	StartedAt  time.Time  `reform:"started_at"`
	FinishedAt *time.Time `reform:"finished_at"`
	Status     string     `reform:"status"`
	Info       string     `reform:"info"`
}
