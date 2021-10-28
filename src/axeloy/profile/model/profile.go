package model

type Profile struct {
	Fields map[string][]string
}

func (p *Profile) GetFields() map[string][]string {
	return p.Fields
}
