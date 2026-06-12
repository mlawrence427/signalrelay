# SignalRelay

SignalRelay is a future SimpleStates concept for stale-aware local read infrastructure for externally sourced state facts.

It is not an active production project yet.

Core boundary:

> SignalRelay does not make stale state safe. It makes stale state visible before the application decides.

> Freshness is evidence. Risk tolerance is application logic.

Current status: concept / research note.

Read the concept note:

* [SignalRelay concept note](docs/signalrelay.md)

SignalRelay is not a service mesh, feature flag system, policy engine, workflow engine, webhook platform, generic cache, or authorization layer.

It is intended to explore local availability for externally sourced state facts while keeping application decisions application-owned.

## Local development

Run the local service:

```sh
go run ./cmd/signalrelay
```

Write a Stripe subscription state envelope:

```sh
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

```sh
curl "http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_123"
```
