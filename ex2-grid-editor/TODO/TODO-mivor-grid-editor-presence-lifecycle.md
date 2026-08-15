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

### DI-dizut

- ID: DI-dizut
- Date: 2026-08-14
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Keep `live-awareness` aging embodiment-local: browser and Neovim refresh after accepted awareness events and again at their next exact lifecycle boundary, using `schedulePresenceRefresh`, `cancelPresenceRefresh`, and `presenceRefreshTimer`. Use injected deterministic scheduling tests that prove each display transition and removal.
- Intent: Preserve live presence as each embodiment's interpretation of its latest accepted local observation rather than grant any relay membership authority or create a new durable protocol event.
- Constraints: Do not add relay expiry, peer-removal messages, a new pCID, durable membership state, or historical collaboration records. The existing normal and demo windows remain unchanged. Timers are local presentation mechanics and must be cancelled when an embodiment disconnects or changes documents.
- Affects: `ex2-grid-editor/{web/src/main.js,nvim/lua/grid_editor/init.lua,web/src/*_test.mjs,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

### DI-dazin

- ID: DI-dazin
- Date: 2026-08-14
- Status: active
- Author: jj@thesalleys.com (JJ)
- Decision: Put Ex2's pure browser-local presence lifecycle classifier and next-boundary scheduler in `web/src/presence.js`, with deterministic clock and scheduler injection tested in `web/src/presence.test.mjs`; keep `web/src/main.js` as the UI integration point.
- Intent: Make the local presentation timer independently testable without importing or simulating the whole browser application, while retaining the DI-dizut rule that relays have no expiry or membership authority.
- Constraints: The module does not send traffic, persist data, alter the `live-awareness` pCID contract, or define durable activity. It receives only already accepted `last_seen_at` observations and reports local display timing.
- Affects: `ex2-grid-editor/{web/src/presence.js,web/src/presence.test.mjs,web/src/main.js,web/app.js,docs/testing.md,TODO}`, `TODO/handle-namespace.tsv`

Goal: Implement the chosen live presence aging policy and keep historical collaboration signals out of the main peer roster.

- [x] mivor.1 Add relay-side or client-side peer freshness tracking so live awareness entries can age from live to stale to offline and then disappear.
- [x] mivor.2 Update the browser UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [x] mivor.3 Update the Neovim UI to render `live`, `stale`, and `offline` peer states distinctly before removal.
- [x] mivor.4 Add tests covering awareness expiration thresholds and peer removal timing.
- [ ] mivor.5 Design separate follow-up surfaces for durable collaboration signals:
  document activity
  comments
  version history
  `last viewed`
  `last edited`
