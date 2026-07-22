# TODO tasuk - correct downstream ex5 PromiseGrid summary doc drift

## Decision Intent Log

ID: DI-tasuk
Date: 2026-07-22 14:18:44 -0700
Status: active
Decision: Correct the remaining downstream docs so they stop describing
already-shipped signed-envelope and peer-exchange behavior as absent.
Intent: Keep the secondary ex5 doc layer aligned with the shipped runtime and
the main PromiseGrid claims docs.
Constraints: Preserve the honest distinction between the shipped local adapter
surface and unshipped websocket/relay-network follow-on work.
Affects: docs/product-overview.md; docs/features-guide.md; TODO/TODO.md

## Goal

Correct the secondary ex5 docs that still describe already-shipped PromiseGrid
behavior as not yet shipped.

## Why this exists

The main claims docs are now fairly accurate, but at least
`docs/product-overview.md` and `docs/features-guide.md` still describe
relay-visible peer exchange and signed-envelope-on-wire behavior as absent,
which no longer matches the shipped runtime.

## Tasks

- [x] tasuk.1 Update `docs/product-overview.md` so its PromiseGrid boundary
  section matches the shipped peer-exchange/runtime state.
- [x] tasuk.2 Update `docs/features-guide.md` so its “not yet included” section
  stops listing already-shipped signed-envelope and peer-exchange behavior.
- [x] tasuk.3 Sweep adjacent secondary docs for the same stale phrasing and
  correct any remaining drift found in that pass.

## Status

- closed
- created from the post-112 PromiseGrid alignment pass
