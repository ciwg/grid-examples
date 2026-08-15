# TODO lavem - ex5 Chromium embodiment bug

## Decision Intent Log

ID: DI-bahak
Date: 2026-08-10 20:05:49 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Retain Chromium-capable implementation code, but advertise, script, document, and verify only Chrome as the supported browser embodiment until Chromium has a separate passing verification path.
Intent: Keep the shipped PromiseGrid embodiment claim evidence-backed without discarding the existing Chromium-capable path or conflating capability with verified support.
Constraints: The Chrome demo must fail closed when Chrome is unavailable; Chromium remains explicitly deferred under TODO 153 and must not be selected as a substitute; browser metadata must name the verified Chrome requirement.
Affects: `ex5-operational-knowledge-system/service/app.go`, `ex5-operational-knowledge-system/service/app_test.go`, `ex5-operational-knowledge-system/service/server_test.go`, `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/scripts/*`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/*`

ID: DI-lavem
Date: 2026-07-23 10:30:11 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track Chromium embodiment failure as a deferred bug and stop spending demo-prep time on it today; use Chrome as the working browser path for the live ex5 demo instead.
Intent: Preserve the evidence that ex5 currently works in Chrome but not in the tested Chromium environment, without blocking today’s demo-prep work.
Constraints: Do not conflate this with the working Chrome path; treat it as a separate browser-embodiment bug around extension/native-host behavior in Chromium.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-vamid
Date: 2026-08-14 18:57:11 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Retain verified Google Chrome as Ex5's only supported browser embodiment and defer Chromium Snap support.
Intent: Keep Ex5's browser claim evidence-backed while avoiding a fragile Snap-specific native-host installer without a verified registration and confinement contract.
Constraints: Do not advertise Chromium support, alter the verified Chrome path, or add an HTTP fallback. Chromium-capable extension code remains unclaimed until a future separately approved support design and passing verification exist.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `TODO/TODO-lavem-ex5-chromium-embodiment-bug.md`, `docs/thought-experiments/TE-tilun-ex5-chromium-snap-support.md`, browser support documentation

## Goal

Capture the current Chromium-specific browser embodiment failure so it can be
troubleshot later without derailing today’s working Chrome demo path.

## Tasks

- [x] lavem.1 Reproduce the Chromium failure cleanly and record the exact browser/build/runtime conditions.
- [x] lavem.2 Determine whether the failure is Chromium native-messaging lookup, Snap confinement, or extension launch behavior.
- [x] lavem.3 Decide whether Chromium support needs a separate setup/verification path from Chrome.

## Status

- closed
- resolved by retaining verified Google Chrome as the supported browser embodiment; Chromium Snap support requires a future separately approved design. Source: `DI-vamid`.
- 2026-08-14 reproduction: Chromium `151.0.7922.108 snap` loaded the shipped
  Ex5 extension (`miagfmaampfgjkojhccdilogehbjijpe`) in a fresh temporary
  profile, but the page handshake failed with `Specified native messaging host
  not found.` against an isolated Ex5 runtime.
- 2026-08-14 classification: registering the generated host manifest in the
  disposable profile, in Chromium Snap's apparent common registry, and beside
  a temporary Snap-local host binary did not change that lookup failure. The
  installed Chromium Snap exposes no native-messaging connector interface.
  The remaining work is therefore a separate Snap-native registration and
  confinement design decision, not an Ex5 extension-load or runtime bug.
