# TODO murev - ex5 embodiment transport metadata unification

## Decision Intent Log

ID: DI-vurak
Date: 2026-07-22 18:16:49 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Replace the current flat embodiment transport fields in `/api/meta` with one top-level `embodiments` table keyed by `browser`, `cli`, and `neovim`.
Intent: Make the embodiment contract explicit in one place and avoid carrying two metadata vocabularies longer than necessary.
Constraints: Do not change underlying transport behavior in this TODO; keep the change limited to metadata shape, tests, and docs.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `docs/thought-experiments/TE-zumek-ex5-embodiment-transport-metadata-shape.md`

ID: DI-zumek
Date: 2026-07-22 18:16:49 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Narrow the embodiment transport metadata decision to either a staged structured embodiment table or a full immediate replacement of the current flat fields.
Intent: Make `129` resolve around one real PromiseGrid contract question instead of accreting more ad hoc top-level metadata fields.
Constraints: Do not change underlying transport behavior in this TODO; this is a metadata-contract refinement only.
Affects: `docs/thought-experiments/TE-zumek-ex5-embodiment-transport-metadata-shape.md`, `ex5-operational-knowledge-system/TODO/TODO-murev-ex5-embodiment-transport-metadata-unification.md`

ID: DI-murev
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to make embodiment transport metadata more explicit and uniform across browser, CLI, and Neovim.
Intent: Improve PromiseGrid alignment in the safest way first by making primary vs compatibility transport semantics more machine-readable before changing larger behavior surfaces.
Constraints: Future-scope only; do not reopen the already-locked browser/CLI/Neovim transport choices in this tracking pass.
Affects: `ex5-operational-knowledge-system/service/types.go`, `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Make the runtime publish a clearer embodiment transport table so primary,
fallback, and compatibility transport semantics are explicit for each shipped
embodiment.

## Tasks

- [x] murev.1 Define the target embodiment transport metadata shape for browser, CLI, and Neovim.
- [x] murev.2 Implement the refined metadata contract and align tests/docs.
- [x] murev.3 Confirm the refined metadata matches the actual shipped transport behavior.

## Status

- closed
- resolved by replacing the flat embodiment transport fields with a structured `embodiments` metadata table
