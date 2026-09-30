.DEFAULT_GOAL := help

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# The coverage floor. It is the same number as .github/workflows/ci.yml, and it
# is a ratchet: raise it as coverage grows, never lower it. It starts low
# because this repo does, not because low is acceptable.
COVER_MIN ?= 30

## build: Build the binary
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home ./cmd/adguard-home/

## install: Install to $GOPATH/bin
install:
	CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" ./cmd/adguard-home/

## verify: The gate. A change is done when this exits 0.
verify: fmt-check vet lint test cover-check
	@echo "✓ verify"

## fmt-check: Fail if anything is unformatted
fmt-check:
	@unformatted=$$(gofmt -l . || true); \
	if [ -n "$$unformatted" ]; then echo "✗ needs gofmt:"; echo "$$unformatted"; exit 1; fi
	@echo "✓ formatting clean"

## vet: Run go vet
vet:
	go vet ./...

## cover-check: Fail if coverage falls below COVER_MIN
cover-check:
	@go test -coverprofile=coverage.out ./... > /dev/null
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/ {print $$3}' | tr -d '%'); \
	awk -v t="$$total" -v min="$(COVER_MIN)" 'BEGIN { if (t+0 < min+0) { printf "✗ coverage %.1f%% < %s%%\n", t, min; exit 1 } printf "✓ coverage %.1f%% ≥ %s%%\n", t, min }'

## test: Run all tests
test:
	go test -race -v ./...

## lint: Run golangci-lint
lint:
	golangci-lint run

## clean: Remove build artifacts
clean:
	rm -f adguard-home adguard-home-*

## build-all: Cross-compile for all platforms
build-all:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home-darwin-arm64 ./cmd/adguard-home/
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home-darwin-amd64 ./cmd/adguard-home/
	GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home-linux-amd64 ./cmd/adguard-home/
	GOOS=linux  GOARCH=arm64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home-linux-arm64 ./cmd/adguard-home/
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home-windows-amd64.exe ./cmd/adguard-home/

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

.PHONY: build install test lint clean build-all help
