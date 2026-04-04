# Contributing

## Prerequisites

- Go 1.24+
- Make

## Setup

```bash
git clone https://github.com/jjuanrivvera/adguard-cli.git
cd adguard-cli
make build
make test
```

## Development Workflow

1. Create a branch: `git checkout -b feat/my-feature`
2. Make changes
3. Run tests: `make test`
4. Build: `make build`
5. Test against a real AdGuard Home instance
6. Commit with conventional commits: `feat:`, `fix:`, `docs:`, `test:`
7. Push and open a PR

## Conventional Commits

```
feat: add dhcp lease management commands
fix: handle nil response in query log
docs: add multi-instance guide
test: add credential encryption tests
chore: update dependencies
```

## Code Style

- One command file per API resource group
- Keep command files under 200 lines
- Use `CLIError` with hints for all user-facing errors
- `cmdutil.Infof` for messages (stderr), `output.Print` for data (stdout)
- Tests use `httptest.NewServer` for API mocking
- Run `go vet ./...` before committing

## Testing

```bash
# Run all tests
make test

# Run with race detection
go test -race ./...

# Coverage
go test -cover ./internal/...
```
