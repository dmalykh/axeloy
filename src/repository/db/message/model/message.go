//go:generate reform
package model

import (
	"github.com/google/uuid"
	"time"
)

//reform:message_message
type Message struct {
	Id            uuid.UUID `reform:"id,pk"`
	SourceWayName string    `reform:"source_way_name"`
	CreatedAt     time.Time `reform:"created_at"`
	Info          string    `reform:"info"`
	Status        string    `reform:"status"`
}
