# AGENTS.md

## Project status

SignalRelay is currently an early SimpleStates concept repo moving toward a minimal local HTTP service.

Do not treat SignalRelay as production software yet.

## Core boundary

SignalRelay is a stale-aware local read replica for externally sourced state facts.

SignalRelay does not:

- evaluate authorization
- decide access
- enforce policy
- retry application actions
- orchestrate workflows
- transform business meaning
- make stale state safe
- replace Stripe as the source of truth

SignalRelay only reports last observed external state + provenance + freshness.

## Architecture doctrine

State is queried.
Decisions are made.
Actions remain application-owned.

Freshness is evidence. Risk tolerance is application logic.

## MVP scope

The initial implementation is limited to Stripe subscription state.

Allowed first-slice behavior:

- local HTTP service
- health endpoint
- in-memory state envelope store
- POST a simplified Stripe subscription state envelope
- GET the latest envelope by customer_id
- compute payload_hash
- compute freshness on read

Do not add real Stripe webhook signature verification, SQLite/Pebble, StateMirror integration, auth, retries, queues, workers, or reconciliation until explicitly requested.

## Code style

Prefer boring, explicit Go code.
Use the Go standard library unless a dependency is clearly necessary.
Keep packages small and purpose-specific.
Run gofmt after editing Go files.

## Commands

Run locally:

```bash
go run ./cmd/signalrelay
```

Format:

```bash
gofmt -w .
```

Test, once tests exist:

```bash
go test ./...
```

## Documentation

Keep README and docs boundary-focused.
Do not use hype language.
Do not make outage-proof, guaranteed-availability, or correctness claims.
Do not frame SignalRelay as AI infrastructure.
