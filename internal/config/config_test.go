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
	if cfg.StripeStaleAfterSeconds != 300 {
		t.Fatalf("StripeStaleAfterSeconds = %d, want %d", cfg.StripeStaleAfterSeconds, 300)
	}
	if cfg.StripeWebhookSecret != "" {
		t.Fatalf("StripeWebhookSecret = %q, want empty", cfg.StripeWebhookSecret)
	}
	if cfg.StripeSignatureToleranceSeconds != 300 {
		t.Fatalf("StripeSignatureToleranceSeconds = %d, want %d", cfg.StripeSignatureToleranceSeconds, 300)
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

func TestLoadCustomStripeStaleAfterSeconds(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", "")
	t.Setenv("SIGNALRELAY_DB_PATH", "")
	t.Setenv("SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS", "60")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.StripeStaleAfterSeconds != 60 {
		t.Fatalf("StripeStaleAfterSeconds = %d, want %d", cfg.StripeStaleAfterSeconds, 60)
	}
}

func TestLoadInvalidStripeStaleAfterSecondsReturnsError(t *testing.T) {
	cases := []string{"0", "-1", "not-a-number"}

	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			t.Setenv("SIGNALRELAY_ADDR", "")
			t.Setenv("SIGNALRELAY_STORE", "")
			t.Setenv("SIGNALRELAY_DB_PATH", "")
			t.Setenv("SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS", value)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
		})
	}
}

func TestLoadCustomStripeWebhookSecret(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", "")
	t.Setenv("SIGNALRELAY_DB_PATH", "")
	t.Setenv("SIGNALRELAY_STRIPE_WEBHOOK_SECRET", "whsec_test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.StripeWebhookSecret != "whsec_test" {
		t.Fatalf("StripeWebhookSecret = %q, want %q", cfg.StripeWebhookSecret, "whsec_test")
	}
}

func TestLoadCustomStripeSignatureToleranceSeconds(t *testing.T) {
	t.Setenv("SIGNALRELAY_ADDR", "")
	t.Setenv("SIGNALRELAY_STORE", "")
	t.Setenv("SIGNALRELAY_DB_PATH", "")
	t.Setenv("SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS", "60")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.StripeSignatureToleranceSeconds != 60 {
		t.Fatalf("StripeSignatureToleranceSeconds = %d, want %d", cfg.StripeSignatureToleranceSeconds, 60)
	}
}

func TestLoadInvalidStripeSignatureToleranceSecondsReturnsError(t *testing.T) {
	cases := []string{"0", "-1", "not-a-number"}

	for _, value := range cases {
		t.Run(value, func(t *testing.T) {
			t.Setenv("SIGNALRELAY_ADDR", "")
			t.Setenv("SIGNALRELAY_STORE", "")
			t.Setenv("SIGNALRELAY_DB_PATH", "")
			t.Setenv("SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS", value)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
		})
	}
}
