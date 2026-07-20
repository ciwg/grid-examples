# TODO savuk - ex5 workflow navigation polish

## Decision Intent Log

ID: DI-savuk
Date: 2026-07-20 11:33:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Improve `ex5` with safe browser and documentation polish that does not depend on unresolved product choices about websocket collaboration, collaborative editing scope, or a future Neovim embodiment.
Intent: Make the current standalone `ex5` implementation easier to inspect and demo by adding better in-browser detail views, contextual navigation, clearer live-draft conflict handling, and explicit open-items documentation.
Constraints: Do not assume the unresolved collaboration decisions are settled; do not add `ex3` websocket transport yet; keep `ex5` standalone.
Affects: `ex5-operational-knowledge-system/web/**`, `ex5-operational-knowledge-system/docs/**`, `ex5-operational-knowledge-system/README.md`, `TODO/TODO.md`, `TODO/TODO-savuk-ex5-workflow-navigation-polish.md`

## Goal

Polish the current `ex5` workflow slice with better browser inspection and
navigation plus clearer documentation of what remains open.

## Tasks

- [x] savuk.1 Add browser detail and contextual navigation views for places, resources, responsibilities, items, runs, and search results.
- [x] savuk.2 Improve the live-draft conflict and selection UX without changing the underlying transport model.
- [x] savuk.3 Document the unresolved `ex5` collaboration decisions and the current browser workflow honestly.
- [x] savuk.4 Add tests for the new browser structure and any changed behavior.
