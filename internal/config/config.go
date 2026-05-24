package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type OIDCProvider struct {
	Name         string
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	PostLoginURL string // configurable per provider
}

type Config struct {
	Addr            string
	DatabaseURL     string
	SessionLifetime time.Duration
	SessionSecure   bool
	SessionCookie   string
	OIDCProvider    map[string]OIDCProvider
}

func Load() *Config {
	return &Config{
		Addr:            getEnv("APP_ADDR", ":8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/app?sslmode=disable"),
		SessionLifetime: getEnvDuration("SESSION_LIFETIME_HOURS", 24) * time.Hour,
		SessionSecure:   getEnvBool("SESSION_COOKIE_SECURE", true),
		SessionCookie:   getEnv("SESSION_COOKIE_NAME", "webbuilder_session"),
		OIDCProvider:    loadOIDCProviders(),
	}
}

func loadOIDCProviders() map[string]OIDCProvider {
	providers := make(map[string]OIDCProvider)

	// reads OIDC_<NAME>_* env vars, supports multiple providers
	for name := range strings.SplitSeq(getEnv("OIDC_PROVIDERS", ""), ",") {
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			continue
		}
		prefix := "OIDC_" + strings.ToUpper(name) + "_"
		providers[name] = OIDCProvider{
			Name:         name,
			IssuerURL:    getEnv(prefix+"ISSUER_URL", ""),
			ClientID:     getEnv(prefix+"CLIENT_ID", ""),
			ClientSecret: getEnv(prefix+"CLIENT_SECRET", ""),

			RedirectURL:  getEnv(prefix+"REDIRECT_URL", ""),
			Scopes:       strings.Split(getEnv(prefix+"SCOPES", "openid,email,profile"), ","),
			PostLoginURL: getEnv(prefix+"POST_LOGIN_URL", "/"),
		}
	}
	return providers
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
