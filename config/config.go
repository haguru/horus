package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	CONFIG_PATH = "../res/config.yaml"
)

type ServiceConfig struct {
	Name     string   `yaml:"name"`
	LogLevel string   `yaml:"loglevel"`
	Database Database `yaml:"database"`
}

type Database struct {
	Host       string        `yaml:"host"`
	Port       int           `yaml:"port"`
	Name       string        `yaml:"name"`
	Collection string        `yaml:"collection"`
	Options    ServerOptions `yaml:"options"`
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
