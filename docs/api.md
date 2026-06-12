# SignalRelay Local API

This document describes the current local HTTP API for the SignalRelay prototype.

The API is limited to storing and reading stale-aware Stripe subscription state envelopes. The default local prototype store is in memory. Optional SQLite persistence can be enabled explicitly with `SIGNALRELAY_STORE=sqlite`.

The API includes a separate signature-verified Stripe webhook endpoint. It does not authenticate general requests or evaluate access.

## Local Storage

The default local prototype store is in memory.

Runtime configuration:

```bash
SIGNALRELAY_ADDR=:8080
SIGNALRELAY_STORE=memory
SIGNALRELAY_DB_PATH=signalrelay.db
SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS=300
SIGNALRELAY_STRIPE_WEBHOOK_SECRET=
SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS=300
```

To use optional SQLite persistence:

```bash
SIGNALRELAY_ADDR=:8080 SIGNALRELAY_STORE=sqlite SIGNALRELAY_DB_PATH=signalrelay.db go run ./cmd/signalrelay
```

On Windows PowerShell:

```powershell
$env:SIGNALRELAY_ADDR=":8080"
$env:SIGNALRELAY_STORE="sqlite"
$env:SIGNALRELAY_DB_PATH="signalrelay.db"
go run ./cmd/signalrelay
```

## GET /healthz

Returns service health.

Response:

```json
{ "ok": true }
```

## POST /v1/stripe/subscription-state

Stores the last observed Stripe subscription state envelope in the configured local store.

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

Static example payloads are available in `examples/`. Refresh their `observed_at` and `stale_after` values if you need a response to report `fresh`.

## POST /v1/stripe/events

Accepts an unsigned demo Stripe-shaped event payload and converts supported subscription events into the existing SignalRelay state envelope.

This is demo ingestion only. This endpoint does not verify `Stripe-Signature`, call the Stripe API, handle secrets, or claim production webhook behavior.

Supported event types:

* `customer.subscription.created`
* `customer.subscription.updated`
* `customer.subscription.deleted`

Unsupported event types return HTTP 400:

```json
{ "error": "unsupported_stripe_event_type" }
```

Simplified event shape:

```json
{
  "id": "evt_123",
  "type": "customer.subscription.updated",
  "created": 1760000000,
  "data": {
    "object": {
      "id": "sub_123",
      "object": "subscription",
      "customer": "cus_123",
      "status": "active",
      "current_period_end": 1762600000,
      "cancel_at_period_end": false
    }
  }
}
```

Required simplified event fields:

* `id`
* `type`
* `created`
* `data.object`
* `data.object.object`
* `data.object.id`
* `data.object.customer`
* `data.object.status`

`created` must be a positive Unix timestamp in seconds.

`data.object.object` must be `subscription`.

Mapping:

* `source` is `stripe`
* `subject` is `data.object.customer`
* `state_type` is `stripe.subscription`
* `observed_at` is `created` converted from Unix seconds to RFC3339 time
* `stale_after` is `observed_at` plus `SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS`
* `source_event_id` is `id`
* `source_object_id` is `data.object.id`
* `payload` is the raw `data.object` JSON
* `payload_hash` is computed from the raw `data.object` JSON

`SIGNALRELAY_STRIPE_STALE_AFTER_SECONDS` defaults to `300` and must be a positive integer. Invalid values fail startup.

The response is the stored envelope response. It never includes `allowed` or `denied`.

Duplicate event handling:

SignalRelay records ingested Stripe event ids by `source_event_id`. Reposting the same event id does not rewrite the stored envelope.

Duplicate response:

```json
{
  "duplicate": true,
  "source_event_id": "evt_123",
  "subject": "cus_123"
}
```

This is duplicate event protection for the local observed-state store. It is not workflow retry orchestration.

An example event payload is available at `examples/stripe-event-subscription-updated.json`.

## POST /v1/stripe/webhook

Accepts a Stripe-shaped event payload only after verifying the `Stripe-Signature` header with `SIGNALRELAY_STRIPE_WEBHOOK_SECRET`.

The verified endpoint uses the raw request body for verification before parsing JSON. The signature base string is:

```text
timestamp + "." + raw_body
```

Verification uses HMAC SHA-256 and the `v1` signature from the `Stripe-Signature` header.

`SIGNALRELAY_STRIPE_WEBHOOK_SECRET` is required when calling this endpoint. The secret is not required at startup.

`SIGNALRELAY_STRIPE_SIGNATURE_TOLERANCE_SECONDS` defaults to `300` and must be a positive integer. Invalid values fail startup.

After verification, this endpoint reuses the same subscription event ingestion path as `POST /v1/stripe/events`:

* same supported event types
* same event-shape validation
* same envelope mapping
* same duplicate event handling by `source_event_id`
* same memory and SQLite behavior
* same freshness computation

The response is the stored envelope response, or the duplicate event response when the event id has already been ingested. Responses never include `allowed` or `denied`.

Signature verification errors return HTTP 400 with JSON error bodies:

* `stripe_webhook_secret_not_configured`
* `stripe_signature_header_required`
* `stripe_signature_timestamp_invalid`
* `stripe_signature_timestamp_outside_tolerance`
* `stripe_signature_v1_required`
* `stripe_signature_mismatch`

This endpoint verifies the local webhook signature for this prototype. It does not call Stripe APIs, enforce policy, decide access, or make SignalRelay production software.

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
