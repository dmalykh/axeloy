//go:generate reform

package model

import (
	"github.com/google/uuid"
)

//reform:route_route
type Track struct {
	Id        uuid.UUID `reform:"id,pk"`
	WayId     uuid.UUID
	MessageId uuid.UUID
}
