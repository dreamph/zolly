package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type GatewayConfig struct {
	Server   Server             `yaml:"server"`
	Plugins  []PluginDefinition `yaml:"plugins"`
	Services []ServiceConfig    `yaml:"services"`
}

type PluginDefinition struct {
	Name     string            `yaml:"name"`
	Path     string            `yaml:"path"`
	Settings map[string]string `yaml:"settings"`
}

type Server struct {
	Port      string      `yaml:"port"`
	BodyLimit int         `yaml:"bodyLimit"`
	SSL       *SSLConfig  `yaml:"ssl"`
	Log       *LogConfig  `yaml:"log"`
	Cors      *CorsConfig `yaml:"cors"`
}

type SSLConfig struct {
	Enable      bool               `yaml:"enable"`
	GenerateKey *GenerateKeyConfig `yaml:"generateKey"`
	Key         *KeyConfig         `yaml:"key"`
}

type GenerateKeyConfig struct {
	Enable    bool                   `yaml:"enable"`
	KeyConfig *GenerateKeyConfigInfo `yaml:"keyConfig"`
}

type GenerateKeyConfigInfo struct {
	CommonName string `yaml:"commonName"`
	File       string `yaml:"file"`
	Password   string `yaml:"password"`
}

type KeyConfig struct {
	File     string `yaml:"file"`
	Password string `yaml:"password"`
}

type LogConfig struct {
	Enable bool `yaml:"enable"`
}

type CorsConfig struct {
	Enable bool `yaml:"enable"`
}

type ServiceConfig struct {
	Path      string   `yaml:"path"`
	Timeout   int      `yaml:"timeout"`
	Servers   []string `yaml:"servers"`
	StripPath bool     `yaml:"stripPath"`
	Plugins   []string `yaml:"plugins"`
}

func LoadConfig(filename string) (*GatewayConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config GatewayConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
