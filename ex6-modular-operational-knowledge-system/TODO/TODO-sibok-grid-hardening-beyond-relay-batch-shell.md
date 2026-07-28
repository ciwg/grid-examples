# TODO sibok - grid hardening beyond relay batch shell

## Decision Intent Log

ID: DI-sibok
Date: 2026-07-28 18:20:00
Status: active
Decision: Harden the current relay path by deduplicating exact-byte records during append/import and validating relay batch metadata before accepting a batch.
Intent: Keep repeated relay exchange from inflating append-only history with identical records while preserving unknown-family exact-byte carriage and making malformed batches fail early.
Constraints: Keep the current batch wire shape; do not invent a speculative final PromiseGrid transport; dedupe exact bytes only, not semantic near-matches.
Affects: `store/history.go`, `kernel/runtime.go`, `grid/batch.go`, runtime tests, and ex6 current-state docs.

## Goal

Make the current relay shell safer and less noisy under repeated exchange.

## Scope

- exact-byte record dedupe
- relay batch metadata validation
- repeated import idempotence tests
- current-state hardening note

## Why

Without this pass, repeated imports can append the same raw record forever and
malformed relay batches are not rejected early enough.
