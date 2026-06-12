# StateMirror Integration Example

SignalRelay provides observed external state with freshness and provenance metadata. The application remains responsible for deciding what to do. StateMirror can then record the decision evidence snapshot.

## Flow

1. Stripe sends a subscription webhook.
2. SignalRelay verifies the webhook signature.
3. SignalRelay stores the observed subscription state envelope.
4. The application queries SignalRelay.
5. The application evaluates access using its own policy.
6. The application records a StateMirror snapshot containing:
   * the SignalRelay state envelope
   * the application-owned decision
   * freshness metadata
   * provenance fields
   * payload hash
   * source event id

SignalRelay does not decide access.
StateMirror does not decide access.
The application owns the decision.

## TypeScript-style pseudocode

```typescript
type SignalRelayEnvelope = {
  source: string;
  subject: string;
  state_type: string;
  observed_at: string;
  stale_after: string;
  freshness: "fresh" | "stale";
  source_event_id: string;
  source_object_id: string;
  payload_hash: string;
  payload: {
    status?: string;
    [key: string]: unknown;
  };
};

async function checkPremiumAccess(customerId: string) {
  const response = await fetch(
    `http://localhost:8080/v1/state/stripe/subscription?customer_id=${encodeURIComponent(customerId)}`,
  );

  if (!response.ok) {
    return {
      allowed: false,
      reason: "subscription_state_missing",
      decided_by: "application",
    };
  }

  const signalrelay = (await response.json()) as SignalRelayEnvelope;

  const subscriptionIsActive = signalrelay.payload.status === "active";
  const stateIsFresh = signalrelay.freshness === "fresh";

  const decision = {
    allowed: subscriptionIsActive && stateIsFresh,
    reason: subscriptionIsActive
      ? stateIsFresh
        ? "stripe_subscription_active_and_fresh"
        : "stripe_subscription_active_but_stale"
      : "stripe_subscription_not_active",
    decided_by: "application",
  };

  const snapshot = {
    snapshot_type: "premium_access.check",
    subject: customerId,
    captured_at: new Date().toISOString(),
    decision,
    evidence: {
      signalrelay,
    },
    boundary: {
      signalrelay_role: "observed external state",
      application_role: "decision owner",
      statemirror_role: "decision evidence recorder",
    },
  };

  await recordStateMirrorSnapshot(snapshot);

  if (decision.allowed) {
    return grantPremiumAccess(customerId);
  }

  return denyPremiumAccess(customerId, decision.reason);
}
```

The SignalRelay response contains observed state and metadata. It does not contain an access decision.

The application evaluates freshness and subscription status, records a StateMirror-style evidence snapshot, and then acts based on its own decision.

See `examples/statemirror-premium-access-snapshot.json` for a concrete snapshot shape.
