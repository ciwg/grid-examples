# TODO lavom - ex6 inventory package

## Decision Intent Log

ID: DI-lavom
Date: 2026-07-28 18:00:00
Status: active
Decision: Implement the first inventory package as a first-party built-in package with explicit inventory item, count, and reconciliation families over shared context, knowledge, and runs records.
Intent: Add the next real domain package without collapsing inventory behavior into the runtime or hiding count and reconciliation truth inside generic run notes.
Constraints: Keep the package under `packages/inventory/`; preserve runtime-mediated durable history; compose with `context`, `knowledge`, and `runs` instead of duplicating their cores.
Affects: `packages/inventory`, `cmd/moks`, runtime tests, and ex6 status docs.

## Goal

Build the first-party `inventory` package.

## Scope

- inventory item kind
- count recording
- reconciliation recording
- inventory inspection and listing
- durable family declarations

## Why

Inventory is the next real ex5-derived domain surface and should exist as its
own package instead of waiting behind later embodiments.
