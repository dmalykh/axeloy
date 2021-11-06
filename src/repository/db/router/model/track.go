//go:generate reform

package model

import (
	"github.com/google/uuid"
)

//reform:route_track
type Track struct {
	Id        uuid.UUID `reform:"id,pk"`
	WayName   string    `reform:"way_name"`
	MessageId uuid.UUID `reform:"message_id"`
	Attempts  int       `reform:"attempts"`
	Info      string    `reform:"info"`
	Status    string    `reform:"status"`
}
