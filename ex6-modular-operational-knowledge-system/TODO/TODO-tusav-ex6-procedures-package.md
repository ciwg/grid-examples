# TODO tusav - ex6 procedures package

## Decision Intent Log

ID: DI-tusav
Date: 2026-07-28 14:45:00
Status: active
Decision: Implement the first procedures package as a first-party built-in domain package that composes the shared knowledge and runs families while also declaring procedure-specific families.
Intent: Prove that business behavior can live above the shared packages instead of collapsing back into the runtime core.
Constraints: Keep the package under `packages/procedures/`; preserve durable history writes; keep procedure semantics above the shared generic families.
Affects: `packages/procedures`, `cmd/moks`, runtime tests, and ex6 current-state docs.

## Goal

Build the first real domain package: `procedures`.

## Scope

- procedure item kind over the knowledge package
- procedure-specific commands
- example package records and workflows
- package docs showing how a real package is structured

## Why

This package should become the first serious proof that ex6 can carry ex5-like
business behavior as a modular package.
