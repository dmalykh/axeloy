//go:generate reform

package model

import (
	"github.com/google/uuid"
)

//reform:route_route
type Track struct {
	Id        uuid.UUID `reform:"id,pk"`
	WayId     uuid.UUID `reform:"way_id"`
	MessageId uuid.UUID `reform:"message_id"`
	Attempts  int       `reform:"attempts"`
	Status    string    `reform:"status"`
}
