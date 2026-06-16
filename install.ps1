#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Install wc26 — FIFA World Cup 2026 CLI
.DESCRIPTION
    Builds and installs the wc26 CLI tool from source.
    Requires Go 1.22+.
#>

$ErrorActionPreference = 'Stop'
$RepoBase = "https://github.com/Infran/wc26"
$AppName = "wc26"
$Version = if ($env:WC26_VERSION) { $env:WC26_VERSION } else { "latest" }

# --- Prerequisites ---
$goVer = go version 2>$null
if (-not $goVer) {
    Write-Host "Go is not installed. Installing Go..." -ForegroundColor Yellow
    $goUrl = "https://go.dev/dl/go1.24.0.windows-amd64.msi"
    $goInstaller = "$env:TEMP\go.msi"
    try {
        Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller -UseBasicParsing
        Start-Process msiexec -ArgumentList "/i `"$goInstaller`" /quiet /norestart" -Wait
        $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
        Write-Host "Go installed. Restart your terminal if detection fails." -ForegroundColor Green
    } catch {
        Write-Error "Failed to install Go. Install manually from https://go.dev/dl/"
        exit 1
    }
}

Write-Host "Found: $(go version)" -ForegroundColor Green

# --- Determine install directory ---
$installDir = if ($env:WC26_INSTALL_DIR) {
    $env:WC26_INSTALL_DIR
} else {
    Join-Path $env:USERPROFILE ".local\bin"
}

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

# --- Clone or use local ---
$localPath = Join-Path $PWD $AppName
if (Test-Path (Join-Path $localPath "go.mod")) {
    Write-Host "Using local source at $localPath" -ForegroundColor Cyan
    $srcDir = $localPath
} else {
    $srcDir = Join-Path $env:TEMP $AppName
    if (Test-Path $srcDir) { Remove-Item -Recurse -Force $srcDir }
    Write-Host "Cloning $RepoBase ..." -ForegroundColor Yellow
    git clone --depth 1 $RepoBase $srcDir 2>$null
    if (-not $?) {
        Write-Host "Git not available, trying go install directly..." -ForegroundColor Yellow
        & go install "$RepoBase/cmd/$AppName@$Version" 2>$null
        if ($?) {
            Write-Host "Installed via 'go install'" -ForegroundColor Green
            goto :post_install
        }
        Write-Error "Cannot clone or install. Install Git or clone manually."
        exit 1
    }
}

# --- Build ---
Write-Host "Building $AppName..." -ForegroundColor Yellow
$binary = Join-Path $srcDir "$AppName.exe"
Push-Location $srcDir
try {
    $ldflags = "-X main.Version=$Version"
    & go build -ldflags $ldflags -o $binary ./cmd/$AppName/
    if (-not $?) {
        Write-Error "Build failed"
        exit 1
    }
} finally {
    Pop-Location
}

# --- Install ---
$target = Join-Path $installDir "$AppName.exe"
Copy-Item $binary $target -Force
Write-Host "Installed: $target" -ForegroundColor Green

# --- Add to PATH ---
$currentPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    $newPath = "$installDir;$currentPath"
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$installDir;$env:Path"
    Write-Host "Added $installDir to PATH" -ForegroundColor Yellow
}

# --- Post-install ---
& $target config init
Write-Host ""
Write-Host "✓ wc26 installed successfully!" -ForegroundColor Green
Write-Host "  Run 'wc26 --help' to get started"
Write-Host "  Run 'wc26 auth login <email> <password>' to authenticate"
Write-Host "  Run 'wc26 health' to check API status"
