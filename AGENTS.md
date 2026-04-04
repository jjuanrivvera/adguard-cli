# AGENTS.md

Instructions for AI coding assistants working on this repository.

## Project Overview

`adguard-cli` is a Go CLI for managing AdGuard Home DNS filtering instances via their REST API. It provides full API coverage (~90%+ of AdGuard Home's 81 API operations) with structured output (table, JSON, YAML), multi-instance support, and an interactive setup wizard.

## Architecture

```
adguard-cli/
├── cmd/adguard-home/main.go    # Entry point — injects version, calls NewRootCommand()
├── commands/                    # One file per command group (Cobra commands)
│   ├── root.go                  # NewRootCommand() constructor, GlobalFlags, getClient()
│   ├── status.go                # Server status + protection toggle
│   ├── stats.go                 # DNS query statistics
│   ├── clients.go               # Client CRUD (list/find/add/delete)
│   ├── services.go              # Globally blocked services
│   ├── rewrites.go              # DNS rewrites CRUD
│   ├── log.go                   # Query log viewer
│   ├── filters.go               # Filter lists management
│   ├── dhcp.go                  # DHCP server management
│   ├── tls.go                   # TLS/HTTPS configuration
│   ├── dns.go                   # DNS config, cache, host checking
│   ├── safebrowsing.go          # Safe browsing toggle
│   ├── parental.go              # Parental control toggle
│   ├── safesearch.go            # Safe search per-engine config
│   ├── access.go                # Access control lists
│   ├── version.go               # Update check + trigger
│   ├── doctor.go                # Diagnostic checks
│   └── setup.go                 # Interactive first-run wizard
├── internal/
│   ├── api/client.go            # AdGuard Home REST API client (all endpoints)
│   ├── config/config.go         # Viper config (~/.adguard-cli/config.yaml)
│   ├── errors/errors.go         # CLIError with Code + Hint pattern
│   ├── output/formatter.go      # Table/JSON/YAML output formatters
│   └── cmdutil/util.go          # Infof/Infoln (stderr), HandleError
└── test/                        # Unit and integration tests
```

## Key Patterns

### Constructor root command
The root command uses `NewRootCommand(version, commit, date)` returning `*cobra.Command` — testable, no package-level state.

### CLIError with hints
Every error includes a `Hint` field with actionable fix instructions:
```go
errors.ConnectionFailed(url, err)
// → "Cannot connect to AdGuard Home at http://... \n Hint: Check that AdGuard Home is running..."
```

### stderr/stdout separation
- `cmdutil.Infof()` / `cmdutil.Infoln()` → stderr (progress, confirmations)
- `output.Print()` / `output.PrintJSON()` → stdout (structured data)

This means `adguard-home clients list -o json | jq .` works cleanly.

### Output formatting
All read commands support `--output/-o` with `table` (default), `json`, or `yaml`. The `output.Print()` function takes a format, raw data, table headers, and a row-builder function.

### Config + env var precedence
1. Environment variables (`ADGUARD_URL`, `ADGUARD_USERNAME`, `ADGUARD_PASSWORD`)
2. Config file (`~/.adguard-cli/config.yaml`)

Env vars are checked first for CI/automation. Config file supports multiple instances.

## Development Commands

```bash
# Build
go build -o adguard-home ./cmd/adguard-home/

# Run
./adguard-home doctor
./adguard-home clients list -o json

# Test
go test ./...

# Lint
golangci-lint run

# Cross-compile
GOOS=darwin GOARCH=arm64 go build -o adguard-home-darwin-arm64 ./cmd/adguard-home/
GOOS=linux GOARCH=amd64 go build -o adguard-home-linux-amd64 ./cmd/adguard-home/
```

## Adding a New Command

1. Create `commands/{resource}.go`
2. Define `func new{Resource}Cmd() *cobra.Command` with subcommands
3. Add API methods to `internal/api/client.go` if needed
4. Register in `root.go` → `root.AddCommand(new{Resource}Cmd())`
5. Use `output.Print(getFormat(), data, headers, rowBuilder)` for output
6. Use `clierrors.WithHint()` for user-friendly errors
7. Use `cmdutil.Infof()` for confirmations (to stderr)

## API Client

`internal/api/client.go` wraps the AdGuard Home REST API with typed methods. All endpoints use `/control/` prefix. Auth is HTTP Basic. The client handles JSON marshaling, error responses, and timeouts (15s default).

## Conventions

- One command file per API resource group, keep under 200 lines
- Command files define both the parent and subcommands
- No package-level state — pass everything via constructors or function args
- Error messages must include hints for the user
- All data output goes to stdout, all messages to stderr
- Tests use `httptest.NewServer` for API mocking
