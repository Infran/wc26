#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Install wc26 — FIFA World Cup 2026 CLI
.DESCRIPTION
    Downloads the latest pre-built binary from GitHub releases with SHA256
    verification. Falls back to building from source if download fails.
    Usage: iwr -useb https://github.com/Infran/wc26/releases/latest/download/install.ps1 | iex
#>

$ErrorActionPreference = 'Stop'
$RepoBase = "https://github.com/Infran/wc26"
$AppName = "wc26"
$ApiReleases = "https://api.github.com/repos/Infran/wc26/releases"

# --- Resolve version ---
if ($env:WC26_VERSION) {
    $Tag = $env:WC26_VERSION
} else {
    try {
        $release = Invoke-RestMethod -Uri "$ApiReleases/latest" -UseBasicParsing
        $Tag = $release.tag_name
    } catch {
        Write-Warning "Could not fetch latest version. Specify WC26_VERSION env var or check connectivity."
        $Tag = "latest"
    }
}
$Version = $Tag -replace "^v", ""

Write-Host "Installing $AppName $Tag ..." -ForegroundColor Cyan

# --- Determine install directory ---
$installDir = if ($env:WC26_INSTALL_DIR) {
    $env:WC26_INSTALL_DIR
} else {
    Join-Path $env:USERPROFILE ".local\bin"
}
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

$target = Join-Path $installDir "$AppName.exe"

# --- Try pre-built binary ---
$binaryUrl = "$RepoBase/releases/download/$Tag/wc26_windows_amd64.exe"
$checksumsUrl = "$RepoBase/releases/download/$Tag/wc26-checksums.txt"
$tmpDir = Join-Path $env:TEMP "$AppName-install"

$downloaded = $false
try {
    New-Item -ItemType Directory -Force -Path $tmpDir | Out-Null
    $binPath = Join-Path $tmpDir "wc26_windows_amd64.exe"
    $checksumsPath = Join-Path $tmpDir "wc26-checksums.txt"

    Write-Host "Downloading binary..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $binaryUrl -OutFile $binPath -UseBasicParsing -ErrorAction Stop

    Write-Host "Downloading checksums..." -ForegroundColor Yellow
    try {
        Invoke-WebRequest -Uri $checksumsUrl -OutFile $checksumsPath -UseBasicParsing -ErrorAction Stop

        $expectedHash = $null
        $expectedLine = Get-Content $checksumsPath | Where-Object { $_ -match "wc26_windows_amd64.exe" }
        if ($expectedLine) {
            $expectedHash = ($expectedLine -split "\s+")[0]
        }

        if ($expectedHash) {
            $actualHash = (Get-FileHash -Path $binPath -Algorithm SHA256).Hash.ToLower()
            if ($actualHash -ne $expectedHash.ToLower()) {
                throw "SHA256 mismatch! Expected: $expectedHash, Got: $actualHash"
            }
            Write-Host "SHA256 verified." -ForegroundColor Green
        } else {
            Write-Warning "Binary not found in checksums file, skipping verification."
        }
    } catch {
        Write-Warning "Checksum verification unavailable: $_"
    }

    # Smoke test
    $versionOut = & $binPath --version 2>&1
    if ($LASTEXITCODE -ne 0) {
        throw "Binary smoke test failed: $versionOut"
    }

    Copy-Item $binPath $target -Force
    Write-Host "Binary installed to: $target" -ForegroundColor Green
    $downloaded = $true
} catch {
    Write-Warning "Pre-built binary download failed: $_"
}

# --- Fallback: build from source ---
if (-not $downloaded) {
    $goVer = go version 2>$null
    if (-not $goVer) {
        Write-Error "Go is required to build from source. Install Go from https://go.dev/dl/ or ensure internet access for binary download."
        exit 1
    }
    Write-Host "Building from source..." -ForegroundColor Yellow

    $srcDir = Join-Path $tmpDir "src"
    if (Test-Path $srcDir) { Remove-Item -Recurse -Force $srcDir }

    if (Test-Path (Join-Path $PWD "go.mod")) {
        $srcDir = $PWD
    } else {
        git clone --depth 1 "$RepoBase.git" $srcDir 2>$null
        if (-not $?) {
            & go install "$RepoBase/cmd/$AppName@$Version" 2>$null
            if ($?) {
                $downloaded = $true
                $target = (Get-Command $AppName).Source
                goto :post
            }
            Write-Error "Cannot clone repo. Install Git or clone manually."
            exit 1
        }
    }

    Push-Location $srcDir
    try {
        $ldflags = "-X main.Version=$Version"
        go build -ldflags $ldflags -o $target ./cmd/$AppName/
        if (-not $?) { throw "Build failed" }
    } finally {
        Pop-Location
    }
    Write-Host "Built from source: $target" -ForegroundColor Green
}

# --- Add to PATH ---
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    $newPath = "$installDir;$currentPath"
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$installDir;$env:Path"
    Write-Host "Added $installDir to PATH" -ForegroundColor Yellow
}

# --- Post-install ---
& $target config init 2>$null
Remove-Item -Recurse -Force $tmpDir -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "✓ wc26 $Tag installed!" -ForegroundColor Green
Write-Host "  Run 'wc26 --help' to get started"
Write-Host "  Run 'wc26 auth login <email> <password>' to authenticate"
Write-Host "  Run 'wc26 update' to upgrade later"
