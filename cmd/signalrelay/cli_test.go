package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestHandleCLIVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	handled, code := handleCLI([]string{"--version"}, &stdout, &stderr)
	if !handled {
		t.Fatal("handled = false, want true")
	}
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if got := strings.TrimSpace(stdout.String()); got != version {
		t.Fatalf("stdout = %q, want %q", got, version)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestHandleCLIHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	handled, code := handleCLI([]string{"--help"}, &stdout, &stderr)
	if !handled {
		t.Fatal("handled = false, want true")
	}
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}

	out := stdout.String()
	for _, want := range []string{
		"Usage:",
		"SIGNALRELAY_ADDR",
		"SIGNALRELAY_STORE",
		"SIGNALRELAY_DB_PATH",
		"SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS",
		"GET /healthz",
		"POST /v1/stripe/subscription-state",
		"POST /v1/stripe/events",
		"GET /v1/state/stripe/subscription?customer_id=...",
		"does not decide access",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("help output missing %q:\n%s", want, out)
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestHandleCLIUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	handled, code := handleCLI([]string{"--bogus"}, &stdout, &stderr)
	if !handled {
		t.Fatal("handled = false, want true")
	}
	if code == 0 {
		t.Fatal("code = 0, want nonzero")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}

	errOut := stderr.String()
	if !strings.Contains(errOut, "unknown flag: --bogus") {
		t.Fatalf("stderr missing unknown flag message:\n%s", errOut)
	}
	if !strings.Contains(errOut, "Usage:") {
		t.Fatalf("stderr missing usage:\n%s", errOut)
	}
}

func TestHandleCLINoArgs(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	handled, code := handleCLI(nil, &stdout, &stderr)
	if handled {
		t.Fatal("handled = true, want false")
	}
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
}
