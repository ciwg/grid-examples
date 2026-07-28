# TODO figar - ex6 links package

## Decision Intent Log

ID: DI-figar
Date: 2026-07-28 14:20:00
Status: active
Decision: Implement the first links package as a first-party built-in egg with one typed-link family while the installed-package mutation contract remains shallow.
Intent: Get cross-record relationships into ex6 code now so packages can relate durable records without private side channels.
Constraints: Keep the package under `packages/links/`; preserve durable history writes; keep the relation model explicit and typed.
Affects: `packages/links`, `cmd/moks`, runtime tests, and ex6 current-state docs.

## Goal

Build the first-party `links` package.

## Scope

- typed links
- link validation
- command surface for creation and inspection
- durable family declarations

## Why

This package keeps cross-record relationships explicit instead of burying them
inside package-private payloads.
