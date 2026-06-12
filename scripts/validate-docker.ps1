$ErrorActionPreference = "Stop"

$imageName = "signalrelay:local"
$containerName = "signalrelay-local-test"
$healthUrl = "http://127.0.0.1:8080/healthz"

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Error "Docker is not available on PATH. Install or start Docker before running local container validation."
    exit 1
}

$containerStarted = $false

try {
    $existingContainer = docker ps -a --filter "name=^/$containerName$" --format "{{.Names}}"
    if ($existingContainer -eq $containerName) {
        docker rm -f $containerName | Out-Null
    }

    Write-Host "Building $imageName"
    docker build -t $imageName .

    Write-Host "Starting $containerName"
    docker run --rm -d --name $containerName -p 127.0.0.1:8080:8080 $imageName | Out-Null
    $containerStarted = $true

    $healthBody = $null
    for ($i = 0; $i -lt 30; $i++) {
        try {
            $healthBody = Invoke-WebRequest -UseBasicParsing -Uri $healthUrl
            if ($healthBody.Content -match '"ok"\s*:\s*true') {
                break
            }
        }
        catch {
            Start-Sleep -Milliseconds 500
        }
    }

    if ($null -eq $healthBody -or $healthBody.Content -notmatch '"ok"\s*:\s*true') {
        throw "Health check did not return ok=true from $healthUrl"
    }

    Write-Host "Docker validation passed."
}
finally {
    if ($containerStarted) {
        docker stop $containerName | Out-Null
    }
}
