# SignalRelay

SignalRelay is a future SimpleStates concept for stale-aware local read infrastructure for externally sourced state facts.

It is not an active production project yet.

Core boundary:

> SignalRelay does not make stale state safe. It makes stale state visible before the application decides.

> Freshness is evidence. Risk tolerance is application logic.

Current status: concept / research note.

Read the concept note:

* [SignalRelay concept note](docs/signalrelay.md)

SignalRelay is not a service mesh, feature flag system, policy engine, workflow engine, webhook platform, generic cache, or authorization layer.

It is intended to explore local availability for externally sourced state facts while keeping application decisions application-owned.
