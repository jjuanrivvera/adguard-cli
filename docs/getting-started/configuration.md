# Configuration

## Interactive Setup

The fastest way to configure your instance:

```bash
adguard-home setup
```

This wizard will:

1. Ask for your AdGuard Home URL (e.g., `http://192.168.0.105:8001`)
2. Ask for your username and password (password is masked)
3. Test the connection
4. Save the config to `~/.adguard-cli/config.yaml`
5. Store the password securely in your system keyring

## Config File

The configuration lives at `~/.adguard-cli/config.yaml`:

```yaml
instances:
  default:
    url: http://192.168.0.105:8001
    username: admin
  secondary:
    url: http://10.0.0.1:3000
    username: admin
current_instance: default
output:
  format: table
  color: auto
```

!!! warning "Passwords are never stored in the config file"
    Passwords are stored in your system keyring (macOS Keychain, GNOME Keyring, KWallet) or in an AES-256-GCM encrypted file as fallback on headless servers.

## Environment Variables

For CI/CD and automation, use environment variables:

| Variable | Description |
|----------|-------------|
| `ADGUARD_URL` | AdGuard Home base URL |
| `ADGUARD_USERNAME` | HTTP Basic auth username |
| `ADGUARD_PASSWORD` | HTTP Basic auth password |

Environment variables take precedence over the config file.

```bash
export ADGUARD_URL="http://192.168.0.105:8001"
export ADGUARD_USERNAME="admin"
export ADGUARD_PASSWORD="your-password"

adguard-home clients list -o json
```

## Verify Configuration

```bash
adguard-home doctor
```

This runs 4 diagnostic checks: config, connectivity, authentication, and protection status.
