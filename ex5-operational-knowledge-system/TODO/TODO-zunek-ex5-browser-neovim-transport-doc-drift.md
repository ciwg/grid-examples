# TODO zunek - ex5 browser and neovim transport doc drift

## Decision Intent Log

ID: DI-zunek
Date: 2026-07-22 19:58:07 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to remove the remaining high-visibility transport wording drift for the browser and Neovim embodiments.
Intent: Keep the PromiseGrid docs honest now that the shipped browser contract is Chrome/Chromium native messaging and Neovim is socket-first with explicit compatibility mode.
Constraints: Update stale docs rather than deleting explanatory material; preserve DI-backed wording where it is still accurate.
Affects: `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-vubem
Date: 2026-07-22 20:14:41 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Apply the `136B.1` high-visibility transport sweep to the already-approved six doc surfaces only, aligning browser wording to the shipped Chrome/Chromium native-messaging embodiment and Neovim wording to the shipped socket-first contract with explicit compatibility mode.
Intent: Make the main PromiseGrid reader path honest end to end without turning `136` into a full-corpus editorial rewrite.
Constraints: Update stale wording rather than deleting explanatory material; keep the browser fail-closed unsupported-browser requirement explicit; describe Neovim compatibility as explicit opt-in instead of silent websocket demotion.
Affects: `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/docs/architecture.md`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/docs/promisegrid-implementation-claims.md`, `docs/thought-experiments/TE-lurem-ex5-browser-neovim-transport-doc-drift.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Bring the remaining long-form ex5 docs into line with the shipped browser and
Neovim transport model.

## Tasks

- [x] zunek.1 Sweep the README and long-form guides for stale websocket-first or HTTP-primary wording. See `../../docs/thought-experiments/TE-lurem-ex5-browser-neovim-transport-doc-drift.md`.
- [x] zunek.2 Update the affected passages to describe the shipped embodiment contracts directly.
- [x] zunek.3 Recheck the main PromiseGrid doc surfaces for transport wording consistency.

## Status

- completed
- the high-visibility transport sweep now describes the shipped browser native-messaging contract and Neovim explicit compatibility mode consistently across the approved doc set
