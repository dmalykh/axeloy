//go:generate reform
package model

//reform:way_way
type Way struct {
	Name       string `reform:"name,pk"`
	Title      string `reform:"title"`
	Type       string `reform:"type"`
	DriverName string `reform:"driver_name"`
	Params     string `reform:"params"`
}
