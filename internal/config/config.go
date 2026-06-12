package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	StoreMemory = "memory"
	StoreSQLite = "sqlite"
)

type Config struct {
	Addr   string
	Store  string
	DBPath string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:   getEnv("SIGNALRELAY_ADDR", ":8080"),
		Store:  strings.ToLower(getEnv("SIGNALRELAY_STORE", StoreMemory)),
		DBPath: getEnv("SIGNALRELAY_DB_PATH", "signalrelay.db"),
	}

	switch cfg.Store {
	case StoreMemory, StoreSQLite:
		return cfg, nil
	default:
		return Config{}, StoreError(cfg.Store)
	}
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
