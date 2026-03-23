# ddog

`ddog` is a Go-based Datadog CLI for humans and coding agents.

It uses the official Datadog Go SDK, supports env-first authentication, optional local `.env` loading, concise terminal output, and stable JSON output for automation.

## Features in v1

- env-first auth via `DATADOG_API_KEY` and `DATADOG_APP_KEY`
- optional local `.env` loading
- explicit Datadog site selection
- self-discoverable help and built-in docs
- read-only commands for monitors, dashboards, hosts, metrics, and logs
- automated tests

## Installation

### Quick Install (macOS/Linux)

You can install the latest release using the provided install script:

```bash
curl -fsSL https://raw.githubusercontent.com/nazar256/datadog-cli/main/install.sh | sh
```

### Go Install

If you have Go installed, you can build and install it directly:

```bash
go install github.com/nazar256/datadog-cli/cmd/ddog@latest
```

## Build

```bash
go build -o ddog ./cmd/ddog
```

## Examples

```bash
ddog config doctor
ddog docs summary
ddog monitor list --output json
ddog metric query --query 'avg:system.load.1{*}' --last 1h
```

See [docs/usage.md](docs/usage.md) for more examples.
