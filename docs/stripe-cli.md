# Stripe CLI validation

SignalRelay supports a signature-verified endpoint:

```text
POST /v1/stripe/webhook
```

This can be tested locally with Stripe CLI.

Important:

* This is local validation.
* SignalRelay is still an early prototype.
* Stripe signature verification exists, but this does not make the project production-ready.
* SignalRelay still does not decide access, enforce policy, or make stale state safe.

Stripe CLI can forward sandbox events to a local endpoint with `stripe listen --forward-to`. The listen command prints a webhook signing secret that can be used to validate local webhook signatures.

## 1. Start SignalRelay with a local webhook secret

Terminal 1:

```powershell
stripe listen --forward-to localhost:8080/v1/stripe/webhook
```

Copy the printed webhook signing secret. It usually begins with `whsec_`.

Terminal 2:

```powershell
$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_REPLACE_ME"
go run ./cmd/signalrelay
```

Optional SQLite mode:

```powershell
$env:SIGNALRELAY_STORE="sqlite"
$env:SIGNALRELAY_DB_PATH="signalrelay-dev.db"
$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_REPLACE_ME"
go run ./cmd/signalrelay
```

## 2. Trigger a subscription event

Terminal 3:

```powershell
stripe trigger customer.subscription.updated
```

Stripe CLI triggers can create related events. SignalRelay only accepts supported subscription event types:

* `customer.subscription.created`
* `customer.subscription.updated`
* `customer.subscription.deleted`

Unsupported event types are rejected.

## 3. Query local state

Use the customer id from the event payload Stripe generated:

```powershell
Invoke-RestMethod "http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_..."
```

SignalRelay returns the last observed state envelope with provenance, payload hash, source event identity, and freshness metadata.

## 4. Troubleshooting

`stripe_webhook_secret_not_configured`

The verified endpoint was called without `SIGNALRELAY_STRIPE_WEBHOOK_SECRET` set.

`stripe_signature_required`

The request did not include a usable Stripe signature header. For the current local API, the exact error string may be `stripe_signature_header_required`.

`stripe_signature_verification_failed`

The signature could not be verified. For the current local API, the exact error string may be `stripe_signature_mismatch`, `stripe_signature_timestamp_invalid`, `stripe_signature_timestamp_outside_tolerance`, or `stripe_signature_v1_required`.

`unsupported_stripe_event_type`

Stripe CLI forwarded an event type outside the supported subscription event set.

`subscription_state_missing`

No stored subscription state exists for that `customer_id`. Confirm the customer id from the generated event payload and retry the query.

## Boundaries

The unsigned demo endpoint remains available at:

```text
POST /v1/stripe/events
```

The verified webhook endpoint does not decide access, enforce policy, call Stripe APIs, or make stale state safe.
