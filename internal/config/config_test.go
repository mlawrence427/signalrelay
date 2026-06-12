package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", "")
	t.Setenv("SIGNALRELAY_DB_PATH", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Addr != ":8080" {
		t.Fatalf("Addr = %q, want %q", cfg.Addr, ":8080")
	}
	if cfg.Store != StoreMemory {
		t.Fatalf("Store = %q, want %q", cfg.Store, StoreMemory)
	}
	if cfg.DBPath != "signalrelay.db" {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, "signalrelay.db")
	}
}

func TestLoadCustomAddr(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "127.0.0.1:9090")
	t.Setenv("SIGNALRELAY_STORE", "")
	t.Setenv("SIGNALRELAY_DB_PATH", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Addr != "127.0.0.1:9090" {
		t.Fatalf("Addr = %q, want %q", cfg.Addr, "127.0.0.1:9090")
	}
}

func TestLoadSQLiteWithDefaultDBPath(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", StoreSQLite)
	t.Setenv("SIGNALRELAY_DB_PATH", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Store != StoreSQLite {
		t.Fatalf("Store = %q, want %q", cfg.Store, StoreSQLite)
	}
	if cfg.DBPath != "signalrelay.db" {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, "signalrelay.db")
	}
}

func TestLoadSQLiteWithCustomDBPath(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", StoreSQLite)
	t.Setenv("SIGNALRELAY_DB_PATH", "custom.db")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DBPath != "custom.db" {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, "custom.db")
	}
}

func TestLoadUnknownStoreReturnsError(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", "other")
	t.Setenv("SIGNALRELAY_DB_PATH", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}
