package model

import "github.com/dmalykh/axeloy/axeloy/profile"

type Profile struct {
	Fields profile.Fields
}

func (p *Profile) GetFields() profile.Fields {
	return p.Fields
}
