//go:generate reform

package model

import "github.com/google/uuid"

//reform:route_message
type MessageRoute struct {
	MessageId uuid.UUID `reform:"id,pk"`
	RouteId   uuid.UUID
}
