# Suggested GitHub metadata

These settings are configured in the GitHub UI, not from repository files.

## Repository description

Suggested description:

> Single-binary Datadog CLI for AI agents, automation, and terminal-first engineers when MCP is unavailable.

## Homepage

Suggested homepage until a dedicated docs site exists:

> `https://github.com/nazar256/datadog-cli#readme`

## Topics

Suggested topics:

- `datadog`
- `cli`
- `automation`
- `ai-agents`
- `terminal`
- `observability`
- `logs`
- `metrics`
- `monitors`
- `dashboards`
- `golang`

## Social preview direction

Use a simple terminal screenshot or mockup that shows:

- `ddog --help`
- one JSON example such as `ddog log search --output json`
- a short caption like: `Datadog CLI for AI agents and automation when MCP is unavailable`

## First public release notes should cover

1. what `ddog` is and who it is for
2. supported Datadog surfaces in v1: monitors, dashboards, hosts, metrics, logs
3. install options: release installer, release assets, source build
4. auth model: `DATADOG_API_KEY`, `DATADOG_APP_KEY`, optional `DATADOG_SITE`
5. help-driven discovery and `--output json`
6. checksum-verified release assets and supported platforms
