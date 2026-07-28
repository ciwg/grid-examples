# TODO ramek - ex6 maintenance package

## Decision Intent Log

ID: DI-ramek
Date: 2026-07-28 17:20:00
Status: active
Decision: Implement the first maintenance package as a first-party built-in package with explicit maintenance item, service, and finding families over shared context, knowledge, and runs records.
Intent: Add the next real domain package without collapsing maintenance behavior into the runtime or hiding resource-specific findings inside generic run notes.
Constraints: Keep the package under `packages/maintenance/`; preserve runtime-mediated durable history; compose with `context`, `knowledge`, and `runs` instead of duplicating their cores.
Affects: `packages/maintenance`, `cmd/moks`, runtime tests, and ex6 status docs.

## Goal

Build the first-party `maintenance` package.

## Scope

- maintenance item kind
- maintenance service recording
- maintenance finding recording
- maintenance inspection and listing
- durable family declarations

## Why

Maintenance is one of the next real ex5-derived domain surfaces and should
exist as its own package instead of waiting behind later embodiments.
