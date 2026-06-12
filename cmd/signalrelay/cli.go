package main

import (
	"fmt"
	"io"
)

const version = "signalrelay dev"

func handleCLI(args []string, stdout io.Writer, stderr io.Writer) (bool, int) {
	if len(args) == 0 {
		return false, 0
	}

	switch args[0] {
	case "--version":
		fmt.Fprintln(stdout, version)
		return true, 0
	case "--help", "-h":
		fmt.Fprint(stdout, usageText())
		return true, 0
	default:
		fmt.Fprintf(stderr, "unknown flag: %s\n\n", args[0])
		fmt.Fprint(stderr, usageText())
		return true, 2
	}
}

func usageText() string {
	return `Usage:
  signalrelay [--version] [--help]

Environment:
  SIGNALRELAY_ADDR
  SIGNALRELAY_STORE
  SIGNALRELAY_DB_PATH
  SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS
  SIGNALRELAY_STRIPE_WEBHOOK_SECRET
  SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS

Local endpoints:
  GET /healthz
  POST /v1/stripe/subscription-state
  POST /v1/stripe/events
  POST /v1/stripe/webhook
  GET /v1/state/stripe/subscription?customer_id=...

Boundary:
  SignalRelay reports observed state, provenance, payload hash, and freshness metadata. It does not decide access.
`
}
