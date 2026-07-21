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

ID: DI-rafek
Date: 2026-07-20 20:17:05 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Close the presence-lifecycle TODO because the browser and Neovim already render `live`/`stale`/`offline`, the timing policy is now covered by browser tests, and the separate durable collaboration surfaces shipped later as Phase 3 review/history features.
Intent: Bring the TODO queue back in sync with the actual ex3 implementation instead of leaving a completed feature slice open.
Constraints: This closes the lifecycle split as implemented today; it does not claim that all future durable collaboration surfaces are complete beyond the current review/history set.
Affects: `ex3-grid-editor-websocket/web`, `ex3-grid-editor-websocket/nvim`, `ex3-grid-editor-websocket/docs`, `ex3-grid-editor-websocket/TODO`

Goal: Implement the chosen live presence aging policy and keep historical collaboration signals out of the main peer roster.

- [x] mivor.1 Add relay-side or client-side peer freshness tracking so live awareness entries can age from live to stale to offline and then disappear.
- [x] mivor.2 Update the browser UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [x] mivor.3 Update the Neovim UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [x] mivor.4 Add tests covering awareness expiration thresholds and peer removal timing.
- [x] mivor.5 Design separate follow-up surfaces for durable collaboration signals:
  document activity
  comments
  version history
  `last viewed`
  `last edited`

Resolved shape:
- live presence remains in the `Peers` card and Neovim peer roster
- durable collaboration signals now live in the separate review/history surfaces:
  activity, comments, saved versions, recent participants, and local timestamps
