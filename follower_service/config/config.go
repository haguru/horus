package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	CONFIG_PATH = "./res/config.yaml"
)

type ServiceConfig struct {
	Name     string   `yaml:"name" validate:"required"`
	LogLevel string   `yaml:"loglevel" validate:"required"`
	Port     int      `yaml:"port" validate:"required"`
	Database Database `yaml:"database" validate:"required"`
	Metrics  Metrics  `yaml:"metrics" validate:"required"`
}

type Database struct {
	Host         string        `yaml:"host" validate:"required"`
	Port         int           `yaml:"port" validate:"required"`
	Name         string        `yaml:"name" validate:"required"`
	Timeout      string        `yaml:"timeout" validate:"required"`
	PingInterval string        `yaml:"ping_interval" validate:"required"`
	Collection   string        `yaml:"collection" validate:"required"`
	Options      ServerOptions `yaml:"options" validate:"required"`
}

type Metrics struct {
	Port int `yaml:"port" validate:"required"`
}

type ServerOptions struct {
	SetStrict            bool `yaml:"setstrict"`
	SetDeprecationErrors bool `yaml:"setdeprecationerrors"`
}

func ReadLocalConfig(configPath string) (*ServiceConfig, error) {
	config := &ServiceConfig{}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
