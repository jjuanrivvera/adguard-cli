# Installation

## Homebrew (macOS/Linux)

```bash
brew tap jjuanrivvera/adguard-cli
brew install adguard-cli
```

## Go Install

```bash
go install github.com/jjuanrivvera/adguard-cli/cmd/adguard-home@latest
```

## Binary Download

Download the latest release for your platform from the [Releases](https://github.com/jjuanrivvera/adguard-cli/releases) page.

### Linux (amd64)

```bash
curl -LO https://github.com/jjuanrivvera/adguard-cli/releases/latest/download/adguard-cli_linux_amd64.tar.gz
tar xzf adguard-cli_linux_amd64.tar.gz
sudo mv adguard-home /usr/local/bin/
```

### macOS (Apple Silicon)

```bash
curl -LO https://github.com/jjuanrivvera/adguard-cli/releases/latest/download/adguard-cli_darwin_arm64.tar.gz
tar xzf adguard-cli_darwin_arm64.tar.gz
sudo mv adguard-home /usr/local/bin/
```

## From Source

```bash
git clone https://github.com/jjuanrivvera/adguard-cli.git
cd adguard-cli
make build
sudo mv adguard-home /usr/local/bin/
```

## Verify Installation

```bash
adguard-home --version
```
