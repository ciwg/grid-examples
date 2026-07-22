# TODO ruvot - ex5 browser error handling

## Decision Intent Log

ID: DI-ruvot
Date: 2026-07-21 16:11:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track browser form/search failure handling as its own ex5 review TODO.
Intent: Keep common operator mistakes and server validation failures inside the browser UX instead of leaking them as unhandled async failures.
Constraints: Stay inside `ex5`; preserve the current browser surfaces and forms; improve error handling without changing the broader workflow model.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-ruvot-ex5-browser-error-handling.md`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/features-guide.md`

## Goal

Make the browser handle ordinary request failures cleanly with in-app feedback
instead of depending on unhandled promise behavior.

## Review finding

The browser form submit handlers and search path await network calls directly
with no local `try/catch`, and search does not check `response.ok` before
parsing JSON. Server validation failures can therefore bypass the intended toast
feedback path.

## Tasks

- [x] ruvot.1 Wrap browser form submit flows in consistent request-failure handling.
- [x] ruvot.2 Make browser search treat non-`200` responses as handled errors instead of blindly parsing JSON.
- [x] ruvot.3 Add browser regression coverage for representative failure cases and update docs if needed.

## Status

- done
- derived from the 2026-07-21 deep ex5 review
