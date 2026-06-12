Write-Host "SignalRelay Stripe CLI local validation"
Write-Host ""
Write-Host "Terminal 1: start Stripe CLI forwarding"
Write-Host 'stripe listen --forward-to localhost:8080/v1/stripe/webhook'
Write-Host ""
Write-Host "Copy the printed webhook signing secret. It usually begins with whsec_."
Write-Host ""
Write-Host "Terminal 2: start SignalRelay with the copied secret"
Write-Host '$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_REPLACE_ME"'
Write-Host 'go run ./cmd/signalrelay'
Write-Host ""
Write-Host "Optional SQLite mode:"
Write-Host '$env:SIGNALRELAY_STORE="sqlite"'
Write-Host '$env:SIGNALRELAY_DB_PATH="signalrelay-dev.db"'
Write-Host '$env:SIGNALRELAY_STRIPE_WEBHOOK_SECRET="whsec_REPLACE_ME"'
Write-Host 'go run ./cmd/signalrelay'
Write-Host ""
Write-Host "Terminal 3: trigger a supported subscription event"
Write-Host 'stripe trigger customer.subscription.updated'
Write-Host ""
Write-Host "Query the stored state with the customer id from the generated event payload"
Write-Host 'Invoke-RestMethod "http://localhost:8080/v1/state/stripe/subscription?customer_id=cus_..."'
