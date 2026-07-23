# TODO fovek - ex5 browser host readiness proof

## Decision Intent Log

ID: DI-fovek
Date: 2026-07-22 19:58:07 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a follow-on pass to make browser embodiment readiness prove native-host availability instead of only proving the page bridge is installed.
Intent: Tighten the shipped Chrome/Chromium embodiment so startup truth matches the actual direct browser contract and operators fail early when the native host is missing or misregistered.
Constraints: Keep the browser fail-closed for unsupported or misconfigured direct-contract environments; do not silently demote back into the older HTTP browser path.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-salov
Date: 2026-07-22 20:00:39 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Prove browser embodiment readiness with a one-shot typed `runtime_ready` operation that crosses the page bridge, extension worker, native host, and local runtime socket before startup marks the browser ready.
Intent: Make browser startup truth match the actual direct browser contract without over-coupling readiness to a full live-draft session.
Constraints: Keep the browser fail-closed on missing or misregistered native-host paths, reuse the existing bridge handshake lane instead of inventing a second startup protocol, and document the readiness rule explicitly.
Affects: `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/chrome-extension/content.js`, `ex5-operational-knowledge-system/service/local_socket.go`, `ex5-operational-knowledge-system/service/browser_host_test.go`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/docs/http-api-guide.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Make the browser embodiment declare itself ready only after it has proved the
extension-to-native-host path is actually reachable.

## Tasks

- [x] fovek.1 Define the readiness probe boundary across page bridge, extension worker, and native host. See `../../docs/thought-experiments/TE-funek-ex5-browser-host-readiness-proof.md`.
- [x] fovek.2 Add coverage for missing-host, misregistered-host, and healthy-host startup outcomes.
- [x] fovek.3 Update runtime and operator docs so the readiness rule is stated explicitly.

## Status

- completed
- browser readiness now proves the native-host path with a typed `runtime_ready` round-trip before startup marks the browser ready
