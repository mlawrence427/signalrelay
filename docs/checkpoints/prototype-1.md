# SignalRelay Prototype Checkpoint 1

SignalRelay is now an early local prototype for stale-aware Stripe subscription state envelopes.

It can accept demo Stripe-shaped subscription events, map them into local state envelopes, protect against duplicate demo events, store envelopes in memory or SQLite, and expose the last observed state through a localhost read API.

It is not production webhook handling.

## Current capabilities

SignalRelay currently includes:

* local HTTP service
* configurable listen address with `SIGNALRELAY_ADDR`
* memory store by default
* optional SQLite store with `SIGNALRELAY_STORE=sqlite`
* configurable SQLite path with `SIGNALRELAY_DB_PATH`
* manual envelope POST endpoint
* unsigned demo Stripe-shaped event ingestion endpoint
* supported demo event types:
  * `customer.subscription.created`
  * `customer.subscription.updated`
  * `customer.subscription.deleted`
* demo event idempotency by `source_event_id`
* `payload_hash` computation
* freshness computed on read from `stale_after`
* API validation
* JSON responses
* API contract docs
* example payloads
* PowerShell smoke-test scripts

## Current API surface

Current local endpoints:

* `GET /healthz`
* `POST /v1/stripe/subscription-state`
* `POST /v1/stripe/events`
* `GET /v1/state/stripe/subscription?customer_id=...`

`POST /v1/stripe/events` is unsigned demo ingestion only. It is not production Stripe webhook handling.

## Boundary rules

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

## Freshness rule

Freshness is evidence. Risk tolerance is application logic.

SignalRelay computes freshness from the stored envelope timestamps. It does not decide whether stale state is acceptable.

## Prototype status

This checkpoint represents a runnable local prototype, not production software.

The prototype is useful for:

* testing the state envelope model
* demonstrating local state reads
* demonstrating stale-but-active subscription state
* demonstrating duplicate event protection
* validating the boundary between SignalRelay and application-owned decisions

## Next safe milestones

Possible next milestones, not yet implemented:

* real Stripe signature verification
* stronger event-shape validation
* CLI/version output
* release binaries
* StateMirror example integration
* Dockerfile
* operational hardening
