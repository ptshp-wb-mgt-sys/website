// Package config contains the configuration for the application
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadCfg() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	cfg := &Config{
		Port: os.Getenv("PORT"),
	}

	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.Port == "" {
		return fmt.Errorf("port is required")
	}

	return nil
}
