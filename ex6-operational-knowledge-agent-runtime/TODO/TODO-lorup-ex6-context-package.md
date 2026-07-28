# TODO lorup - ex6 context package

## Decision Intent Log

ID: DI-lorup
Date: 2026-07-28 13:00:00
Status: active
Decision: Implement the first context package as a first-party built-in package while the installed-package mutation contract remains too shallow for serious durable writes.
Intent: Get the first real ex6 package into code without pretending the external package ABI is further along than it is.
Constraints: Keep the package under `packages/context/`; preserve the runtime contract; do not bypass durable history.
Affects: `packages/context`, `cmd/moks`, runtime tests, and ex6 current-state docs.

## Goal

Build the first-party `context` package.

## Scope

- responsibilities
- places
- resources
- command surface for create, inspect, and list
- durable family declarations

## Why

This package gives the rest of ex6 a shared operational context layer.
