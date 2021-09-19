//go:generate reform
package model

import "github.com/google/uuid"

//reform:route_way
type Way struct {
	Id      uuid.UUID `reform:"id,pk"`
	RouteId uuid.UUID `reform:"route_id"`
	WayId   uuid.UUID `reform:"way_id"`
}
