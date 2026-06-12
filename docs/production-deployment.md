# Production Deployment Boundary Guide

SignalRelay is still an early prototype. This guide explains what would need to be true before production deployment. It is not a production deployment guarantee.

## Current status

* signature-verified Stripe webhook endpoint exists
* local signed webhook smoke test passes
* memory and SQLite modes exist
* Dockerfile exists
* Docker validation is documented and scripted, but must be run in the local environment before release
* Stripe CLI live forwarding is documented and scripted, but must be manually validated before release
* no official release/tag exists yet

## Deployment shape

A minimal safe shape would require:

* expose only `POST /v1/stripe/webhook` publicly
* keep `POST /v1/stripe/events` private/local only
* terminate HTTPS with trusted infrastructure
* configure `SIGNALRELAY_STRIPE_WEBHOOK_SECRET` securely
* run with durable storage if state must survive restarts
* query SignalRelay from the application over private/internal networking

## Required environment variables

* `SIGNALRELAY_ADDR`
* `SIGNALRELAY_STORE`
* `SIGNALRELAY_DB_PATH`
* `SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS`
* `SIGNALRELAY_STRIPE_WEBHOOK_SECRET`
* `SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS`

Secrets must not be committed.

Secrets must not be logged.

Local test secrets must not be reused.

## Storage choice

Memory mode is restart-volatile.

SQLite mode persists locally.

The SQLite file location must be durable if persistence is expected.

Backups and file permissions are operator-owned.

## Network boundary

Public ingress should be limited to the verified webhook route.

The state query route should be internal/private.

The unsigned demo ingestion route should not be exposed publicly.

## Application boundary

Freshness is evidence. Risk tolerance is application logic.

SignalRelay reports observed Stripe subscription state.

The application decides access.

SignalRelay does not return allowed or denied.

SignalRelay does not enforce policy.

SignalRelay does not replace Stripe.

SignalRelay does not make stale state safe.

## Before production checklist

* [ ] Docker build/run validated
* [ ] Stripe CLI live forwarding validated
* [ ] real Stripe webhook endpoint tested
* [ ] secret rotation plan documented
* [ ] storage backup plan documented
* [ ] logging/retention plan documented
* [ ] restart/rollback process documented
* [ ] health check behavior reviewed
* [ ] operations/security checklist reviewed
* [ ] release/tag created
* [ ] versioned binaries or deployment artifact created
