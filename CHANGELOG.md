# Changelog

## Unreleased

### Added

* local HTTP service
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion
* stricter demo Stripe event validation
* duplicate demo event protection by `source_event_id`
* local read API by `customer_id`
* dynamic freshness computation
* in-memory store by default
* optional SQLite persistence
* runtime configuration
* CLI inspection flags
* examples and PowerShell smoke tests
* Dockerfile and `.dockerignore` for local container builds
* API contract documentation
* prototype checkpoint documentation

### Not implemented

* production webhook handling
* Stripe signature verification
* authorization evaluation
* access decisions
* policy enforcement
* workflow orchestration
* Stripe API calls
* StateMirror integration
* production release artifacts
