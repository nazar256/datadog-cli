# For AI agents

`ddog` is useful when an agent needs Datadog access from a terminal, script, or sandbox where MCP is unavailable or impractical.

## Authenticate

Use environment variables or an explicit `.env` file:

```bash
export DATADOG_API_KEY=...
export DATADOG_APP_KEY=...
export DATADOG_SITE=datadoghq.com
```

Or:

```bash
ddog --env-file .env config doctor
```

## Discover the command tree

Start with built-in help and docs:

```bash
ddog --help
ddog docs summary
ddog docs commands --output json
ddog monitor --help
ddog log search --help
```

Use `--help` for the actual command tree. `ddog docs commands --output json` is a compact machine-readable summary of the command taxonomy, not a full command listing.

## Prefer machine-readable output when parsing results

```bash
ddog version --output json
ddog config doctor --output json
ddog monitor list --limit 10 --output json
ddog log search --query 'service:web status:error' --last 15m --output json
```

## Good agent workflows

Check auth and site before live calls:

```bash
ddog config doctor --output json
```

Discover monitors related to a service:

```bash
ddog monitor list --name api --limit 20 --output json
```

Pull recent logs for a narrow incident query:

```bash
ddog log search --query 'service:web status:error' --last 15m --limit 20 --output json
```

Inspect a metric window:

```bash
ddog metric query --query 'avg:system.cpu.user{env:prod}' --last 1h --output json
```

Current scope is intentionally read-oriented: monitors, dashboards, hosts, metrics, and logs.
