.DEFAULT_GOAL := help

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

## build: Build the binary
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o adguard-home ./cmd/adguard-home/

## install: Install to $GOPATH/bin
install:
	CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" ./cmd/adguard-home/

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
