package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	KafkaOutputMode  = "kafka"
	DirectOutputMode = "console"
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

func validate(cfg *Config) error {
	if cfg.App.OutputMode != DirectOutputMode && cfg.App.OutputMode != KafkaOutputMode {
		return errors.New("unknown output mode")
	}
	return nil
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

	if err = validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
