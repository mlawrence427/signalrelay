# SignalRelay

SignalRelay is a SimpleStates concept and early local prototype for stale-aware local read infrastructure for externally sourced state facts.

Core boundary:

> SignalRelay does not make stale state safe. It makes stale state visible before the application decides.

> Freshness is evidence. Risk tolerance is application logic.

Current status: early local prototype.

SignalRelay currently includes a minimal local HTTP service that can store and return stale-aware Stripe subscription state envelopes in memory.

It is not production software.

Read the concept note:

* [SignalRelay concept note](docs/signalrelay.md)

SignalRelay is not a service mesh, feature flag system, policy engine, workflow engine, webhook platform, generic cache, or authorization layer.

It is intended to explore local availability for externally sourced state facts while keeping application decisions application-owned.

The prototype does not evaluate authorization, decide access, return `allowed` or `denied`, enforce policy, or orchestrate workflows.

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

Available local endpoints:

* `GET /healthz`
* `POST /v1/stripe/subscription-state`
* `GET /v1/state/stripe/subscription?customer_id=...`

API contract:

* [Local API](docs/api.md)

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
