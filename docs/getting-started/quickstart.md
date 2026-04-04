# Quick Start

After [installing](installation.md) and [configuring](configuration.md), try these commands:

## Check Your Setup

```bash
adguard-home doctor
adguard-home status
adguard-home stats
```

## Manage Clients

```bash
# List all clients
adguard-home clients list

# Find which client owns an IP
adguard-home clients find 192.168.0.57

# Add a new client
adguard-home clients add "Smart TV" "192.168.0.110,192.168.0.71"
```

## Manage Blocked Services

```bash
# See what's blocked globally
adguard-home services blocked

# Block or unblock services
adguard-home services block tiktok,instagram
adguard-home services unblock youtube
```

## DNS Operations

```bash
# Check if a domain is blocked
adguard-home dns check youtube.com

# Add a local DNS entry
adguard-home rewrites add homelab.local 192.168.0.100

# View recent queries
adguard-home log -n 20

# Clear DNS cache
adguard-home dns cache-clear
```

## Use JSON for Scripting

```bash
# Pipe to jq
adguard-home clients list -o json | jq '.[].name'

# Count blocked queries
adguard-home stats -o json | jq '.num_blocked_filtering'
```
