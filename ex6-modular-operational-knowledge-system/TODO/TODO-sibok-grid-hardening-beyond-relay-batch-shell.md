# TODO sibok - grid hardening beyond relay batch shell

## Decision Intent Log

ID: DI-sibok
Date: 2026-07-28 18:20:00
Status: active
Decision: Harden the current relay path by deduplicating exact-byte records during append/import and validating relay batch metadata before accepting a batch.
Intent: Keep repeated relay exchange from inflating append-only history with identical records while preserving unknown-family exact-byte carriage and making malformed batches fail early.
Constraints: Keep the current batch wire shape; do not invent a speculative final PromiseGrid transport; dedupe exact bytes only, not semantic near-matches.
Affects: `store/history.go`, `kernel/runtime.go`, `grid/batch.go`, runtime tests, and ex6 current-state docs.

ID: DI-nasek
Date: 2026-07-28 18:35:00
Status: active
Decision: Add a minimal live peer exchange layer with HTTP `relay serve`, `relay pull`, and `relay push` commands that carry the existing validated relay batch format.
Intent: Move ex6 beyond file-only batch exchange while keeping the current protocol narrow, testable, and compatible with the hardened batch layer.
Constraints: Reuse the current batch wire shape; no speculative peer discovery or trust protocol yet; keep endpoints explicit and local-runtime owned.
Affects: `cmd/moks`, relay CLI tests, and ex6 current-state docs.

ID: DI-rupem
Date: 2026-07-28 18:50:00
Status: active
Decision: Add runtime-owned peer allow rules, stable local peer identity, and batch-to-peer identity matching for live relay pull/push/import.
Intent: Make multi-peer exchange safer by requiring explicit peer registration and by rejecting live imports whose declared batch identity does not match the allowed peer entry.
Constraints: Keep trust policy explicit and local; no automatic peer discovery; no cryptographic signing layer yet.
Affects: `grid/peers.go`, `kernel/runtime.go`, `cmd/moks`, relay tests, and ex6 current-state docs.

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
