package config

import (
	"errors"
	"os"
	"time"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	AccessTTL   time.Duration
}

func Load() (Config, error) {
	cfg := Config{}
	cfg.HTTPAddr = getenv("HTTP_ADDR", ":8080")
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.AccessTTL = parseDuration("JWT_ACCESS_TTL", "8760h")
	return cfg, validate(cfg)
}

func validate(cfg Config) error {
	if cfg.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	return nil
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseDuration(key, fallback string) time.Duration {
	value := getenv(key, fallback)
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return time.Minute
	}
	return parsed
}
