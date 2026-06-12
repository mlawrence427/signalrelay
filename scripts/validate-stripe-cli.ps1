$ErrorActionPreference = "Stop"

if (-not (Get-Command stripe -ErrorAction SilentlyContinue)) {
    Write-Error "Stripe CLI is not available on PATH. Install Stripe CLI and log in before running live forwarding validation."
    exit 1
}

Write-Host "SignalRelay Stripe CLI live forwarding validation"
Write-Host ""
Write-Host "Terminal 1: start Stripe CLI forwarding:"
Write-Host "stripe listen --forward-to localhost:8080/v1/stripe/webhook"
Write-Host ""
Write-Host "Stripe CLI prints a webhook signing secret. It usually begins with whsec_."
$secret = Read-Host "Paste the webhook signing secret"
if ([string]::IsNullOrWhiteSpace($secret)) {
    Write-Error "Webhook signing secret is required."
    exit 1
}

Write-Host ""
Write-Host "Terminal 2: start SignalRelay with that secret:"
Write-Host "`$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET=`"$secret`""
Write-Host "go run ./cmd/signalrelay"
Write-Host ""
$confirm = Read-Host "Is SignalRelay running with that secret? Type yes to continue"
if ($confirm -ne "yes") {
    Write-Error "SignalRelay must be running before triggering the Stripe event."
    exit 1
}

Write-Host ""
Write-Host "Triggering customer.subscription.updated with Stripe CLI"
stripe trigger customer.subscription.updated
if ($LASTEXITCODE -ne 0) {
    Write-Error "stripe trigger customer.subscription.updated failed."
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "Query local state after identifying the generated customer id:"
Write-Host 'Invoke-RestMethod "http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_..."'
Write-Host ""
Write-Host "Inspect the forwarded event payload or SignalRelay logs to determine the customer id."
Write-Host "Stripe CLI forwarding validation commands completed."
