package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	RedisAddr   string
	JWTSecret   string
	JWTExpires  time.Duration
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	jwtExpires := 24 * time.Hour
	if raw := os.Getenv("JWT_EXPIRES"); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JWT_EXPIRES: %w", err)
		}
		jwtExpires = parsed
	}

	cfg := &Config{
		HTTPAddr:    os.Getenv("HTTP_ADDR"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpires:  jwtExpires,
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
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	return cfg, nil
}
