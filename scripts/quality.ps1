param()

$ErrorActionPreference = "Stop"

function Invoke-Checked {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Command,

        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Arguments
    )

    & $Command @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code ${LASTEXITCODE}: $Command $($Arguments -join ' ')"
    }
}

if (-not $env:GOCACHE) {
    $env:GOCACHE = Join-Path (Get-Location) ".gocache"
}
if (-not $env:GOMODCACHE) {
    $env:GOMODCACHE = Join-Path (Get-Location) ".gomodcache"
}

$unformatted = gofmt -l .\cmd .\pkg
if ($LASTEXITCODE -ne 0) {
    throw "gofmt failed"
}

if ($unformatted) {
    $unformatted
    throw "Go files need formatting. Run gofmt before committing."
}

Invoke-Checked go vet ./...
Invoke-Checked go test ./...
