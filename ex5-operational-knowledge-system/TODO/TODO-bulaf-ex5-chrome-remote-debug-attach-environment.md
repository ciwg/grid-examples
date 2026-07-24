# TODO bulaf - ex5 Chrome remote-debug attach environment

## Decision Intent Log

ID: DI-bulaf
Date: 2026-07-23 11:21:46 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track the Chrome remote-debug attach failure as a separate environment/tooling problem so it can be debugged explicitly instead of being confused with ex5 browser-flow correctness or PromiseGrid backend behavior.
Intent: Preserve the evidence that the ex5 browser path can work manually while Playwright still cannot attach to the same environment through Chrome remote debugging.
Constraints: Keep this separate from Chromium embodiment work and from product-level browser interaction bugs; focus on why Chrome launched with `--remote-debugging-port=9222` is not yielding a reliable attachable endpoint in this environment.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `/home/jj/lab/cswg/grid-examples-browser-checks/ex5/**`

## Goal

Determine why Chrome remote debugging is not yielding a stable attachable
endpoint for the external ex5 browser-check harness, even when the manual ex5
browser page can be working at the same time.

## Tasks

- [ ] bulaf.1 Record one clean reproduction of the Chrome launch command, observed browser state, `ps` output, and `ss`/`curl` results for `9222`.
- [ ] bulaf.2 Determine whether the failure is launch-environment specific (Wayland, profile chooser, crash/restart behavior, or Chrome session replacement) rather than an ex5 extension/runtime problem.
- [ ] bulaf.3 Prove one attachable Chrome launch recipe that keeps `9222` alive long enough for Playwright to connect.

## Status

- open
- created after repeated local runs showed that manual ex5 browser behavior can work while attach-mode automation still cannot reach a stable Chrome DevTools endpoint
- 2026-07-23 progress:
  - proved that attach-mode can work when Chrome is launched against the copied profile under `/tmp/ex5-demo-browser/chrome-profile-real` and `curl -s http://127.0.0.1:9222/json/version` returns a live DevTools endpoint
  - proved that the remaining failure after attach is no longer Playwright itself: the page handshake in that attached session returns `{"__oks_bridge":true,"direction":"bridge->page","kind":"handshake","ok":false}`
  - therefore the current blocker is the extension/native-host/runtime leg inside the copied-profile Chrome session, not the PromiseGrid backend and not the Playwright CDP attach step
