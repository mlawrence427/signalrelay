$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$payloadPath = Join-Path $root "examples\stripe-subscription-active.json"
$baseURL = "http://localhost:8080"

Write-Host "SignalRelay SQLite smoke test"
Write-Host "Expecting SignalRelay to already be running with SQLite enabled:"
Write-Host '  $env:SIGNALRELAY_STORE="sqlite"'
Write-Host '  $env:SIGNALRELAY_DB_PATH="signalrelay-dev.db"'
Write-Host "Posting $payloadPath"
Write-Host "Refresh observed_at/stale_after in the example if you need freshness to be fresh."

$payload = Get-Content -LiteralPath $payloadPath -Raw
Invoke-RestMethod -Method Post -Uri "$baseURL/v1/stripe/subscription-state" -ContentType "application/json" -Body $payload | Out-Null

$response = Invoke-RestMethod -Uri "$baseURL/v1/state/stripe/subscription?customer_id=cus_123"
$json = $response | ConvertTo-Json -Depth 10
Write-Host $json

$requiredFields = @("source", "subject", "freshness", "payload_hash")
foreach ($field in $requiredFields) {
    if (-not $response.PSObject.Properties.Name.Contains($field)) {
        Write-Error "Missing required response field: $field"
        exit 1
    }
}

if ($json -match '"allowed"' -or $json -match '"denied"') {
    Write-Error "Response included an access decision"
    exit 1
}

Write-Host "Smoke test passed"
