# Multi-Instance Support

Manage multiple AdGuard Home servers from a single CLI installation.

## Add Instances

Run `adguard-home setup` for each instance, giving each a unique name:

```bash
# First instance (becomes "default")
adguard-home setup
# Enter: http://192.168.0.105:8001, admin, password, "default"

# Second instance
adguard-home setup
# Enter: http://10.0.0.1:3000, admin, password, "secondary"
```

## Switch Instances

Use the `--instance` flag on any command:

```bash
# Query the default instance
adguard-home clients list

# Query the secondary instance
adguard-home --instance secondary clients list

# Compare stats between instances
adguard-home stats -o json
adguard-home --instance secondary stats -o json
```

## Config File

```yaml
instances:
  default:
    url: http://192.168.0.105:8001
    username: admin
  secondary:
    url: http://10.0.0.1:3000
    username: admin
current_instance: default
```

Each instance's password is stored separately in the credential store, keyed by instance name.
