package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	LogLevel           string        `json:"log_level"`
	BindAddr           string        `json:"bind_addr"`
	ReadHandlerTimeout time.Duration `json:"read_handler_timeout"`
	ReadTimeout        time.Duration `json:"read_timeout"`
	WriteTimeout       time.Duration `json:"write_timeout"`
	IdleTimeout        time.Duration `json:"idle_timeout"`
	configPath         string
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
