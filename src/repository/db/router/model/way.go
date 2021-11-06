//go:generate reform
package model

import "github.com/google/uuid"

//reform:route_way
type Way struct {
	Id      uuid.UUID `reform:"id,pk"`
	RouteId uuid.UUID `reform:"route_id"`
	WayName string    `reform:"way_name"`
}
