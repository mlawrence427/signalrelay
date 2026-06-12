# SignalRelay

SignalRelay is an early local prototype for stale-aware Stripe subscription state envelopes.

It is not production software.

The `POST /v1/stripe/events` endpoint is unsigned demo ingestion only.

Real Stripe webhook signature verification is not implemented yet.

Do not expose the demo ingestion endpoint to the public internet.

Core boundary:

> SignalRelay does not make stale state safe. It makes stale state visible before the application decides.

> Freshness is evidence. Risk tolerance is application logic.

## Current capabilities

SignalRelay currently includes:

* local HTTP service
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion
* duplicate demo event protection by `source_event_id`
* local read API by `customer_id`
* dynamic freshness computation
* in-memory store by default
* optional SQLite persistence
* examples and PowerShell smoke scripts

## Documentation

* [SignalRelay concept note](docs/signalrelay.md)
* [Local API](docs/api.md)
* [Prototype checkpoint](docs/checkpoints/prototype-1.md)

## Non-capabilities

SignalRelay does not:

* evaluate authorization
* decide access
* return `allowed` or `denied`
* enforce policy
* orchestrate workflows
* retry application actions
* call Stripe APIs
* verify Stripe signatures yet
* expose a production webhook surface
* replace Stripe as the source of truth
* integrate with StateMirror yet

SignalRelay is not a service mesh, feature flag system, policy engine, workflow engine, webhook platform, generic cache, or authorization layer.

It is intended to explore local availability for externally sourced state facts while keeping application decisions application-owned.

## Local development

Run the local service:

```bash
go run ./cmd/signalrelay
```

The default local prototype store is in memory.

Runtime configuration:

```bash
SIGNALRELAY_ADDR=:8080
SIGNALRELAY_STORE=memory
SIGNALRELAY_DB_PATH=signalrelay.db
SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS=300
```

To use optional SQLite persistence:

```bash
SIGNALRELAY_ADDR=:8080 SIGNALRELAY_STORE=sqlite SIGNALRELAY_DB_PATH=signalrelay.db go run ./cmd/signalrelay
```

On Windows PowerShell:

```powershell
$env:SIGNALRELAY_ADDR=":8080"
$env:SIGNALRELAY_STORE="sqlite"
$env:SIGNALRELAY_DB_PATH="signalrelay.db"
go run ./cmd/signalrelay
```

CLI inspection:

```bash
go run ./cmd/signalrelay --version
go run ./cmd/signalrelay --help
```

Available local endpoints:

* `GET /healthz`
* `POST /v1/stripe/subscription-state`
* `POST /v1/stripe/events`
* `GET /v1/state/stripe/subscription?customer_id=...`

Write a Stripe subscription state envelope:

```bash
curl -X POST http://localhost:8080/v1/stripe/subscription-state \
  -H "Content-Type: application/json" \
  -d '{
    "source": "stripe",
    "subject": "cus_123",
    "state_type": "subscription",
    "observed_at": "2026-06-11T18:00:00Z",
    "stale_after": "2099-01-01T00:00:00Z",
    "source_event_id": "evt_123",
    "source_object_id": "sub_123",
    "payload": {
      "customer_id": "cus_123",
      "subscription_id": "sub_123",
      "status": "active"
    }
  }'
```

Read the local state envelope:

```bash
curl "http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_123"
```

The prototype also accepts unsigned demo Stripe-shaped subscription events at `POST /v1/stripe/events`. Real Stripe webhook signature verification is not implemented.

Repeated demo Stripe events with the same event id are treated as duplicates and do not rewrite the stored envelope. This is duplicate event protection, not workflow retry orchestration.

## Examples and smoke tests

The `examples/` directory includes static Stripe subscription envelopes. Refresh `observed_at` and `stale_after` if you need an example response to report `fresh`.

Memory mode:

```powershell
go run ./cmd/signalrelay
.\scripts\smoke-memory.ps1
```

SQLite mode:

```powershell
$env:SIGNALRELAY_STORE="sqlite"
$env:SIGNALRELAY_DB_PATH="signalrelay-dev.db"
go run ./cmd/signalrelay
.\scripts\smoke-sqlite.ps1
```
