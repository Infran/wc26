# wc26 — FIFA World Cup 2026 CLI

A Go CLI for the [World Cup 2026 API](https://worldcup26.ir). Query teams, groups, matches, stadiums, and live scores from your terminal. Works with the hosted API or any self-hosted instance.

## Install

### Windows (PowerShell)
```powershell
# From cloned repo
.\install.ps1

# Or directly
iex "& { $(Invoke-RestMethod https://raw.githubusercontent.com/Infran/wc26/main/install.ps1) }"
```

### macOS / Linux
```bash
# From cloned repo
chmod +x install.sh && ./install.sh

# Or directly
curl -fsSL https://raw.githubusercontent.com/Infran/wc26/main/install.sh | bash
```

### Go install
```bash
go install github.com/Infran/wc26/cmd/wc26@latest
```

## Quick Start

```bash
# Authenticate (one time)
wc26 auth register "Your Name" you@example.com yourpassword
# or login if already registered:
wc26 auth login you@example.com yourpassword

# Check server health
wc26 health

# List all teams
wc26 teams

# List teams in Group A
wc26 teams --group A

# Get a specific team
wc26 team Brazil
wc26 team 37

# List all groups
wc26 groups

# Get single group
wc26 group A

# List matches
wc26 matches
wc26 matches --type group
wc26 matches --team Brazil
wc26 matches --matchday 1

# Get a specific match
wc26 match 1

# List stadiums
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

Point to a self-hosted instance:

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

Config is stored at:

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
wc26 auth logout
```

## Build from Source

```bash
git clone https://github.com/Infran/wc26.git
cd wc26
go build -ldflags "-X main.Version=$(git describe --tags 2>/dev/null || echo dev)" -o wc26 ./cmd/wc26/
```
