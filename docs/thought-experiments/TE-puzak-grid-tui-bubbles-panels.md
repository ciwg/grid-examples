# Focused Bubbles panels for Grid TUI

TE ID: TE-puzak

## Status

needs DF

## Decision under test

How Grid TUI exposes new Bubbles components from its existing top-level menus
without changing the Ex3 document, awareness, relay, browser, or Neovim
protocol behavior. This analyzes `maduh.1` through `maduh.7`.

## Assumptions and boundary

Grid TUI's document and awareness transport remain implementation-specific
local UI machinery. The bundled PromiseGrid Development Guide 0.1.0, App Devs
guidance, treats UI structure as local implementation detail rather than a
protocol contract. Alice and Bob may use different UI clients against the same
document; neither choice changes promises, pCIDs, persisted text, or presence
payloads.

## Alternatives

1. Focused panels and dialogs opened by the existing custom menu bar.
2. Replace the menu bar with one command-palette overlay.
3. Add components inline beside the editor continuously.

## Scenario analysis

### Normal collaboration

Alice selects `Help` or `Collaborate` while editing. Focused panels preserve
the editor as the dominant surface and return focus predictably when closed.
A command palette is fast for experts but makes existing labeled menu
navigation disappear. Inline components permanently reduce space for shared
text and peer context.

### Relay failure or incomplete data

Bob opens catalog or published-version metadata while the relay is slow or
returns no records. A focused panel can show a spinner, empty state, and retry
without hiding the current document. A command palette conflates navigation
with asynchronous results. Inline output risks stale data becoming visually
ambiguous beside active editing.

### Concurrent or mixed-version clients

Alice's newer TUI may have Help, List, File picker, Table, or Progress while
Bob's browser and Neovim do not. Focused local panels add no message fields and
therefore preserve interoperability. The same is true of the other alternatives,
but replacing the menu changes local muscle memory more broadly.

### Long-horizon evolution

Panels give each component an explicit job: key help, document/result choice,
file selection, metadata, or operation state. That makes later replacement or
removal narrow. A palette becomes a second navigation architecture; permanent
inline surfaces accumulate layout pressure as features grow.

### Trust and scale

Panels do not add durable state or broaden trust. Lists and tables can bound
large results with filtering or scroll regions. A future progress indicator is
only honest when the operation exposes real incremental work; otherwise a
spinner is the correct transient signal.

## Conclusions

Alternative 1 survives: it preserves the established menu vocabulary and
gives each Bubbles component a focused, inspectable purpose. Alternative 2 is
rejected for this task because it replaces rather than improves an existing
interaction model. Alternative 3 is rejected because collaboration needs text
space more than permanent controls.

The remaining implementation choice is whether this conclusion locks the
named Bubbles components as focused menu-opened panels/dialogs. If locked,
their local state stays transient and protocol-neutral.

## Implications for open work

`maduh.1` through `maduh.7` may be implemented as focused panels/dialogs,
with deterministic view and update tests. `maduh.8` and `maduh.9` remain
deferred.
