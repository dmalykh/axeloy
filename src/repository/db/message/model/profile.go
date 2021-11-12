//go:generate reform

package model

import (
	"github.com/google/uuid"
)

//reform:message_profile
type Profile struct {
	Id                 uuid.UUID `reform:"id,pk"`
	MessageId          uuid.UUID `reform:"message_id"`
	DestinationWayName *string   `reform:"destination_way_name"`
	Key                string    `reform:"key"`
	Value              string    `reform:"value"`
}
