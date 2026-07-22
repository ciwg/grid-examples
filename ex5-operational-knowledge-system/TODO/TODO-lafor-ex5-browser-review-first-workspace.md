# TODO lafor - ex5 browser review-first workspace

## Decision Intent Log

ID: DI-lafor
Date: 2026-07-22 03:35:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Implement the ex5 browser redesign as a review-first single-page workspace while preserving all current browser functionality.
Intent: Make the browser usable for operators by reorganizing the page around review, authoring, and action-taking without deleting any current form, drilldown, or debug capability.
Constraints: Keep the existing local HTTP API and single-page browser shell, preserve every current browser workflow, keep raw JSON and debug payloads reachable behind disclosure, and prefer context-aware defaults over removing manual override paths.
Affects: `ex5-operational-knowledge-system/TODO/TODO.md`, `ex5-operational-knowledge-system/TODO/TODO-lafor-ex5-browser-review-first-workspace.md`, `ex5-operational-knowledge-system/web/index.html`, `ex5-operational-knowledge-system/web/style.css`, `ex5-operational-knowledge-system/web/app.js`, `ex5-operational-knowledge-system/web/browser_smoke_test.go`, `ex5-operational-knowledge-system/web/assets_test.go`, `ex5-operational-knowledge-system/README.md`, `ex5-operational-knowledge-system/docs/browser-ui-guide.md`, `ex5-operational-knowledge-system/docs/user-guide.md`

## Goal

Preserve the full browser feature set while reorganizing the page into a
review-first operational workspace with clearer workflow hierarchy and easier
high-frequency actions.

## Tasks

- [x] lafor.1 Reorganize the browser shell into review, author, operate, create, and browse zones.
- [x] lafor.2 Add scoped browser status/error feedback instead of leaking failures into unrelated panels.
- [x] lafor.3 Add context-aware form defaults and selection helpers while preserving manual override paths.
- [x] lafor.4 Add direct hotspot/search/inspector handoffs and hide debug payloads behind disclosure.
- [x] lafor.5 Update browser smoke coverage and user-facing docs to match the redesigned shell.
