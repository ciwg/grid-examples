# ex5 PromiseGrid embodiment staging

This note records the first shipped embodiment-contract tightening step in
`ex5`. Source: `DI-vabek`; `DI-rovuz`.

## Current state

Browser, CLI, and Neovim still share one runtime through the current local HTTP
adapter. That remains the delivery surface today, but the adapter now exposes
runtime capability metadata for peer exchange, CAS, and shared-draft storage
through `Meta`. Source: `DI-sobek`; `DI-vabek`; `DI-rovuz`; `DI-bavuk`.

## Tightening trigger

The first tightening trigger was both runtime layers landing in implemented
form:

- the first relay-visible exchange layer
- the additive CAS-backed storage layer

That first tightening step now ships as:

- adapter-visible runtime capability metadata
- doc wording that describes HTTP as an adapter over a peer/CAS-capable
  runtime, not only as a local app surface
- explicit confirmation that HTTP remains the sole embodiment adapter for now

Source: `DI-vabek`; `DI-rovuz`; `DI-bavuk`.

## What not to do yet

- do not rewrite current browser, CLI, or Neovim as if they bypass the adapter
- do not confuse adapter-visible capability metadata with a finished transport
  migration
- do not invent a second embodiment contract unless a narrower future TE shows
  the local HTTP adapter is actually blocking needed runtime behavior

Source: `DI-vabek`; `DI-rovuz`; `DI-bavuk`.
