//go:generate reform

package model

import (
	"github.com/google/uuid"
)

type Type string

const (
	Source      Type = "source"
	Destination Type = "destination"
)

//reform:route_profile
type Profile struct {
	Id      uuid.UUID `reform:"id,pk"`
	RouteId uuid.UUID `reform:"route_id"`
	Type    Type      `reform:"enum"`
	Key     string    `reform:"key"`
	Value   string    `reform:"value"`
}
