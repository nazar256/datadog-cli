# Datadog CLI usage

## Quick start

Build the CLI:

```bash
go build -o ddog ./cmd/ddog
```

Check local configuration:

```bash
./ddog config doctor
./ddog --env-file .env config doctor --output json
```

## Authentication

`ddog` reads credentials from:

- `DATADOG_API_KEY`
- `DATADOG_APP_KEY`
- optional `DATADOG_SITE`

You can also use a local `.env` file in the current directory, or pass an explicit file with `--env-file`.

Secrets are never accepted as CLI flags.

## Command discovery

Use help output as the primary interface:

```bash
ddog --help
ddog docs summary
ddog docs commands --output json
ddog monitor --help
ddog log search --help
```

Shell completion is also available:

```bash
ddog completion bash
ddog completion zsh
ddog completion fish
```

## Output modes

- default: concise text for terminals
- `--output json`: stable machine-readable output

Examples:

```bash
ddog monitor list --output json
ddog docs commands --output json
```

## Read-only v1 commands

### Monitors

```bash
ddog monitor list
ddog monitor list --name api --limit 20
ddog monitor get 123456
```

### Dashboards

```bash
ddog dashboard list --count 20
ddog dashboard get abc-def-ghi
```

### Hosts

```bash
ddog host list --filter web
ddog host get web-01
```

### Metrics

```bash
ddog metric query --query 'avg:system.load.1{*}' --last 1h
ddog metric query --query 'avg:system.cpu.user{env:prod}' --from 2026-03-21T09:00:00Z --to 2026-03-21T10:00:00Z
```

### Logs

```bash
ddog log search --query 'service:web status:error' --last 15m
ddog log search --query 'env:prod' --index main --limit 20 --output json
```

## Notes

- All shipped v1 Datadog commands are read-only.
- Empty results are valid outcomes.
- Supported sites are: `us1`, `us3`, `us5`, `eu`, `ap1`, `ap2`, `us1-fed`, and their canonical hostnames.
- Install and release details live in [install.md](install.md).
- AI-agent-specific guidance lives in [for-ai-agents.md](for-ai-agents.md).
