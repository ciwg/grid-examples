# TODO mivor - grid-editor presence lifecycle

## Decision Intent Log

ID: DI-mivor
Date: 2026-07-13 09:55:00 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Treat the main `Peers` list as live presence only, with testing-friendly aging windows of `0-1 minute` live, `1-5 minutes` stale/dimmed, `5-15 minutes` offline, and `15+ minutes` removed; keep historical collaboration signals such as document activity, comments, version history, `last viewed`, and `last edited` in separate surfaces instead of the live peer roster.
Intent: Preserve the mental model of "who is here now?" while still leaving room for richer historical collaboration features that do not make the live presence UI look broken or sticky during testing.
Constraints: This decision defines the intended UX policy and separation of concerns; it does not yet define the exact storage/query model for durable activity, comments, or version history, which may require later TE/DI work.
Affects: `ex2-grid-editor/protocols/live-awareness.md`, `ex2-grid-editor/docs/grid-editor-ui-example.md`, `ex2-grid-editor/web`, `ex2-grid-editor/nvim`, `ex2-grid-editor/service`

Goal: Implement the chosen live presence aging policy and keep historical collaboration signals out of the main peer roster.

- [ ] mivor.1 Add relay-side or client-side peer freshness tracking so live awareness entries can age from live to stale to offline and then disappear.
- [ ] mivor.2 Update the browser UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [ ] mivor.3 Update the Neovim UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [ ] mivor.4 Add tests covering awareness expiration thresholds and peer removal timing.
- [ ] mivor.5 Design separate follow-up surfaces for durable collaboration signals:
  document activity
  comments
  version history
  `last viewed`
  `last edited`
