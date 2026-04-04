# adguard-cli

The missing command-line interface for [AdGuard Home](https://github.com/AdguardTeam/AdGuardHome).

Manage your AdGuard Home instance from the terminal: clients, blocked services, DNS rewrites, query logs, filters, DHCP, TLS, and more. Full API coverage with structured output.

## Why

AdGuard Home has a web UI and a REST API, but no official CLI. The only existing CLI ([adctl](https://github.com/ewosborne/adctl)) covers ~20% of the API. This project covers 90%+ of AdGuard Home's 81 API operations.

Built for homelab operators, sysadmins, and anyone who automates their DNS infrastructure.

## Install

### From source

```bash
go install github.com/jjuanrivvera/adguard-cli/cmd/adguard-home@latest
```

### From binary

Download the latest release from the [Releases](https://github.com/jjuanrivvera/adguard-cli/releases) page.

## Quick Start

```bash
# Configure your instance
adguard-home setup

# Check connectivity
adguard-home doctor

# View server status
adguard-home status

# List configured clients
adguard-home clients list

# Check DNS stats
adguard-home stats
```

Or use environment variables for CI/automation:

```bash
export ADGUARD_URL="http://192.168.0.105:8001"
export ADGUARD_USERNAME="admin"
export ADGUARD_PASSWORD="your-password"

adguard-home clients list -o json | jq '.[].name'
```

## Commands

| Command | Description |
|---------|-------------|
| `status` | Server status, enable/disable protection |
| `stats` | DNS query statistics, reset |
| `clients` | List, find, add, delete clients |
| `services` | List, block, unblock services globally |
| `rewrites` | List, add, delete DNS rewrites |
| `log` | View DNS query log |
| `filters` | List, add, remove, refresh filter lists |
| `dhcp` | DHCP status, leases, static lease management |
| `tls` | TLS/HTTPS configuration status |
| `dns` | DNS config, cache clear, host blocking check |
| `safebrowsing` | Enable/disable safe browsing |
| `parental` | Enable/disable parental control |
| `safesearch` | Safe search enforcement per engine |
| `access` | Allowed/disallowed clients and blocked hosts |
| `check-update` | Check for AdGuard Home updates |
| `update` | Trigger AdGuard Home update |
| `doctor` | Run diagnostic checks |
| `setup` | Interactive configuration wizard |

## Output Formats

All read commands support `--output` / `-o`:

```bash
# Table (default)
adguard-home clients list

# JSON (for scripting)
adguard-home clients list -o json

# YAML
adguard-home stats -o yaml
```

## Examples

### Client Management

```bash
# List all clients with their IDs and blocked services
adguard-home clients list

# Find which client owns an IP
adguard-home clients find 192.168.0.57

# Add a new client
adguard-home clients add "Smart TV" "192.168.0.110,192.168.0.71"

# Delete a client
adguard-home clients delete "Smart TV"
```

### Blocked Services

```bash
# See what's blocked globally
adguard-home services blocked

# Block TikTok and Instagram
adguard-home services block tiktok,instagram

# Unblock YouTube
adguard-home services unblock youtube
```

### DNS Rewrites

```bash
# List all rewrites
adguard-home rewrites list

# Add a local DNS entry
adguard-home rewrites add homelab.local 192.168.0.100

# Remove a rewrite
adguard-home rewrites delete homelab.local 192.168.0.100
```

### DNS Diagnostics

```bash
# Check if a domain is blocked
adguard-home dns check youtube.com

# View recent query log
adguard-home log -n 50

# Clear DNS cache
adguard-home dns cache-clear
```

### Filter Lists

```bash
# List all filter lists with rule counts
adguard-home filters list

# Add a new filter list
adguard-home filters add "OISD Big" "https://big.oisd.nl"

# Refresh all filters
adguard-home filters refresh
```

### DHCP

```bash
# Show DHCP status and lease count
adguard-home dhcp status

# List all leases
adguard-home dhcp leases

# Add a static lease
adguard-home dhcp add-lease "AA:BB:CC:DD:EE:FF" "192.168.0.50" "my-server"
```

## Configuration

Config lives at `~/.adguard-cli/config.yaml`:

```yaml
instances:
  default:
    url: http://192.168.0.105:8001
    username: admin
    password: your-password
  secondary:
    url: http://10.0.0.1:3000
    username: admin
    password: other-password
current_instance: default
output:
  format: table
  color: auto
```

Switch instances: `adguard-home --instance secondary clients list`

### Environment Variables

| Variable | Description |
|----------|-------------|
| `ADGUARD_URL` | AdGuard Home base URL |
| `ADGUARD_USERNAME` | HTTP Basic auth username |
| `ADGUARD_PASSWORD` | HTTP Basic auth password |

Environment variables override the config file.

## API Coverage

| Category | Operations | Covered |
|----------|-----------|---------|
| Status/Global | status, protection toggle, DNS config, cache clear | Yes |
| Clients | list, find, add, update, delete | Yes |
| Blocked Services | list all, list blocked, get, set | Yes |
| DNS Rewrites | list, add, delete | Yes |
| Query Log | query, clear | Yes |
| Filtering | status, add, remove, refresh, check host | Yes |
| DHCP | status, interfaces, leases, static lease CRUD, reset | Yes |
| TLS | status, configure | Yes |
| Safe Browsing | status, enable, disable | Yes |
| Parental | status, enable, disable | Yes |
| Safe Search | status per engine | Yes |
| Access Control | list allowed/disallowed/blocked | Yes |
| Version | check update, trigger update | Yes |
| Stats | get, reset | Yes |

## Building from Source

```bash
git clone https://github.com/jjuanrivvera/adguard-cli.git
cd adguard-cli
go build -o adguard-home ./cmd/adguard-home/
```

Cross-compile:

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o adguard-home-darwin-arm64 ./cmd/adguard-home/

# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o adguard-home-linux-amd64 ./cmd/adguard-home/

# Windows
GOOS=windows GOARCH=amd64 go build -o adguard-home.exe ./cmd/adguard-home/
```

## License

MIT
