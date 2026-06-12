# Prototype Release Checklist

Purpose: Document what must be true before tagging `v0.0.1-prototype`.

## Release meaning

`v0.0.1-prototype` would mean:

* local prototype only
* signature-verified Stripe webhook ingestion exists
* memory and SQLite storage exist
* local signed webhook smoke test passes
* docs explain boundaries and non-production status
* not production-ready
* not a commercial SimpleStates artifact yet

## Required validation before tag

* [ ] `go test ./...` passes
* [ ] `gofmt` produces no changes
* [ ] local signed webhook smoke test passes
* [ ] Docker build/run validation completed or explicitly noted as not completed
* [ ] Stripe CLI live forwarding validation completed or explicitly noted as not completed
* [ ] run `.\scripts\validate-stripe-cli.ps1` when Stripe CLI is available
* [ ] record whether Stripe CLI live forwarding passed or remains unvalidated
* [ ] README prototype warnings reviewed
* [ ] CHANGELOG reviewed
* [ ] operations/security checklist reviewed

## Suggested release artifacts

* source archive from GitHub release
* optional compiled binaries:
  * Windows amd64
  * Linux amd64
  * macOS arm64
* optional checksum file
* optional Docker image later, but not required for prototype release

## Manual tag commands

Document commands only. Do not run them.

```bash
git status
git tag v0.0.1-prototype
git push origin v0.0.1-prototype
```

## GitHub release notes draft

### v0.0.1-prototype

SignalRelay `v0.0.1-prototype` is an early local prototype for stale-aware Stripe subscription state envelopes.

Implemented capabilities:

* local HTTP service
* manual state-envelope POST
* unsigned demo Stripe-shaped event ingestion
* signature-verified Stripe webhook ingestion
* duplicate event protection by `source_event_id`
* local read API by `customer_id`
* dynamic freshness computation
* memory store by default
* optional SQLite persistence
* runtime configuration
* CLI inspection flags
* local examples and smoke-test scripts
* local container build files
* API, operations/security, and prototype checkpoint documentation

This release is not production-ready.

SignalRelay does not decide access, enforce policy, orchestrate workflows, retry application actions, call Stripe APIs, replace Stripe as source of truth, integrate with StateMirror, or make stale state safe.

Validation status placeholders:

* `go test ./...`: TODO
* `gofmt` clean: TODO
* local signed webhook smoke test: TODO
* Docker build/run validation: TODO
* Stripe CLI live forwarding validation: TODO
