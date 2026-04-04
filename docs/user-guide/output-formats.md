# Output Formats

All read commands support the `--output` / `-o` flag with three formats.

## Table (default)

```bash
adguard-home clients list
```

```
NAME       IDS                       BLOCKED SERVICES  FILTERING
Desktop    192.168.1.50, 10.0.0.2    none              on
Laptop     192.168.1.51              none              on
```

## JSON

```bash
adguard-home clients list -o json
```

```json
[
  {
    "name": "Desktop",
    "ids": ["192.168.1.50", "10.0.0.2"],
    "blocked_services": [],
    "filtering_enabled": true
  }
]
```

## YAML

```bash
adguard-home stats -o yaml
```

```yaml
num_dns_queries: 245831
num_blocked_filtering: 38104
avg_processing_time: 0.0423
```

## Piping and Scripting

All informational messages (progress, confirmations) go to **stderr**. Structured data goes to **stdout**. This means pipes work cleanly:

```bash
# Extract client names
adguard-home clients list -o json | jq '.[].name'

# Count blocked queries
adguard-home stats -o json | jq '.num_blocked_filtering'

# Export rewrites to a file
adguard-home rewrites list -o json > rewrites-backup.json
```
