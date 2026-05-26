param(
    [switch]$Race,
    [switch]$Coverage,
    [decimal]$MinimumCoverage = 70
)

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

Invoke-Checked go test ./...

if ($Race) {
    $cgoEnabled = (& go env CGO_ENABLED).Trim()
    if ($LASTEXITCODE -ne 0) {
        throw "Could not inspect CGO_ENABLED"
    }

    if ($cgoEnabled -eq "1") {
        Invoke-Checked go test -race ./...
    } else {
        Write-Warning "Skipping race tests because CGO_ENABLED is $cgoEnabled. The GitHub Actions race job still runs on Linux."
    }
}

if ($Coverage) {
    Invoke-Checked go test "-coverprofile=coverage.out" ./...
    $coverageOutput = & go tool cover "-func=coverage.out"
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code ${LASTEXITCODE}: go tool cover -func=coverage.out"
    }

    $coverageOutput

    $totalLine = $coverageOutput | Where-Object { $_ -match '^total:\s+\(statements\)\s+([0-9.]+)%' } | Select-Object -Last 1
    if (-not $totalLine) {
        throw "Could not find total coverage line"
    }
    if ($totalLine -notmatch '^total:\s+\(statements\)\s+([0-9.]+)%') {
        throw "Could not parse total coverage line: $totalLine"
    }

    $totalCoverage = [decimal]::Parse($Matches[1], [System.Globalization.CultureInfo]::InvariantCulture)
    if ($totalCoverage -lt $MinimumCoverage) {
        throw "Total coverage $totalCoverage% is below required minimum $MinimumCoverage%"
    }
}
