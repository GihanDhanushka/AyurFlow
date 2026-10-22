package config

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Environment string
	HTTPAddr    string
	DatabaseURL string
	Auth        AuthPolicy
}

// AuthPolicy is deliberately unset until DEC07 values are supplied by the operator.
type AuthPolicy struct {
	Enabled         bool
	IdleTimeout     time.Duration
	AbsoluteTimeout time.Duration
	AttemptLimit    int
	AttemptWindow   time.Duration
	BlockDuration   time.Duration
	AllowedOrigin   string
}

func Load() (Config, error) {
	cfg := Config{
		Environment: valueOrDefault("APP_ENV", "development"),
		HTTPAddr:    valueOrDefault("HTTP_ADDR", "127.0.0.1:8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	policy, err := loadAuthPolicy(os.Getenv)
	if err != nil {
		return Config{}, err
	}
	cfg.Auth = policy
	return cfg, nil
}

func loadAuthPolicy(getenv func(string) string) (AuthPolicy, error) {
	if getenv("AUTH_ENABLED") == "" {
		return AuthPolicy{}, nil
	}
	if getenv("AUTH_ENABLED") != "true" {
		return AuthPolicy{}, errors.New("AUTH_ENABLED must be true or unset")
	}
	keys := []string{
		"AUTH_IDLE_TIMEOUT", "AUTH_ABSOLUTE_TIMEOUT", "AUTH_ATTEMPT_WINDOW", "AUTH_BLOCK_DURATION",
	}
	values := make([]time.Duration, len(keys))
	for i, key := range keys {
		value, err := time.ParseDuration(getenv(key))
		if err != nil || value <= 0 {
			return AuthPolicy{}, errors.New(key + " must be a positive duration")
		}
		values[i] = value
	}
	limit, err := strconv.Atoi(getenv("AUTH_ATTEMPT_LIMIT"))
	if err != nil || limit < 1 || limit > 100 {
		return AuthPolicy{}, errors.New("AUTH_ATTEMPT_LIMIT must be between 1 and 100")
	}
	if values[0] > values[1] {
		return AuthPolicy{}, errors.New("AUTH_IDLE_TIMEOUT cannot exceed AUTH_ABSOLUTE_TIMEOUT")
	}
	origin := getenv("AUTH_ALLOWED_ORIGIN")
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
		parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.String() != origin {
		return AuthPolicy{}, errors.New("AUTH_ALLOWED_ORIGIN must be an HTTPS origin")
	}
	return AuthPolicy{
		Enabled: true, IdleTimeout: values[0], AbsoluteTimeout: values[1],
		AttemptLimit: limit, AttemptWindow: values[2], BlockDuration: values[3], AllowedOrigin: origin,
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
