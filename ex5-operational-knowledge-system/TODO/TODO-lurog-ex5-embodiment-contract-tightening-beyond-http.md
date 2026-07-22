# TODO lurog - tighten the ex5 embodiment contract beyond the current HTTP adapter

## Decision Intent Log

ID: DI-bavuk
Date: 2026-07-22 14:10:17 -0700
Status: active
Decision: Implement `112A`: keep HTTP as the sole embodiment adapter for now,
and tighten the embodiment contract by exposing the richer runtime/capability
boundary more explicitly instead of adding a second transport surface.
Intent: Keep one stable embodiment adapter while the runtime underneath
becomes more PromiseGrid-native, avoiding adapter drift and unnecessary
embodiment churn.
Constraints: Do not add a second transport contract, do not bypass the local
HTTP adapter for browser/CLI/Neovim, and keep the tightening step grounded in
real runtime capabilities rather than speculative future transport.
Affects: service/types.go; service/app.go; service/server_test.go; README.md;
docs/promisegrid-implementation-claims.md; docs/architecture.md;
docs/practical-implementation.md

## Goal

Decide whether later `ex5` embodiments should keep routing strictly through
the local HTTP adapter or expose a more direct grid-native runtime contract.

## Why this exists

The current embodiment story is now honest and stable, but it is still
explicitly HTTP-adapter-first. If `ex5` is to become more fully on-grid, that
boundary may need a later tightening pass.

## Tasks

- [x] lurog.1 Run the required TE for later embodiment-contract tightening.
- [x] lurog.2 Lock whether browser, CLI, and Neovim stay adapter-first or gain
  a more direct runtime-facing contract.
- [x] lurog.3 Implement the next concrete tightening step, if any.
- [x] lurog.4 Update docs, claims, and operator guidance to match the chosen
  embodiment contract.

## Status

- closed
- created from the post-109 PromiseGrid review
