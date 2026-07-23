# TODO jovek - ex5 implementation promise publication alignment

## Decision Intent Log

ID: DI-jovek
Date: 2026-07-22 21:58:33 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a PromiseGrid dev-guide alignment pass for ex5's published implementation promise claims so the B-side `CHANGELOG.md` remains current, auditable, and honest about which shipped components implement which protocol surfaces.
Intent: Close the gap where ex5's runtime and embodiment layers are strongly PromiseGrid-aligned but the published implementation-promise surface still reflects an older transport era and under-describes component boundaries.
Constraints: Prefer one honest publication surface over a fragmented documentation scheme unless a split is clearly justified; keep the fix scoped to implementation-promise publication, not a broader architecture rewrite.
Affects: `docs/thought-experiments/*`, `ex5-operational-knowledge-system/CHANGELOG.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-nubor
Date: 2026-07-22 22:00:27 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Lock `147B` and keep `CHANGELOG.md` as the single implementation-promise publication surface, but expand it into a component-aware claim document that publishes both the eight frozen family claims and the shipped component-level implementation/delegation boundaries.
Intent: Meet the PromiseGrid dev guide's two strongest app-developer publication requirements at once: exact auditable family claims by doc-CID, and explicit component-level honesty about which shipped pieces implement which parts of the contract.
Constraints: Do not fragment the publication surface into multiple files in this pass; keep summary docs pointing back to `CHANGELOG.md` as the canonical B-side implementation-promise source.
Affects: `../../docs/thought-experiments/TE-satek-ex5-implementation-promise-publication-boundary.md`, `ex5-operational-knowledge-system/CHANGELOG.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`
Supersedes: DI-jovek

## Goal

Bring ex5's published implementation promise claims into line with the
PromiseGrid dev guide's app-developer publication discipline.

## Tasks

- [x] jovek.1 Compare the current ex5 publication surface against the dev guide's implementation-promise requirements. See `../../docs/thought-experiments/TE-satek-ex5-implementation-promise-publication-boundary.md`.
- [x] jovek.2 Lock the smallest honest publication boundary for family claims versus component claims. Locked to `147B`: one component-aware `CHANGELOG.md`.
- [x] jovek.3 Update the publication surface and aligned summary docs so they match the shipped embodiment and component reality.

## Status

- completed
- created from the PromiseGrid dev-guide review finding that ex5's implementation-promise publication layer is behind the shipped runtime and embodiment boundary
- TE complete: `TE-satek` recommends a single component-aware `CHANGELOG.md` as the most PromiseGrid-aligned surviving choice
- locked to `147B`: `CHANGELOG.md` remains the single canonical publication surface, now widened to cover both frozen families and component-level implementation/delegation truth
