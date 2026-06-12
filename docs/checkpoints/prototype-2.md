# SignalRelay Prototype Checkpoint 2

SignalRelay is now a runnable local prototype for stale-aware Stripe subscription state envelopes, with memory and SQLite storage, demo Stripe-shaped ingestion, duplicate event protection, stricter event-shape validation, CLI inspection, and local container support.

It remains unsigned demo ingestion only and is not production software.

## Current capabilities

SignalRelay currently includes:

* local HTTP service
* configurable listen address
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion
* supported demo subscription event types
* stricter demo event-shape validation
* duplicate demo event protection by `source_event_id`
* local read API by `customer_id`
* `payload_hash` computation
* freshness computed on read
* in-memory store by default
* optional SQLite persistence
* runtime configuration
* CLI inspection flags
* API contract documentation
* example payloads
* PowerShell smoke-test scripts
* Dockerfile and `.dockerignore` for local container builds

## Current API surface

Current local endpoints:

* `GET /healthz`
* `POST /v1/stripe/subscription-state`
* `POST /v1/stripe/events`
* `GET /v1/state/stripe/subscription?customer_id=...`

`POST /v1/stripe/events` is unsigned demo ingestion only. It is not production Stripe webhook handling.

## Storage modes

The default store is memory.

SQLite is opt-in with:

```bash
SIGNALRELAY_STORE=sqlite
```

The SQLite path is configured with:

```bash
SIGNALRELAY_DB_PATH=signalrelay.db
```

Freshness is computed on read from envelope timestamps. It is not trusted from storage as authoritative state.

## Runtime configuration

Current runtime configuration:

* `SIGNALRELAY_ADDR`
* `SIGNALRELAY_STORE`
* `SIGNALRELAY_DB_PATH`
* `SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS`

## Boundaries

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
* make stale state safe
* replace Stripe as the source of truth
* integrate with StateMirror yet

SignalRelay only reports observed state, provenance, payload hash, and freshness metadata.

## Next possible milestones

Possible next milestones, not current capabilities:

* real Stripe signature verification
* Docker build/run validation
* release binaries
* stronger operational docs
* StateMirror example integration
* versioned prototype release
