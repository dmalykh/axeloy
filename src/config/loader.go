package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type dbconfig struct {
	Driver string `hcl:"driver,label"`
	Dsn    string `hcl:"dsn"`
}

type DriverConfig struct {
	Name   string   `hcl:"name,label"`
	Path   string   `hcl:"path,label"`
	Config hcl.Body `hcl:",remain"`
}

type Config struct {
	Database dbconfig       `hcl:"database,block"`
	Drivers  []DriverConfig `hcl:"driver,block"`
}

func LoadFile(path string) (*Config, error) {
	var c Config
	if err := hclsimple.DecodeFile(path, nil, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func Unmarshal(dc DriverConfig, v interface{}) error {
	if err := gohcl.DecodeBody(dc.Config, nil, v); err != nil {
		return err
	}
	return nil
}
