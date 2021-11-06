//go:generate reform

package model

import (
	"github.com/google/uuid"
)

//reform:route_route
type Route struct {
	Id uuid.UUID `reform:"id,pk"`
}
