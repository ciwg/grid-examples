# TODO lanis - ex5 browser smoke doc honesty

## Decision Intent Log

ID: DI-lanis
Date: 2026-07-21 16:14:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the browser smoke coverage overstatement as its own ex5 review TODO.
Intent: Keep ex5 verification claims honest about what the headless browser suite really proves and what it still stubs out.
Constraints: Stay inside `ex5`; either tighten the docs to match the current stubbed-browser suite or expand the suite to run against more of the real service stack.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lanis-ex5-browser-smoke-doc-honesty.md`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`, `ex5-operational-knowledge-system/docs/practical-implementation.md`

## Goal

Align the ex5 verification story with what the browser smoke tests actually
exercise today.

## Review finding

The docs imply the browser smoke runs the real UI against a live test server in
a way that sounds close to full service integration, but the current suite
mostly builds stub muxes that handcraft `/api/*` responses. That is still
useful, but it is not the same proof level.

## Tasks

- [x] lanis.1 Decide whether to narrow the docs, strengthen the browser smoke suite, or both.
- [x] lanis.2 Make the docs describe the actual verification level accurately.
- [x] lanis.3 If the suite stays stub-based, say so explicitly; if it is strengthened, add the missing integration coverage.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
