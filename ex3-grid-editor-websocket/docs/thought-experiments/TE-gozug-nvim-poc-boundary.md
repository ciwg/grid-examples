# Ex3 Neovim POC boundary

TE ID: TE-gozug
## Status
decided

## Decision under test

What native-Neovim proof of concept should Ex3 add so that its existing
collaborative editor is more legible and impressive without changing the
browser, duplicating the editor in a standalone terminal app, or changing the
relay/protocol boundary?

## Assumptions and scope

- This is TODO-numop, a proof of concept limited to
  `ex3-grid-editor-websocket/`.
- The browser embodiment is unchanged.
- The existing Neovim plugin already owns local display state for document ID,
  relay connectivity, presentation identity, and peer awareness; it already
  renders remote cursors and selections.
- The existing `grid-nvim-sidecar` command and its newline-delimited local
  stdio protocol remain the embodiment-local bridge. TE-zorud locked that
  hybrid boundary.
- Names and colors are presentation hints, not signing-key continuity or
  authority. The current PromiseGrid Development Guide treats the app surface
  as provisional and requires protocol/implementation boundaries to remain
  explicit.
- No new pCID, envelope field, browser route, authorization rule, or durable
  record is in scope.

## Alternatives

1. **Local session dashboard.** Add a native Neovim floating-window dashboard
   and a bounded local activity list, derived only from existing plugin state
   and sidecar events.
2. **Raw transport inspector.** Extend the sidecar protocol to emit decoded
   WebSocket/relay traffic and build a Neovim inspector for it.
3. **Standalone Charm TUI.** Build a Bubble Tea terminal program beside the
   Neovim plugin.
4. **Browser enhancement.** Improve Ex3's browser flow/trace card instead of
   the plugin.

## Scenario analysis

### Alice opens a shared document while Bob edits in a browser

Alternative 1 gives Alice one native view showing document, relay status,
participant label, current peers, and join/leave/connect changes. Alice keeps
editing in the same Neovim buffer and still sees Bob's cursor and selection.
It makes the existing collaboration visible without introducing a second
editor.

Alternative 2 gives a deeper demonstration of the transport, but raw decoded
payloads compete with the editing task and require a new sidecar contract.

Alternative 3 gives a polished terminal program but duplicates the plugin's
core purpose and asks Alice to leave Neovim.

Alternative 4 may be valuable independently, but does not satisfy the stated
Neovim-only POC boundary.

### The relay disconnects or the sidecar emits malformed output

Alternative 1 can display the current disconnected state and record a local
activity entry while retaining existing error notifications. It neither retries
silently nor changes admission.

Alternative 2 creates more data to parse and more opportunities for a tracing
surface to misrepresent malformed or incomplete traffic as a protocol fact.

Alternative 3 has the same connection failure to explain, plus a second
client/process lifecycle.

Alternative 4 does not improve the Neovim operator's local diagnosis.

### Carol joins with a different client version or an unfamiliar presentation label

Alternative 1 can show Carol as a peer exactly as the existing awareness
surface reports her and mark the state as local awareness. It does not infer
trusted role continuity from a name or color.

Alternative 2 risks encouraging users to treat an available transport trace as
complete semantic interpretation across version boundaries.

Alternative 3 and Alternative 4 add no better answer to the distinction.

### The example evolves and activity grows

Alternative 1 can cap an in-memory local activity list and expose its source
and retention limit in the dashboard. No migration or persistent data format
is needed.

Alternative 2 needs filtering, redaction, storage/retention decisions, and
mixed-version compatibility for a new sidecar message shape.

Alternative 3 creates an independent presentation model that must evolve with
the plugin. Alternative 4 expands an explicitly excluded embodiment.

### Mallory attempts to exploit operator confusion at a trust boundary

Alternative 1 can retain explicit labels: relay state is local, peer labels
are presentation hints, and the dashboard does not create promises or grant
access. Its small scope minimizes misleading authority cues.

Alternative 2 increases the chance that raw bytes, a pCID, or a sender label
is mistaken for an authorization result. Alternative 3 and Alternative 4 do
not improve this boundary while widening the surface.

## Conclusions

- Reject Alternative 3: it duplicates the Neovim embodiment.
- Reject Alternative 4: it changes the browser, which is out of scope.
- Defer Alternative 2: a transport inspector is promising only after a
  separately scoped decision about sidecar schema, retention, filtering, and
  trust labelling.
- Alternative 1 survives as the recommended POC: a native Neovim dashboard
  plus a bounded, in-memory local activity view using existing state/events.

## Output to Decision Framing

The remaining user decisions are:

1. Lock Alternative 1 as the POC architecture, or defer the work.
2. Decide whether the dashboard is opened by new explicit commands only, or
   also automatically on `:GridEditorOpen`.
3. Approve the exact implementation and test paths after the command/behavior
   boundary is locked.

## Implications for TODOs and pending DIs

- TODO-numop tracks the proof of concept.
- A future DI must lock the approved architecture, behavior, implementation
  boundary, command names, and paths before code changes begin.
- The new work must preserve TE-zorud's sidecar boundary and the existing
  DI-samuv/DI-gafit presentation and awareness intent.

## Decision status

locked via DI-loril
