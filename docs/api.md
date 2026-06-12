# SignalRelay Local API

This document describes the current local HTTP API for the SignalRelay prototype.

The API is limited to storing and reading stale-aware Stripe subscription state envelopes in memory. It does not verify real Stripe webhook signatures, persist state, authenticate requests, or evaluate access.

## GET /healthz

Returns service health.

Response:

```json
{ "ok": true }
```

## POST /v1/stripe/subscription-state

Stores the last observed Stripe subscription state envelope in memory.

Request fields:

* `source`
* `subject`
* `state_type`
* `observed_at`
* `stale_after`
* `source_event_id`
* `source_object_id`
* `payload`

`subject` is used as the local lookup key. For the current Stripe subscription path, callers should use the Stripe customer id as the subject.

`payload_hash` is computed by SignalRelay from the raw payload JSON. Clients do not provide it.

`freshness` is not trusted from input. It is computed on read from `stale_after`.

Example request:

```json
{
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
}
```

## GET /v1/state/stripe/subscription?customer_id=cus_123

Returns the last observed envelope for `customer_id`.

Response fields:

* `source`
* `subject`
* `state_type`
* `observed_at`
* `stale_after`
* `freshness`
* `source_event_id`
* `source_object_id`
* `payload_hash`
* `payload`

Supported `freshness` values:

* `fresh`
* `stale`

`freshness` is computed at read time:

* `fresh` when the current time is before `stale_after`
* `stale` when the current time is after `stale_after`

Missing customer response:

```json
{ "error": "subscription_state_missing" }
```

## Validation Errors

Invalid POST bodies return HTTP 400 with a JSON error body.

Current error strings:

* `invalid_json`
* `subject_required`
* `state_type_required`
* `observed_at_required`
* `observed_at_invalid`
* `stale_after_required`
* `stale_after_invalid`
* `payload_required`

The current implementation also validates additional envelope fields and may return:

* `source_required`
* `source_event_id_required`
* `source_object_id_required`

## Boundary Note

SignalRelay responses never return `allowed` or `denied`.

SignalRelay does not:

* evaluate authorization
* decide access
* enforce policy
* orchestrate workflows
* make stale state safe

SignalRelay only returns observed state, provenance, and freshness metadata.
