package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db struct {
		Driver string `yaml:"driver"`
		Dsn    string `yaml:"dsn"`
	} `yaml:"database"`
	Ways struct {
		Drivers map[string]struct {
			DriverPath   string      `yaml:"path"`
			DriverConfig interface{} `yaml:"config"`
		} `yaml:"drivers"`
	} `yaml:"ways"`
}

func Load(path string) (*Config, error) {
	var c Config
	if err := cleanenv.ReadConfig(path, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
