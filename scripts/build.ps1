param(
    [string]$OutDir = "dist",
    [string]$Version = "dev"
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

New-Item -ItemType Directory -Path $OutDir -Force | Out-Null

if (-not $env:GOCACHE) {
    $env:GOCACHE = Join-Path (Get-Location) ".gocache"
}
if (-not $env:GOMODCACHE) {
    $env:GOMODCACHE = Join-Path (Get-Location) ".gomodcache"
}

$binary = "prmaven"
if ($IsWindows -or $env:OS -eq "Windows_NT") {
    $binary = "prmaven.exe"
}

$target = Join-Path $OutDir $binary
Invoke-Checked go build "-trimpath" "-ldflags=-s -w -X main.version=$Version" "-o" $target .\cmd\prmaven

Write-Output $target
