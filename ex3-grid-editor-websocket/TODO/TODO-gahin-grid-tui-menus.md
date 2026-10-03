# TODO gahin - Grid TUI menus

## Decision Intent Log

ID: DI-mutoh
Date: 2026-10-03 00:00:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Add terminal-native menus to Ex3 Grid TUI before the remaining queued Charm enhancements.
Intent: Give the collaboration workspace clear, discoverable homes for document, view, collaboration, relay, and help actions before adding further terminal features.
Constraints: Use the existing Charm/Bubble Tea surface; no browser imitation and no changes to the signed relay, Automerge, or WebSocket protocol path.
Affects: `ex3-grid-editor-websocket/cmd/grid-tui`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/TODO`.

## Goal

Add keyboard- and mouse-accessible terminal menus to Grid TUI.

## Tasks

- [x] gahin.1 Add the top-level Document, View, Collaborate, Relay, and Help menus.
- [x] gahin.2 Connect document preview, profile/peer controls, relay inspection, and shortcut help to their menus.
- [x] gahin.3 Make the Glow reader available from Help after TODO budad ships.
- [x] gahin.4 Add deterministic menu navigation and activation tests.
