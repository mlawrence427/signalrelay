# SignalRelay Prototype Checkpoint 3

SignalRelay now supports signature-verified Stripe subscription webhook ingestion through a separate verified endpoint, while preserving the unsigned demo ingestion endpoint for local testing.

It remains an early prototype and does not decide access, enforce policy, or make stale state safe.

## Current capabilities

* local HTTP service
* memory store by default
* optional SQLite persistence
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion at `POST /v1/stripe/events`
* signature-verified Stripe webhook ingestion at `POST /v1/stripe/webhook`
* Stripe-Signature verification using raw request body
* duplicate event protection by `source_event_id`
* stricter subscription event validation
* freshness computed on read
* runtime configuration
* CLI inspection flags
* examples and smoke-test scripts
* Dockerfile and Docker validation docs
* API contract docs

## Current API surface

* `GET /healthz`
* `POST /v1/stripe/subscription-state`
* `POST /v1/stripe/events`
* `POST /v1/stripe/webhook`
* `GET /v1/state/stripe/subscription?customer_id=...`

`/v1/stripe/events` is unsigned demo ingestion.

`/v1/stripe/webhook` verifies Stripe-Signature before parsing JSON.

## Runtime configuration

* `SIGNALRELAY_ADDR`
* `SIGNALRELAY_STORE`
* `SIGNALRELAY_DB_PATH`
* `SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS`
* `SIGNALRELAY_STRIPE_WEBHOOK_SECRET`
* `SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS`

The webhook secret is not required at startup.

The verified webhook endpoint fails if the secret is not configured.

Secrets must not be logged.

## Boundary rules

SignalRelay does not:

* evaluate authorization
* decide access
* return allowed or denied
* enforce policy
* orchestrate workflows
* retry application actions
* call Stripe APIs
* replace Stripe as source of truth
* make stale state safe
* integrate with StateMirror yet

SignalRelay only reports observed state, provenance, payload hash, source event identity, and freshness metadata.

## Freshness rule

Freshness is evidence. Risk tolerance is application logic.

Freshness is computed from timestamps.

Stale state may still contain active billing status.

The application decides whether stale state is acceptable.

## Next possible milestones

Future possibilities, not current capabilities:

* real Stripe CLI/webhook end-to-end validation
* release binaries
* Docker build/run validation
* StateMirror example integration
* stronger operational docs
* v0.0.1 prototype tag
