package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisAddr   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	cfg := &Config{
		HTTPAddr:    os.Getenv("HTTP_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}

	if cfg.HTTPAddr == "" {
		return nil, errors.New("HTTP_ADDR is required")
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.RedisAddr == "" {
		return nil, errors.New("REDIS_ADDR is required")
	}

	return cfg, nil
}
