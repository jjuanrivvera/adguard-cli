# Architecture

## Project Structure

```
adguard-cli/
├── cmd/adguard-home/main.go    # Entry point
├── commands/                    # Cobra command definitions (one file per resource)
├── internal/
│   ├── api/client.go            # AdGuard Home REST API client
│   ├── config/
│   │   ├── config.go            # YAML config management
│   │   └── credentials.go       # Keyring + encrypted file credential store
│   ├── errors/errors.go         # CLIError with hints
│   ├── output/formatter.go      # Table/JSON/YAML output
│   └── cmdutil/util.go          # Stderr helpers
├── tools/gendocs/main.go       # Auto-generate CLI reference docs
├── docs/                        # MkDocs documentation
├── Makefile                     # Build, test, lint, cross-compile
├── .goreleaser.yml              # Release automation
└── .github/workflows/           # CI/CD
```

## Key Design Patterns

### Constructor Root Command

The root command uses `NewRootCommand(version, commit, date)` returning `*cobra.Command`. No package-level state for the command tree, making it testable.

### CLIError with Hints

Every error includes an actionable hint:

```go
errors.ConnectionFailed(url, err)
// Output:
// Error [CONNECTION_ERROR]: cannot connect to AdGuard Home at http://...
//   Cause: dial tcp: connection refused
//   Hint: Check that AdGuard Home is running. Run 'adguard-home doctor' to diagnose.
```

### stderr/stdout Separation

- `cmdutil.Infof()` / `cmdutil.Infoln()` write to stderr (progress, confirmations)
- `output.Print()` / `output.PrintJSON()` write to stdout (structured data)

This makes piping reliable: `adguard-home clients list -o json | jq .`

### Credential Store Interface

```go
type CredentialStore interface {
    Get(instance string) (string, error)
    Set(instance, password string) error
    Delete(instance string) error
}
```

Two implementations: `keyringStore` (system keyring) and `encryptedFileStore` (AES-256-GCM fallback). The `NewCredentialStore()` factory auto-detects which is available.

## API Client

`internal/api/client.go` wraps AdGuard Home's REST API. All endpoints use the `/control/` prefix with HTTP Basic Auth. The client handles JSON serialization, error responses (401/403 → auth error, 4xx → API error), and 15-second timeouts.

## Adding a New Command

1. Create `commands/{resource}.go`
2. Define `func new{Resource}Cmd() *cobra.Command` with subcommands
3. Add API methods to `internal/api/client.go` if needed
4. Register in `root.go` → `root.AddCommand(new{Resource}Cmd())`
5. Use `output.Print(getFormat(), data, headers, rowBuilder)` for output
6. Use `clierrors.WithHint()` for user-friendly errors
7. Run `go run ./tools/gendocs/main.go` to update CLI reference
