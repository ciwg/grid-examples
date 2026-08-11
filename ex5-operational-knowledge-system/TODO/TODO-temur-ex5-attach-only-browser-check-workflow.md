# TODO temur - ex5 attach-only browser check workflow

## Decision Intent Log

ID: DI-bilim
Date: 2026-08-10 21:10:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: The demo launcher owns one Chrome session through a DevTools pipe, loads the unpacked Ex5 extension into that live session, and simultaneously exposes localhost DevTools for the attach-only interaction suite.
Intent: Chrome 146 disables unpacked developer extensions after restart; keeping the sanctioned pipe owner alive is required for a real Chrome session while still allowing the independent suite to attach through the documented port.
Constraints: The pipe owner launches only the dedicated `/tmp/ex5-demo-browser/chrome-profile`, loads only the shipped Ex5 extension and its manifest-derived ID, and is the sole PID recorded and terminated by the demo launcher.
Affects: `ex5-operational-knowledge-system/scripts/install-demo-chrome-extension.mjs`, `ex5-operational-knowledge-system/scripts/launch-demo-browser.sh`, `ex5-operational-knowledge-system/scripts/setup-demo-browser.sh`, Chrome attach-only verification
Supersedes: DI-zuvus

ID: DI-lahoz
Date: 2026-08-10 21:05:00 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Use `miagfmaampfgjkojhccdilogehbjijpe` as the Ex5 extension identity, derived from the shipped manifest key, and verify that identity during setup and native-host authorization.
Intent: Avoid Chrome's unrelated built-in `fignfifoniblkonapihmkfakmlgkbkcf` Google Network Speech component; only the manifest-derived identity can establish the Ex5 bridge.
Constraints: Preserve the shipped manifest key and authorize exactly one Chrome extension origin; the installer must reject any returned ID other than the derived Ex5 ID.
Affects: `ex5-operational-knowledge-system/chrome-extension/native-host/operational_browser_host.json`, `ex5-operational-knowledge-system/chrome-extension/assets_test.go`, `ex5-operational-knowledge-system/scripts/install-demo-chrome-extension.mjs`, `ex5-operational-knowledge-system/scripts/setup-demo-browser.sh`, `ex5-operational-knowledge-system/scripts/verify-demo-browser.sh`
Supersedes: DI-vuduz, DI-zuvus

ID: DI-zuvus
Date: 2026-08-10 21:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Install the shipped unpacked extension into the dedicated Ex5 Chrome profile through Chrome's `Extensions.loadUnpacked` DevTools-pipe API during setup, then use the normal remote-debug attach path for the demo and interaction suite.
Intent: Chrome 146 rejects command-line unpacked-extension loading; its supported pipe API creates the required profile state without claiming Chromium behavior or weakening the attach-only evidence path.
Constraints: The installer is limited to `/tmp/ex5-demo-browser/chrome-profile` and the shipped `chrome-extension` directory, verifies the manifest-derived Ex5 identity, and closes only its installer-owned Chrome process before the regular launcher starts.
Affects: `ex5-operational-knowledge-system/scripts/install-demo-chrome-extension.mjs`, `ex5-operational-knowledge-system/scripts/setup-demo-browser.sh`, `ex5-operational-knowledge-system/scripts/launch-demo-browser.sh`
Supersedes: DI-fofup

ID: DI-fofup
Date: 2026-08-10 20:57:22 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Launch the verified Google Chrome demo with Chrome's supported unpacked-extension flags, `--enable-unsafe-extension-debugging` and `--load-extension`, without `--disable-extensions-except`.
Intent: Remove the restriction flag that Chrome 146 explicitly rejects so the shipped Ex5 extension can load and supply the real page bridge.
Constraints: Keep the dedicated disposable profile, Chrome-only contract, and attach-only DevTools port; do not use Chromium as a fallback.
Affects: `ex5-operational-knowledge-system/scripts/launch-demo-browser.sh`, Chrome attach-only verification

ID: DI-vuduz
Date: 2026-08-10 20:24:45 -0700
Author: jj@thesalleys.com (JJ)
Status: superseded
Decision: Authorize the native host for the Chrome extension ID actually loaded from the shipped unpacked manifest: `fignfifoniblkonapihmkfakmlgkbkcf`.
Intent: Make the verified Chrome browser embodiment carry its handshake into the native host instead of accepting a stale extension origin that Chrome does not use.
Constraints: Keep the authorization limited to the one shipped extension origin and verify that exact origin before a demo is claimed ready.
Affects: `ex5-operational-knowledge-system/chrome-extension/native-host/operational_browser_host.json`, `ex5-operational-knowledge-system/chrome-extension/assets_test.go`, `ex5-operational-knowledge-system/scripts/verify-demo-browser.sh`

ID: DI-danir
Date: 2026-08-10 20:22:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Make the disposable Ex5 demo launcher record and replace only its own prior Chrome PID before creating the next attachable Chrome session.
Intent: Keep attach-only browser evidence reproducible by preventing the launcher from silently reusing a stale demo-profile Chrome session, without targeting ordinary user browser processes.
Constraints: The launcher may signal only the PID stored in `/tmp/ex5-demo-browser/chrome.pid`; it must confirm termination before launch and fail closed if that exact process remains alive.
Affects: `ex5-operational-knowledge-system/scripts/launch-demo-browser.sh`, `ex5-operational-knowledge-system/docs/testing.md`, `/tmp/ex5-demo-browser/chrome.pid`

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
