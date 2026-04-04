# Output Formats

All read commands support the `--output` / `-o` flag with three formats.

## Table (default)

```bash
adguard-home clients list
```

```
NAME    IDS                       BLOCKED SERVICES  FILTERING
JUAN    192.168.0.57, 100.83.15   none              on
VANE    192.168.0.4               none              on
```

## JSON

```bash
adguard-home clients list -o json
```

```json
[
  {
    "name": "JUAN",
    "ids": ["192.168.0.57", "100.83.15.72"],
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
num_dns_queries: 782115
num_blocked_filtering: 125145
avg_processing_time: 0.1247
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
