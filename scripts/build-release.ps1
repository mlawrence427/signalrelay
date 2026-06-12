$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$distDir = Join-Path $repoRoot "dist"

if (Test-Path -LiteralPath $distDir) {
    Remove-Item -LiteralPath $distDir -Recurse -Force
}
New-Item -ItemType Directory -Path $distDir | Out-Null

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "signalrelay-windows-amd64.exe" },
    @{ GOOS = "linux"; GOARCH = "amd64"; Output = "signalrelay-linux-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Output = "signalrelay-darwin-arm64" }
)

foreach ($target in $targets) {
    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    $outputPath = Join-Path $distDir $target.Output

    Write-Host "Building $outputPath"
    go build -o $outputPath ./cmd/signalrelay
}

Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Generated release prototype binaries:"
Get-ChildItem -LiteralPath $distDir | ForEach-Object {
    Write-Host $_.FullName
}
