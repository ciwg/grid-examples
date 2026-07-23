# TODO luneg - ex5 features guide PromiseGrid drift

## Decision Intent Log

ID: DI-lusen
Date: 2026-07-22 17:35:03 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Fully align `docs/features-guide.md` with the shipped signed-envelope, relay, and current Neovim embodiment scope now.
Intent: Remove the last remaining top-level summary drift so the feature guide matches the README, product overview, relay guide, and current runtime behavior.
Constraints: Docs-only scope; do not change product boundaries or runtime semantics.
Affects: `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `docs/thought-experiments/TE-lunok-ex5-features-guide-promisegrid-drift.md`
Supersedes: DI-luneg

ID: DI-luneg
Date: 2026-07-22 17:23:59 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the remaining feature-guide wording that still understates the shipped signed-envelope, relay, and Neovim embodiment scope.
Intent: Finish the PromiseGrid doc-alignment work by bringing the feature guide into line with the README, claims doc, relay guide, and current runtime behavior.
Constraints: This is a doc-honesty pass only; it should not expand product scope or reopen settled runtime decisions.
Affects: `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/product-overview.md`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-luneg-ex5-features-guide-promisegrid-drift.md`

## Goal

Remove the remaining feature-guide statements that still describe ex5 as only
the browser/CLI local runtime layer or as pre-signed-envelope / pre-relay.

## Tasks

- [x] luneg.1 Find the remaining stale scope statements in `docs/features-guide.md`.
- [x] luneg.2 Rewrite them so the shipped signed-envelope, relay, and current Neovim embodiment scope are described consistently.
- [x] luneg.3 Re-run a PromiseGrid wording sweep on the summary doc surfaces after the feature-guide update.

## Status

- closed
- resolved by aligning the feature guide with the shipped PromiseGrid scope
