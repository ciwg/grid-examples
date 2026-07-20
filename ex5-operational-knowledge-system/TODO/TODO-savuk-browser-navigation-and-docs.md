# TODO savuk - browser navigation and docs

## Decision Intent Log

ID: DI-vopuk
Date: 2026-07-20 11:34:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add record inspector, contextual navigation, richer search presentation, and explicit open-items documentation to the current `ex5` browser workflow without changing the underlying ex5 runtime contract.
Intent: Make the current browser implementation easier to use and explain while the larger collaboration-model decisions remain open.
Constraints: Keep the current HTTP/draft model; do not introduce websocket relay behavior in this slice; keep the work focused on safe polish and coverage.
Affects: `web/index.html`, `web/app.js`, `web/style.css`, `docs/**`, `README.md`, `web/assets_test.go`

## Goal

Improve the current `ex5` browser workflow and docs without crossing into the
still-open collaboration architecture questions.

## Tasks

- [x] savuk.1 Add record inspector and contextual navigation.
- [x] savuk.2 Improve search and live-draft conflict presentation.
- [x] savuk.3 Add README open items and align docs.
- [x] savuk.4 Extend browser-asset regression coverage.
