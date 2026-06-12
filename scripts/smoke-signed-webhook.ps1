$ErrorActionPreference = "Stop"

$baseUrl = "http://localhost:8080"
$secret = "whsec_signalrelay_local_test"
$eventPath = Join-Path $PSScriptRoot "..\examples\stripe-event-subscription-updated.json"

Write-Host "Start SignalRelay in another terminal with:"
Write-Host '$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_signalrelay_local_test"; go run ./cmd/signalrelay'
Write-Host ""

$rawBody = Get-Content -Raw -LiteralPath $eventPath
$timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$signatureBase = "$timestamp.$rawBody"

$hmac = [System.Security.Cryptography.HMACSHA256]::new([System.Text.Encoding]::UTF8.GetBytes($secret))
try {
    $hashBytes = $hmac.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($signatureBase))
    $signature = ($hashBytes | ForEach-Object { $_.ToString("x2") }) -join ""
}
finally {
    $hmac.Dispose()
}

$headers = @{
    "Stripe-Signature" = "t=$timestamp,v1=$signature"
}

try {
    $webhookResponse = Invoke-RestMethod `
        -Method Post `
        -Uri "$baseUrl/v1/stripe/webhook" `
        -ContentType "application/json" `
        -Headers $headers `
        -Body $rawBody
}
catch {
    Write-Error "Webhook request failed: $($_.Exception.Message)"
    exit 1
}

Write-Host "Webhook response:"
$webhookJson = $webhookResponse | ConvertTo-Json -Depth 20
Write-Host $webhookJson
Write-Host ""

try {
    $stateResponse = Invoke-RestMethod "$baseUrl/v1/state/stripe/subscription?customer_id=cus_123"
}
catch {
    Write-Error "State query failed: $($_.Exception.Message)"
    exit 1
}

Write-Host "State response:"
$stateJson = $stateResponse | ConvertTo-Json -Depth 20
Write-Host $stateJson

$combined = "$webhookJson`n$stateJson"
if ($combined -match "allowed") {
    Write-Error "Response included allowed"
    exit 1
}
if ($combined -match "denied") {
    Write-Error "Response included denied"
    exit 1
}
if ($stateJson -notmatch "source_event_id") {
    Write-Error "State response did not include source_event_id"
    exit 1
}
if ($stateJson -notmatch "payload_hash") {
    Write-Error "State response did not include payload_hash"
    exit 1
}
if ($stateJson -notmatch "freshness") {
    Write-Error "State response did not include freshness"
    exit 1
}

Write-Host ""
Write-Host "Signed webhook smoke test passed."
