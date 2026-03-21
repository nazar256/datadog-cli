# ADR-0001: Use a domain-first curated CLI with room for a future raw API namespace

## Status

Accepted

## Context

The CLI is intended to grow toward broad Datadog API coverage, but v1 needs to be practical, discoverable, and safe. A narrow one-off command tree would block future expansion, while exposing raw transport details everywhere would hurt usability for humans and AI agents.

## Decision

Use a curated domain-first command taxonomy in v1:

- top-level nouns for Datadog domains (`monitor`, `dashboard`, `host`, `metric`, `log`)
- consistent verbs (`list`, `get`, `query`, `search`)
- shared runtime, output, and time-range behavior

Reserve `api` as a future namespace for broad/raw coverage if the project later adds a generated or thin-wrapper surface.

## Consequences

- v1 help output stays intuitive and compact.
- Domain adapters can shape stable outputs decoupled from SDK transport envelopes.
- Future raw coverage can be added without breaking the curated UX.
- Some duplication may remain between domain adapters until broader abstractions become justified.

## Alternatives considered

- Raw endpoint-first CLI: scalable for coverage, but poor discoverability and agent ergonomics for v1.
- Single flat command namespace: simpler initially, but grows confusing quickly as coverage expands.
