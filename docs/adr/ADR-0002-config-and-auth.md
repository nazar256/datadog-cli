# ADR-0002: Use env-first configuration with optional non-overriding local .env loading

## Status

Accepted

## Context

The preferred Datadog auth model for this project is environment-based credentials using `DATADOG_API_KEY` and `DATADOG_APP_KEY`. Local development convenience is also desired through `.env` loading, but explicit environment variables must remain authoritative and secrets must never be committed.

## Decision

- Use `DATADOG_API_KEY` and `DATADOG_APP_KEY` as the only supported secret inputs in v1.
- Support `DATADOG_SITE` as the canonical site selector.
- Support `--site` as a global override flag.
- Load a local `.env` file only from the current working directory or an explicitly supplied path.
- `.env` loading must not override already-set process environment variables.
- Do not store persistent config on disk in v1.

Precedence:

1. explicit CLI flags
2. process environment variables
3. `.env` values
4. built-in defaults

## Consequences

- Configuration stays simple and easy to explain.
- Local development is convenient without surprising production-style environments.
- Secrets remain outside the CLI flags surface and out of shell history.
- Users wanting named profiles or persistent config will need a later enhancement.

## Alternatives considered

- Persistent config file: more featureful, but unnecessary for v1 and adds state-management complexity.
- Flags for API/app keys: worse security ergonomics and more likely to leak into shell history.
