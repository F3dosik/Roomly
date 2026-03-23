package config

import (
	"fmt"
	"log"

	"github.com/avito-internships/test-backend-1-F3dosik/internal/logger"
	"github.com/caarlos0/env/v6"
)

const (
	defaultLogLevel = string(logger.ModeDevelopment)
)

type Config struct {
	ServerPort  string `env:"SERVER_PORT"`
	DatabaseURL string `env:"DATABASE_URL"`
	LogLevel    string `env:"LOG_LEVEL"`
	JWTSecret   string `env:"JWT_SECRET"`
}

func Load() (*Config, error) {
	config := parseConfig()

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func parseConfig() *Config {
	var config Config
	err := env.Parse(&config)
	if err != nil {
		log.Printf("Warning: failed to parse env config: %v\n", err)
	}

	if config.LogLevel == "" {
		config.LogLevel = defaultLogLevel
	}

	return &config
}

func (c *Config) Validate() error {
	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("database address can't be empty")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT secret can't be empty")
	}

	switch c.LogLevel {
	case string(logger.ModeDevelopment), string(logger.ModeProduction):
	default:
		return fmt.Errorf("invalid log mode: %s, allowed: development, production", c.LogLevel)
	}

	return nil
}
