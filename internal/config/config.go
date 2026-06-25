package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// *** Providers
// *********************
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	PostLoginURI string
}

type ZitadelConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	IssuerURI    string
	PostLoginURI string
}

type Config struct {
	Addr            string
	Environment     string
	DatabaseURL     string
	MigrationsDir   string
	SessionLifetime time.Duration
	SessionSecure   bool
	SessionCookie   string

	// OIDC Providers
	Google  GoogleConfig
	Zitadel ZitadelConfig
}

func Load() *Config {
	return &Config{
		Addr:            getEnv("APP_ADDR", ":8080"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/app?sslmode=disable"),
		MigrationsDir:   getEnv("MIGRATIONS_DIR", "db/migrations"),
		SessionLifetime: getEnvDuration("SESSION_LIFETIME_HOURS", 24) * time.Hour,
		SessionSecure:   getEnvBool("SESSION_COOKIE_SECURE", false),
		SessionCookie:   getEnv("SESSION_COOKIE_NAME", "session"),

		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("GOOGLE_REDIRECT_URI", "http://localhost:8080/auth/google/callback"),
			PostLoginURI: getEnv("GOOGLE_POST_LOGIN_URI", "http://localhost:5173/admin/dashboard"),
		},
		Zitadel: ZitadelConfig{
			ClientID:     getEnv("ZITADEL_CLIENT_ID", ""),
			ClientSecret: getEnv("ZITADEL_CLIENT_SECRET", ""),
			RedirectURI:  getEnv("ZITADEL_REDIRECT_URI", "http://localhost:8080/api/auth/callback/zitadel"),
			IssuerURI:    getEnv("ZITADEL_ISSUER_URI", ""),
			PostLoginURI: getEnv("ZITADEL_POST_LOGIN_URI", "http://localhost:5173/admin/dashboard"),
		},
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

// getEnvRequired checks for an environment variable and completely terminates the application if it is empty.
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("FATAL CONFIGURATION ERROR: Required environment variable %s is missing or empty.", key)
	}
	return value
}
