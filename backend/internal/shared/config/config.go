package config

import (
	"errors"
	"os"
	"time"

	"smarterp/backend/internal/shared/storage"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	AccessTTL   time.Duration
	R2          storage.R2Options
}

func Load() (Config, error) {
	cfg := Config{}
	cfg.HTTPAddr = getenv("HTTP_ADDR", ":8080")
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.AccessTTL = parseDuration("JWT_ACCESS_TTL", "8760h")
	cfg.R2 = loadR2Options()
	return cfg, validate(cfg)
}

func loadR2Options() storage.R2Options {
	return storage.R2Options{
		AccountID:     os.Getenv("R2_ACCOUNT_ID"),
		Bucket:        os.Getenv("R2_BUCKET"),
		AccessKeyID:   os.Getenv("R2_ACCESS_KEY_ID"),
		SecretKey:     os.Getenv("R2_SECRET_ACCESS_KEY"),
		PublicBaseURL: os.Getenv("R2_PUBLIC_BASE_URL"),
	}
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
