package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	LogLevel   string `json:"log_level"`
	BindAddr   string `json:"bind_addr"`
	configPath string
}

func NewConfig() *Config {
	return &Config{
		LogLevel:   "debug",
		configPath: "configs/pr-reviewer-service.json",
	}
}

func (cfg *Config) ReadConfig() error {
	data, err := os.ReadFile(cfg.configPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}

	return nil
}
