package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	StoreMemory = "memory"
	StoreSQLite = "sqlite"
)

type Config struct {
	Addr                            string
	Store                           string
	DBPath                          string
	StripeStaleAfterSeconds         int
	StripeWebhookSecret             string
	StripeSignatureToleranceSeconds int
}

func Load() (Config, error) {
	cfg := Config{
		Addr:                            getEnv("SIGNALRELAY_ADDR", ":8080"),
		Store:                           strings.ToLower(getEnv("SIGNALRELAY_STORE", StoreMemory)),
		DBPath:                          getEnv("SIGNALRELAY_DB_PATH", "signalrelay.db"),
		StripeStaleAfterSeconds:         300,
		StripeWebhookSecret:             os.Getenv("SIGNALRELAY_STRIPE_WEBHOOK_SECRET"),
		StripeSignatureToleranceSeconds: 300,
	}

	staleAfterSeconds, err := parsePositiveIntEnv("SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}
	cfg.StripeStaleAfterSeconds = staleAfterSeconds

	signatureToleranceSeconds, err := parsePositiveIntEnv("SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}
	cfg.StripeSignatureToleranceSeconds = signatureToleranceSeconds

	switch cfg.Store {
	case StoreMemory, StoreSQLite:
		return cfg, nil
	default:
		return Config{}, StoreError(cfg.Store)
	}
}

func (c Config) StripeStaleAfterDuration() time.Duration {
	return time.Duration(c.StripeStaleAfterSeconds) * time.Second
}

func (c Config) StripeSignatureToleranceDuration() time.Duration {
	return time.Duration(c.StripeSignatureToleranceSeconds) * time.Second
}

func StoreError(store string) error {
	return fmt.Errorf("unknown SIGNALRELAY_STORE %q (expected memory or sqlite)", store)
}

func getEnv(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}

func parsePositiveIntEnv(name string, fallback int) (int, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}

	return parsed, nil
}
