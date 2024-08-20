package main

import (
	"github.com/goccy/go-yaml"
	"os"
)

type GatewayConfig struct {
	Server   Server          `yaml:"server"`
	Services []ServiceConfig `yaml:"services"`
}

type Server struct {
	Port string     `yaml:"port"`
	SSL  *SSLConfig `yaml:"ssl"`
}

type SSLConfig struct {
	Enable  bool   `yaml:"enable"`
	KeyType string `yaml:"keyType"`

	P12KeyFile     string `yaml:"p12KeyFile"`
	P12KeyPassword string `yaml:"p12KeyPassword"`
}

type ServiceConfig struct {
	Path      string   `yaml:"path"`
	Timeout   int      `yaml:"timeout"`
	Servers   []string `yaml:"servers"`
	StripPath bool     `yaml:"stripPath"`
}

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func LoadConfig(filePath string) (*GatewayConfig, error) {
	ymlData, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cfg := &GatewayConfig{}
	err = yaml.Unmarshal(ymlData, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
