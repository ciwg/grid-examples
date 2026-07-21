# TODO mitob - ex5 browser participant identity hardening

## Decision Intent Log

ID: DI-mitob
Date: 2026-07-21 13:25:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track browser participant identity hardening as its own ex5 bug fix TODO.
Intent: Prevent ex5 browser startup from failing when local storage or `crypto.randomUUID()` is unavailable in restrictive/private browsing environments.
Constraints: Keep the current local HTTP live-draft model; hardening should stay inside ex5 browser code and tests; update docs in the same pass.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mitob-ex5-browser-participant-identity-hardening.md`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`, `ex5-operational-knowledge-system/README.md`

## Goal

Make browser participant identity allocation resilient when the environment
blocks storage access or does not expose the expected UUID helper.

## Review finding

`getParticipantID()` directly calls `window.localStorage.getItem`,
`crypto.randomUUID()`, and `window.localStorage.setItem` without any fallback.
That can break ex5 browser startup in private or otherwise restricted browser
contexts before the live-draft UI even boots.

## Tasks

- [x] mitob.1 Add safe fallback behavior for participant identity creation when storage access or UUID generation is unavailable.
- [x] mitob.2 Add browser regression coverage and document the restricted-environment behavior.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
