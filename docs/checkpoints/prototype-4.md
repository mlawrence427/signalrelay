# SignalRelay Prototype Checkpoint 4

SignalRelay is now a prototype-ready local service for signature-verified Stripe subscription state ingestion and stale-aware local state reads. It remains an early prototype, not production software.

## What is implemented

* local HTTP service
* memory store
* optional SQLite persistence
* manual state-envelope ingestion
* unsigned demo event ingestion at `POST /v1/stripe/events`
* signature-verified webhook ingestion at `POST /v1/stripe/webhook`
* raw-body Stripe-Signature verification
* duplicate event protection
* stricter Stripe subscription event validation
* freshness metadata computed on read
* local signed webhook smoke test
* Dockerfile
* Docker validation script
* Stripe CLI validation helper
* operations/security checklist
* production deployment boundary guide
* release checklist and local build helper
* StateMirror integration example

## Current validation status

| Validation | Status |
| --- | --- |
| `go test ./...` | passing |
| `gofmt` | clean |
| local signed webhook smoke test | manually passed |
| local release binary build helper | manually passed |
| Docker validation | documented and scripted, not manually completed |
| Stripe CLI live forwarding | documented and scripted, not manually completed |

## Remaining before release/tag

* [ ] decide whether Docker validation can remain incomplete for `v0.0.1-prototype`
* [ ] decide whether Stripe CLI live forwarding can remain incomplete for `v0.0.1-prototype`
* [ ] review README prototype warnings
* [ ] review CHANGELOG
* [ ] review release checklist
* [ ] run `go test ./...`
* [ ] run signed webhook smoke test
* [ ] run build-release script
* [ ] confirm clean git status
* [ ] tag only if comfortable with incomplete Docker/Stripe CLI live validation being documented

## Boundary statement

Freshness is evidence. Risk tolerance is application logic.

SignalRelay reports observed external state, provenance, source event identity, payload hash, and freshness metadata. It does not decide access, enforce policy, orchestrate workflows, retry application actions, replace Stripe, or make stale state safe.

## Website readiness note

SignalRelay is ready to be described on the SimpleStates website as a prototype primitive or labs primitive, not as a commercial product or production-ready component.
