# Local signed webhook smoke test

SignalRelay includes a local smoke test for the verified webhook endpoint:

```text
POST /v1/stripe/webhook
```

The smoke test generates a Stripe-style signature over `examples/stripe-event-subscription-updated.json` and sends it to a locally running SignalRelay process.

This does not replace Stripe CLI validation. It is a local prototype check for the signature verification path when Stripe CLI is not involved.

The script uses the same signature shape as Stripe webhooks:

```text
timestamp + "." + raw_body
```

It is for prototype validation only. SignalRelay still does not decide access, enforce policy, or make stale state safe.

## Usage

Terminal 1:

```powershell
$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_signalrelay_local_test"
go run ./cmd/signalrelay
```

Terminal 2:

```powershell
.\scripts\smoke-signed-webhook.ps1
```

The script:

* loads `examples/stripe-event-subscription-updated.json`
* signs the raw JSON body with `whsec_signalrelay_local_test`
* posts to `http://localhost:8080/v1/stripe/webhook`
* queries `http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_123`
* checks that responses do not include `allowed` or `denied`
* checks that the stored state includes `source_event_id`, `payload_hash`, and `freshness`

## Boundaries

The unsigned demo endpoint remains available at `POST /v1/stripe/events`.

The Stripe CLI validation flow remains documented separately in `docs/stripe-cli.md`.
