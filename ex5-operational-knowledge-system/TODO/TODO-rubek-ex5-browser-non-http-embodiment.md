# TODO rubek - ex5 browser non-http embodiment

## Decision Intent Log

ID: DI-rubek
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to move the browser onto a direct non-HTTP embodiment contract.
Intent: Improve the biggest remaining embodiment-layer PromiseGrid impurity after the terminal surfaces have already moved onto the local socket contract.
Constraints: Treat this as a broader scope wave; it will likely require new transport and embodiment-boundary decisions rather than a small cleanup.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

ID: DI-punek
Date: 2026-07-22 22:04:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Move the browser onto a Chrome/Chromium Manifest V3 native-messaging embodiment that bridges the existing browser UI into the direct local runtime contract, with no silent fallback to the older HTTP browser path.
Intent: Remove the main remaining embodiment-layer adapter impurity by making browser request/response and live drafting ride over the same direct local contract family already used by terminal embodiments, while stating the Chrome/Chromium requirement honestly and explicitly.
Constraints: Keep the existing `web/app.js` UI surface, reuse typed `operation` plus live message semantics instead of inventing a browser-only semantic family, use a fixed extension identity plus native-host origin, and document the unsupported-browser boundary clearly.
Affects: `ex5-operational-knowledge-system/chrome-extension/*`, `ex5-operational-knowledge-system/cmd/operational-browser-host/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define and stage a browser embodiment contract that no longer relies on the
current local HTTP adapter as its primary runtime surface.

## Tasks

- [x] rubek.1 Define the first direct browser embodiment contract candidate. See `../../docs/thought-experiments/TE-sarek-ex5-browser-non-http-embodiment-first-slice.md`.
- [x] rubek.2 Compare it against the current local HTTP adapter behavior and migration cost. See `../../docs/thought-experiments/TE-sarek-ex5-browser-non-http-embodiment-first-slice.md`.
- [x] rubek.3 Lock the first browser non-HTTP slice and stage implementation.

## Status

- completed
- Chrome/Chromium MV3 native-messaging browser embodiment is now the shipped direct browser contract
- `132A` locked and implemented over the existing browser app
