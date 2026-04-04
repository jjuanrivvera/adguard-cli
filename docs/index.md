# AdGuard CLI

The missing command-line interface for [AdGuard Home](https://github.com/AdguardTeam/AdGuardHome).

Manage your AdGuard Home instance from the terminal: clients, blocked services, DNS rewrites, query logs, filters, DHCP, TLS, and more. Full API coverage with structured output.

## Features

- **Full API coverage** -- 90%+ of AdGuard Home's 81 API operations
- **Structured output** -- table, JSON, or YAML for every command
- **Multi-instance** -- manage multiple AdGuard Home servers from one CLI
- **Secure credentials** -- system keyring with encrypted file fallback
- **Cross-platform** -- macOS, Linux, Windows (amd64 and arm64)
- **AI-friendly** -- JSON output and stderr/stdout separation for automation

## Quick Example

```bash
# Check server status
adguard-home status

# List all clients
adguard-home clients list

# Block TikTok globally
adguard-home services block tiktok

# Check if YouTube is blocked
adguard-home dns check youtube.com

# View DNS stats as JSON
adguard-home stats -o json
```

## Why

AdGuard Home has a web UI and a REST API, but no official CLI. The only existing CLI covers about 20% of the API. This project covers 90%+ with proper credential security, multi-instance support, and structured output for scripting.

Built for homelab operators, sysadmins, and anyone who automates their DNS infrastructure.
