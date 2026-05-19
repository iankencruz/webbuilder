package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr            string
	DatabaseURL     string
	SessionLifetime time.Duration
	SessionSecure   bool
}

func Load() Config {
	return Config{
		Addr:            getEnv("APP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/webbuilder?sslmode=disable"),
		SessionLifetime: getEnvDuration("SESSION_LIFETIME_HOURS", 24) * time.Hour,
		SessionSecure:   getEnvBool("SESSION_COOKIE_SECURE", true),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvDuration(key string, fallback int) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return time.Duration(fallback)
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return time.Duration(fallback)
	}
	return time.Duration(n)
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return b
}
