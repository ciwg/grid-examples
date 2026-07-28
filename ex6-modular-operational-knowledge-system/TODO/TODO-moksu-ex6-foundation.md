# TODO moksu - ex6 foundation

## Decision Intent Log

ID: DI-moksu
Date: 2026-07-28 00:00:00
Status: active
Decision: Build ex6 as a basket-first PromiseGrid-native runtime with manifest plus self-check packages, runtime-mediated package communication, and unknown-family exact-byte store-and-relay.
Intent: Keep ex6 as the durable basket that owns grid-facing coordination while letting built-in and installed eggs extend it without direct package-to-package coupling.
Constraints: All artifacts stay inside `ex6-modular-operational-knowledge-system/`; ex5 is reference material only; browser is out of scope; installed packages use executables rather than Go plugin ABI.
Affects: `cmd/moks`, `builtin`, `grid`, `kernel`, `packages`, `records`, `store`, tests, and ex6-local docs.

## Notes

- The repo-local `mint-handle` helper was not present in the checked workspace during implementation, so this ex6-local TODO uses a manually chosen proquint-style handle to preserve the required DI linkage inside the allowed ex6 boundary.
