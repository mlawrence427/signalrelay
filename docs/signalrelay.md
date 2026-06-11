# SignalRelay

SignalRelay is a self-hosted, stale-aware local read replica for externally sourced state facts.

It exists to make externally owned state locally available before an application makes a decision. It does not decide what the application should do with that state.

> SignalRelay does not make stale state safe. It makes stale state visible before the application decides.

> Freshness is evidence. Risk tolerance is application logic.

SignalRelay is pre-decision read infrastructure. StateMirror is post-decision evidence infrastructure.

## What It Is

SignalRelay provides local availability for externally sourced state facts.

The first useful version should be narrow: Stripe subscription state only. Stripe owns the subscription system. SignalRelay observes Stripe events, stores the observed state locally, and exposes that state to the application over localhost.

The application still owns the access decision.

The application should evaluate freshness, state, user context, product rules, and risk tolerance in its own code. For example, `evaluatePremiumAccess` belongs in the application, not in SignalRelay.

## What It Is Not

SignalRelay is not a service mesh, feature flag system, policy engine, workflow engine, webhook platform, generic cache, or authorization layer.

It should not become a place where product rules hide. It should not produce higher-level business conclusions from source facts. It should preserve what was observed, when it was observed, where it came from, and enough metadata for the application to judge whether it is usable.

## MVP Flow

The MVP flow should be:

1. Stripe webhook
2. SignalRelay ingest
3. Local state store
4. App queries localhost
5. App evaluates freshness and state
6. App records the decision in StateMirror
7. App grants or denies access

The first interface should be synchronous localhost pull, not push.

The application asks SignalRelay for the current local view of a source fact. SignalRelay returns an envelope. The application evaluates the envelope. SignalRelay does not call back into the application, retry application actions, or orchestrate workflows.

## State Envelope

A SignalRelay response should return observed state with provenance and freshness metadata:

```ts
type SignalRelayEnvelope<TPayload> = {
  source: "stripe";
  subject: string;
  state_type: "subscription";
  observed_at: string;
  stale_after: string;
  freshness: "fresh" | "stale" | "unknown";
  source_event_id: string;
  source_object_id: string;
  payload_hash: string;
  payload: TPayload;
};
```

The response must never return `allowed` or `denied`.

Access results are decisions, not source facts. SignalRelay should expose evidence for a decision, not the decision itself.

## Example

```ts
type StripeSubscriptionPayload = {
  customer_id: string;
  subscription_id: string;
  status: "active" | "trialing" | "past_due" | "canceled" | "incomplete";
  current_period_end: string;
};

async function evaluatePremiumAccess(userId: string) {
  const response = await fetch(
    `http://localhost:7841/v1/state/stripe/subscription/${userId}`,
  );

  const envelope =
    (await response.json()) as SignalRelayEnvelope<StripeSubscriptionPayload>;

  const decision =
    envelope.freshness === "fresh" && envelope.payload.status === "active"
      ? "allowed"
      : "denied";

  await stateMirror.record({
    decision_type: "premium_access",
    subject: userId,
    decision,
    evidence: envelope,
  });

  return decision;
}
```

This example keeps the boundary clear. SignalRelay returns the Stripe subscription envelope. The application owns `evaluatePremiumAccess`. StateMirror records the exact state envelope and decision.

## StateMirror Snapshot

Stripe can report an active subscription while the local observation is stale. In that case, the application may deny premium access because its risk tolerance requires fresh evidence.

```json
{
  "decision_type": "premium_access",
  "subject": "user_123",
  "decision": "denied",
  "reason": "stripe_subscription_state_stale",
  "decided_at": "2026-06-11T18:42:00Z",
  "evidence": {
    "source": "stripe",
    "subject": "cus_abc123",
    "state_type": "subscription",
    "observed_at": "2026-06-11T17:10:00Z",
    "stale_after": "2026-06-11T17:40:00Z",
    "freshness": "stale",
    "source_event_id": "evt_123",
    "source_object_id": "sub_123",
    "payload_hash": "sha256:0f2b7f4a9c1e5d8a6b3c2d1e9f00112233445566778899aabbccddeeff001122",
    "payload": {
      "customer_id": "cus_abc123",
      "subscription_id": "sub_123",
      "status": "active",
      "current_period_end": "2026-07-11T00:00:00Z"
    }
  }
}
```

## Boundary Rules

SignalRelay does not evaluate authorization.

SignalRelay does not decide access.

SignalRelay does not transform meaning.

SignalRelay does not retry app actions.

SignalRelay does not orchestrate workflows.

SignalRelay does not guarantee correctness.

SignalRelay does not guarantee freshness.

SignalRelay does not make stale state safe.

SignalRelay does not replace Stripe.

SignalRelay does not expose public traffic.

SignalRelay does not become a policy framework.

## Roadmap Fit

SignalRelay should not replace the current StateMirror roadmap. Validate it first through StateMirror examples using simulated freshness envelopes.
