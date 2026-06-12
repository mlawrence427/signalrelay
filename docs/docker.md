# Docker validation

SignalRelay includes a Dockerfile for local container usage of the prototype.

The Dockerfile is intended for local prototype validation. It does not make SignalRelay production-ready.

## Build

```bash
docker build -t signalrelay:local .
```

## Validate

Run the local validation script:

```powershell
.\scripts\validate-docker.ps1
```

The script builds `signalrelay:local`, starts a memory-mode container, binds the published port to `127.0.0.1`, calls `http://127.0.0.1:8080/healthz`, and stops the container at the end.

If Docker is unavailable, the script fails with a clear message.

## Run memory mode

```bash
docker run --rm -p 127.0.0.1:8080:8080 signalrelay:local
```

Then verify:

```bash
curl http://127.0.0.1:8080/healthz
```

Expected:

```json
{ "ok": true }
```

## Run SQLite mode

Use a mounted Docker volume:

```bash
docker run --rm -p 127.0.0.1:8080:8080 \
  -e SIGNALRELAY_STORE=sqlite \
  -e SIGNALRELAY_DB_PATH=/data/signalrelay.db \
  -v signalrelay-data:/data \
  signalrelay:local
```

PowerShell version:

```powershell
docker run --rm -p 127.0.0.1:8080:8080 `
  -e SIGNALRELAY_STORE=sqlite `
  -e SIGNALRELAY_DB_PATH=/data/signalrelay.db `
  -v signalrelay-data:/data `
  signalrelay:local
```

## Smoke test

The existing PowerShell smoke scripts can be run from the host after the container is running:

```powershell
.\scripts\smoke-memory.ps1
.\scripts\smoke-sqlite.ps1
```

## Boundaries

* The demo Stripe ingestion endpoint remains unsigned.
* Do not expose the demo ingestion endpoint publicly.
* The verified webhook endpoint requires Stripe-Signature verification.
* Container support is for local prototype usage only.
* Docker support does not make SignalRelay production-ready.
