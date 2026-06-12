# Changelog

## Unreleased

### Added

* local HTTP service
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion
* signature-verified Stripe webhook endpoint for supported subscription events
* stricter demo Stripe event validation
* duplicate demo event protection by `source_event_id`
* local read API by `customer_id`
* dynamic freshness computation
* in-memory store by default
* optional SQLite persistence
* runtime configuration
* CLI inspection flags
* examples and PowerShell smoke tests
* local signed webhook smoke test
* prototype release checklist
* local release build helper script
* StateMirror integration example documentation
* StateMirror-style premium access snapshot example
* production deployment boundary guide
* Dockerfile and `.dockerignore` for local container builds
* Docker validation script
* operations and security hardening checklist
* Docker validation documentation
* Stripe CLI validation documentation
* API contract documentation
* prototype checkpoint documentation
* Prototype Checkpoint 3 documentation

### Not implemented

* production webhook handling
* authorization evaluation
* access decisions
* policy enforcement
* workflow orchestration
* Stripe API calls
* StateMirror integration
* production release artifacts
