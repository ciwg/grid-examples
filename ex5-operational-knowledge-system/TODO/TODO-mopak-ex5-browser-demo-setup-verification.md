# TODO mopak - ex5 browser demo setup verification

## Decision Intent Log

ID: DI-zuvor
Date: 2026-07-23 08:23:51 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a browser-demo verification wave so the ex5 browser demo can be run from one sheet without hidden prerequisite guessing.
Intent: Convert the current browser demo from an assumed setup into an explicit, fail-closed setup-and-verify path that matches the one-sheet demo promise.
Constraints: Preserve the browser as the primary demo embodiment; fail closed if the direct browser path is not ready; do not silently demote the demo path to CLI-first; keep the setup boundary honest about extension and native-host requirements.
Affects: `docs/thought-experiments/*`, `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-mopak-ex5-browser-demo-setup-verification.md`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/scripts/*`, `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/cmd/operational-browser-host/*`

ID: DI-dabek
Date: 2026-07-23 08:23:51 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the browser-demo fix as one setup script, one launch script, and one verify script, all anchored to a disposable `/tmp/ex5-demo-browser/` state root and an auto-registered Chrome native-host manifest for `operational_browser_host`.
Intent: Make the browser demo actually runnable from one sheet by moving the hidden extension/native-host/runtime setup into one explicit preflight path that fails closed before the live demo starts.
Constraints: Keep the browser as the primary demo embodiment; keep the setup disposable under `/tmp`; register the native host automatically in `/home/jj/.config/google-chrome/NativeMessagingHosts/operational_browser_host.json`; do not claim readiness until the host can forward `runtime_ready` to the real runtime socket.
Affects: `../../docs/thought-experiments/TE-tevok-ex5-browser-demo-setup-verification-boundary.md`, `ex5-operational-knowledge-system/scripts/setup-demo-browser.sh`, `ex5-operational-knowledge-system/scripts/launch-demo-browser.sh`, `ex5-operational-knowledge-system/scripts/verify-demo-browser.sh`, `ex5-operational-knowledge-system/docs/user-guide.md`, `ex5-operational-knowledge-system/TODO/TODO-mopak-ex5-browser-demo-setup-verification.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Make the browser demo runnable from one sheet through explicit setup,
verification, and then exact demo steps.

## Tasks

- [x] mopak.1 Define the browser-demo setup and verification boundary. See `../../docs/thought-experiments/TE-tevok-ex5-browser-demo-setup-verification-boundary.md`.
- [x] mopak.2 Lock the exact setup, launch, and verification path.
- [x] mopak.3 Implement the setup/verify flow and align the one-sheet demo path to it.

## Status

- completed
- created from the failure where the browser demo sheet was not runnable as written
- TE complete: `TE-tevok` recommends a dedicated browser-demo setup and verification path instead of assumed browser readiness
- locked to one explicit setup, launch, and verify path over `/tmp/ex5-demo-browser/`, with auto-registration of `operational_browser_host`
- implemented as `setup-demo-browser.sh`, `launch-demo-browser.sh`, and `verify-demo-browser.sh`, with guide updates that now fail closed on wrong-runtime or missing-readiness states
