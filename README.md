[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/Infran/wc26?color=brightgreen)](https://github.com/Infran/wc26/releases/latest)

# wc26 — FIFA World Cup 2026 CLI

A fast Go CLI for the [World Cup 2026 API](https://worldcup26.ir). Query teams, groups, matches, stadiums, and live scores from your terminal. Works with the hosted API or any self-hosted instance.

## Install

### One-liner (no prerequisites)

**Windows (PowerShell 5+)**:
```powershell
iwr -useb https://github.com/Infran/wc26/releases/latest/download/install.ps1 | iex
```

**macOS / Linux**:
```bash
curl -fsSL https://github.com/Infran/wc26/releases/latest/download/install.sh | bash
```

Downloads a pre-built binary with SHA256 verification. No Go, no Git needed.

### From source

```bash
go install github.com/Infran/wc26/cmd/wc26@latest
```

### Manual

Download the binary for your platform from the [latest release](https://github.com/Infran/wc26/releases/latest) and place it in your PATH.

| Platform | Binary |
|----------|--------|
| Windows x64 | `wc26_windows_amd64.exe` |
| macOS Intel | `wc26_darwin_amd64` |
| macOS Apple Silicon | `wc26_darwin_arm64` |
| Linux x64 | `wc26_linux_amd64` |
| Linux ARM64 | `wc26_linux_arm64` |

## Update

```bash
# Check for update
wc26 update --check

# Download and install latest (atomic .old swap)
wc26 update
```

The update command fetches the latest GitHub release, downloads the binary, swaps it atomically (old binary → `wc26.old`, new binary → `wc26`, removes `.old` on success, rolls back on failure).

## Quick Start

```bash
# Authenticate (one time)
wc26 auth register "Your Name" you@example.com yourpassword

# Or login if already registered
wc26 auth login you@example.com yourpassword

# Show your JWT token (prompts for password to validate)
wc26 auth token

# Check server health
wc26 health

# List all teams
wc26 teams

# List teams in Group A
wc26 teams --group A

# Get a specific team
wc26 team Brazil
wc26 team 37

# List all groups with standings
wc26 groups

# Get single group
wc26 group A

# List all matches
wc26 matches

# Filter matches
wc26 matches --type group
wc26 matches --team Brazil
wc26 matches --matchday 1

# Get a specific match
wc26 match 1

# List all stadiums
wc26 stadiums

# Get a specific stadium
wc26 stadium 1
```

## Output Formats

All commands support `--output` / `-o`:

| Format | Use case |
|--------|----------|
| `table` (default) | Human-readable tables |
| `json` | Pipe to `jq`, scripts, CI/CD |
| `plain` | Simple key:value for grep/awk |

```bash
wc26 teams --output json | jq '.[] | {name: .name_en, code: .fifa_code}'
wc26 matches --output json > matches.json
```

## Custom API URL

Point to a self-hosted instance or override at runtime:

```bash
wc26 --api-url http://localhost:3050 health
wc26 config set api.base_url http://localhost:3050
```

## Shell Completions

```bash
# Bash
source <(wc26 completion bash)

# Zsh
source <(wc26 completion zsh)

# Fish
wc26 completion fish | source

# PowerShell
wc26 completion powershell | Out-String | Invoke-Expression
```

## Configuration

Config is stored at the platform-appropriate location:

- **Windows**: `%APPDATA%\wc26\config.yaml`
- **macOS/Linux**: `~/.config/wc26/config.yaml`

```bash
wc26 config show    # view config
wc26 config path    # config file location
wc26 config set api.base_url https://my-instance.example.com
```

## Authentication

The API requires a JWT token for most endpoints. Register or login once; the CLI saves your token.

```bash
wc26 auth register "Your Name" email@example.com password
wc26 auth login email@example.com password
wc26 auth status
wc26 auth token    # prompts for password, prints token to stdout
wc26 auth logout   # clears saved token
```

## Build from Source

```bash
git clone https://github.com/Infran/wc26.git
cd wc26
go build -ldflags "-X main.Version=$(git describe --tags 2>/dev/null || echo dev)" -o wc26 ./cmd/wc26/
```

## Release

```bash
./scripts/release.sh v0.3.0 "Release notes here"
```

Cross-compiles for all platforms, generates SHA256 checksums, includes install scripts, and publishes to GitHub releases.

## License

MIT License — see [LICENSE](LICENSE).

Copyright (c) 2026 Infran
