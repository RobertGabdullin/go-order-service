package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	Kafka struct {
		Brokers []string `yaml:"brokers"`
		Topic   string   `yaml:"topic"`
		GroupID string   `yaml:"group_id"`
	} `yaml:"kafka"`
	App struct {
		OutputMode string `yaml:"output_mode"`
	} `yaml:"app"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	if cfg.App.OutputMode != "direct" && cfg.App.OutputMode != "kafka" {
		return nil, errors.New("Unknown output mode")
	}
	return &cfg, nil
}
