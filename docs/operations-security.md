# Operations and Security Hardening Checklist

SignalRelay is an early prototype. This checklist documents operational and security considerations before production use. It is not a production certification and does not imply the project is ready for production deployment.

## Current implemented safeguards

* separate unsigned demo endpoint and signature-verified webhook endpoint
* raw-body Stripe-Signature verification for `/v1/stripe/webhook`
* webhook secret supplied by environment variable
* webhook secret is not logged
* signature timestamp tolerance configuration
* duplicate event protection by `source_event_id`
* stricter Stripe subscription event validation
* no allowed/denied responses
* no authorization evaluation
* no policy enforcement
* local signed webhook smoke test
* memory store by default
* optional SQLite persistence

## Required before production use

### Webhook exposure

* [ ] verified endpoint only should be exposed
* [ ] unsigned demo endpoint must not be publicly exposed
* [ ] TLS/HTTPS termination must be handled by trusted infrastructure
* [ ] webhook route should be protected by Stripe signature verification
* [ ] raw body must be preserved before verification

### Secret handling

* [ ] webhook secret must come from a real secret manager or secure runtime config
* [ ] webhook secret must not be committed
* [ ] webhook secret must not be logged
* [ ] rotation procedure must exist
* [ ] local test secrets must never be reused in production

### Storage

* [ ] choose memory or SQLite deliberately
* [ ] memory mode loses state on restart
* [ ] SQLite file location must be durable if persistence is expected
* [ ] backups must be handled outside SignalRelay
* [ ] permissions for the SQLite file should be restricted

### Runtime/container

* [ ] Docker image should be locally validated
* [ ] container should run as non-root when practical
* [ ] expose only required ports
* [ ] mount persistent storage explicitly when using SQLite
* [ ] avoid baking secrets into images

### Operations

* [ ] define health check behavior
* [ ] define restart behavior
* [ ] define log retention policy
* [ ] define upgrade/rollback process
* [ ] define who owns webhook failures
* [ ] define what happens when Stripe delivery is delayed

### Application boundary

* [ ] application remains responsible for access decisions
* [ ] stale state must not be treated as safe by default
* [ ] freshness is evidence, not permission
* [ ] SignalRelay does not replace Stripe as source of truth
* [ ] SignalRelay does not integrate with StateMirror yet

## Production readiness status

| Item | Status |
| --- | --- |
| Signature-verified webhook ingestion | implemented |
| Local signed smoke test | implemented |
| Stripe CLI live forwarding validation | documented, not yet manually validated |
| Docker local build/run validation | documented and scripted, manual execution required |
| Production deployment guide | not written |
| Packaged release artifacts | not created |
| StateMirror integration example | not implemented |
| Operational ownership model | not defined |

## Boundary statement

Freshness is evidence. Risk tolerance is application logic.

SignalRelay reports observed external state, provenance, source event identity, payload hash, and freshness metadata. It does not decide access, enforce policy, orchestrate workflows, retry application actions, or make stale state safe.
