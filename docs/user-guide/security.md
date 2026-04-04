# Security

## Credential Storage

Passwords are **never** stored in the config YAML file. The CLI uses a two-tier credential store:

### Tier 1: System Keyring (default)

On systems with a keyring daemon (macOS Keychain, GNOME Keyring, KWallet), passwords are stored in the OS credential manager. This is the most secure option.

### Tier 2: Encrypted File (fallback)

On headless servers without a keyring (e.g., Docker containers, CI runners, VPS), passwords are stored in AES-256-GCM encrypted files at `~/.adguard-cli/credentials.enc.<instance>`.

Key derivation uses PBKDF2-SHA256 with 100,000 iterations. The key material is derived from the hostname and instance name, providing obfuscation that prevents casual access but is not a security boundary against a determined attacker with file access.

### Tier 3: Environment Variables

For CI/CD pipelines, set `ADGUARD_URL`, `ADGUARD_USERNAME`, and `ADGUARD_PASSWORD` as environment variables. These take precedence over both the keyring and config file.

## Network Security

The CLI communicates with AdGuard Home over HTTP Basic Auth. For production deployments, consider:

- Running AdGuard Home behind a reverse proxy with TLS
- Restricting access to the admin API via firewall rules
- Using a VPN (Tailscale, WireGuard) for remote management

## File Permissions

- `~/.adguard-cli/config.yaml` is created with `0600` permissions (owner read/write only)
- `~/.adguard-cli/` directory is created with `0700` permissions
- Encrypted credential files are created with `0600` permissions
