# TODO temur - ex5 attach-only browser check workflow

## Decision Intent Log

ID: DI-temur
Date: 2026-07-23 11:21:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the follow-on work to make the external ex5 browser checks attach to an already-running, known-good browser session instead of launching their own browser.
Intent: Keep the automation strategy aligned with the browser configuration that already works for live ex5 manual demos.
Constraints: Do not regress into Playwright-launched synthetic browser sessions; prefer attach-only checks against a preverified Chrome session once the remote-debug attach environment is stable.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `/home/jj/lab/cswg/grid-examples-browser-checks/ex5/**`

## Goal

Finish the external ex5 browser-check harness so it attaches to a known-good
browser session and can then identify real UI bugs without reintroducing
browser-launch instability.

## Tasks

- [ ] temur.1 Keep the external `ex5` Playwright harness attach-only and document the required prelaunched browser contract clearly.
- [ ] temur.2 Once Chrome attach is stable, re-run the first hotspot demo-path check and capture the first real UI assertion failure instead of an attach failure.
- [ ] temur.3 Expand the attach-only checks across the draft, hotspot, and `Current Record` demo path so real browser regressions are identified outside the main repo.

## Status

- open
- blocked by `TODO/TODO-bulaf-ex5-chrome-remote-debug-attach-environment.md`
