# TODO lavem - ex5 Chromium embodiment bug

## Decision Intent Log

ID: DI-lavem
Date: 2026-07-23 10:30:11 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track Chromium embodiment failure as a deferred bug and stop spending demo-prep time on it today; use Chrome as the working browser path for the live ex5 demo instead.
Intent: Preserve the evidence that ex5 currently works in Chrome but not in the tested Chromium environment, without blocking today’s demo-prep work.
Constraints: Do not conflate this with the working Chrome path; treat it as a separate browser-embodiment bug around extension/native-host behavior in Chromium.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Capture the current Chromium-specific browser embodiment failure so it can be
troubleshot later without derailing today’s working Chrome demo path.

## Tasks

- [ ] lavem.1 Reproduce the Chromium failure cleanly and record the exact browser/build/runtime conditions.
- [ ] lavem.2 Determine whether the failure is Chromium native-messaging lookup, Snap confinement, or extension launch behavior.
- [ ] lavem.3 Decide whether Chromium support needs a separate setup/verification path from Chrome.

## Status

- open
- deferred for later troubleshooting because Chrome is working and demo-prep should stay on the known-good path today
