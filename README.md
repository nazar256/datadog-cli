# ddog

`ddog` is a single-binary Datadog CLI for AI agents, automation, and terminal-first engineers.

It is designed for the cases where you want Datadog access from a shell, script, CI job, or coding agent runtime where MCP is unavailable, inconvenient, or unnecessary.

`ddog` uses the official Datadog Go SDK, authenticates with environment variables, stays self-discoverable through `--help`, and supports stable JSON output for automation.

## Why this exists

- **MCP fallback for agents**: useful when an agent can run terminal commands but cannot use a Datadog MCP server
- **Fast terminal workflows**: inspect monitors, dashboards, hosts, metrics, and logs without leaving the shell
- **Automation-friendly output**: use concise text for humans or `--output json` for scripts and agents
- **Help-driven discovery**: the command tree is meant to be explored directly from the binary

Current scope is intentionally **read-only**.

## Install

### Recommended: latest release installer

```bash
curl -fsSL https://github.com/nazar256/datadog-cli/releases/latest/download/install.sh | sh
```

This installer selects the correct release archive for your platform and verifies its SHA256 checksum before installing.

### Build from source

```bash
go build -o ddog ./cmd/ddog
./ddog --help
```

### Go install

```bash
go install github.com/nazar256/datadog-cli/cmd/ddog@latest
```

`go install` is useful for source-based workflows, but release binaries are the primary install path and include embedded version metadata.

More install details: [docs/install.md](docs/install.md)

## Authentication

`ddog` reads Datadog credentials from:

- `DATADOG_API_KEY`
- `DATADOG_APP_KEY`
- optional `DATADOG_SITE`

You can also point to a local env file with `--env-file`. By default, `ddog` reads `.env` from the current working directory only.

```bash
export DATADOG_API_KEY=...
export DATADOG_APP_KEY=...
export DATADOG_SITE=datadoghq.com
ddog config doctor
```

Secrets are never accepted as CLI flags.

## Discover commands

Start with built-in help and docs:

```bash
ddog --help
ddog docs summary
ddog docs commands --output json
ddog monitor --help
ddog log search --help
```

## Real examples

```bash
# verify auth, site, and output mode
ddog config doctor --output json

# inspect monitor coverage for a service
ddog monitor list --name api --limit 20 --output json

# fetch dashboards in concise terminal output
ddog dashboard list --count 20

# query a recent metric window
ddog metric query --query 'avg:system.load.1{*}' --last 1h

# search recent logs for an incident query
ddog log search --query 'service:web status:error' --last 15m --limit 20 --output json
```

## Use with AI agents

`ddog` works well for agents that need Datadog access through ordinary shell commands.

Recommended agent flow:

1. Run `ddog --help` for the command tree, and `ddog docs commands --output json` for high-level command taxonomy guidance.
2. Run `ddog config doctor --output json` before live Datadog calls.
3. Prefer `--output json` whenever the result will be parsed.
4. Keep queries narrow and explicit, especially for logs and metrics.

Examples:

```bash
ddog version --output json
ddog config doctor --output json
ddog monitor list --limit 10 --output json
ddog log search --query 'service:web status:error' --last 15m --limit 20 --output json
```

More: [docs/for-ai-agents.md](docs/for-ai-agents.md)

## Output modes

- default: concise terminal text
- `--output json`: stable machine-readable output

Useful JSON entry points:

```bash
ddog version --output json
ddog docs commands --output json
ddog monitor list --output json
```

## Supported v1 command areas

- `config doctor`
- `docs`
- `version`
- `monitor list|get`
- `dashboard list|get`
- `host list|get`
- `metric query`
- `log search`

## Releases

Release archives are intended for:

- Linux amd64
- Linux arm64
- macOS amd64
- macOS arm64

Linux is the main release target today.

Download binaries and checksums from:

- <https://github.com/nazar256/datadog-cli/releases>

## Documentation

- [docs/install.md](docs/install.md)
- [docs/usage.md](docs/usage.md)
- [docs/for-ai-agents.md](docs/for-ai-agents.md)
- [docs/publish-checklist.md](docs/publish-checklist.md)

## Development

```bash
go test ./...
go build ./cmd/ddog
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for lightweight contribution guidance.

## License

[MIT](LICENSE)
