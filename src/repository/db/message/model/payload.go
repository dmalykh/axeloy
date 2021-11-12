//go:generate reform
package model

import (
	"github.com/google/uuid"
)

//reform:message_payload
type Payload struct {
	Id        uuid.UUID `reform:"id,pk"`
	MessageId uuid.UUID `reform:"message_id"`
	Key       string    `reform:"key"`
	Value     string    `reform:"value"`
}
