//go:generate reform

package model

import (
	"github.com/google/uuid"
	"time"
)

//reform:route_route
type Route struct {
	Id        uuid.UUID `reform:"id,pk"`
	CreatedAt time.Time `reform:"created_at"`
}
