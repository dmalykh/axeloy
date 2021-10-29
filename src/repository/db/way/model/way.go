//go:generate reform
package model

import "github.com/google/uuid"

//reform:way_way
type Way struct {
	Id         uuid.UUID `reform:"id,pk"`
	Name       string    `reform:"name"`
	Title      string    `reform:"title"`
	Type       string    `reform:"type"`
	DriverName string    `reform:"driver_name"`
	Params     string    `reform:"params"`
}
