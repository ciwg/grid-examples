# ex5 PromiseGrid embodiment staging

This note records the first shipped embodiment-contract tightening step in
`ex5`. Source: `DI-vabek`; `DI-rovuz`.

## Current state

Browser, CLI, and Neovim still share one runtime through the current local HTTP
adapter. That remains the delivery surface today, but the adapter now exposes
runtime capability metadata for peer exchange and CAS through `Meta`. Source:
`DI-sobek`; `DI-vabek`; `DI-rovuz`.

## Tightening trigger

The first tightening trigger was both runtime layers landing in implemented
form:

- the first relay-visible exchange layer
- the additive CAS-backed storage layer

That first tightening step now ships as:

- adapter-visible runtime capability metadata
- doc wording that describes HTTP as an adapter over a peer/CAS-capable
  runtime, not only as a local app surface

Source: `DI-vabek`; `DI-rovuz`.

## What not to do yet

- do not rewrite current browser, CLI, or Neovim as if they bypass the adapter
- do not confuse adapter-visible capability metadata with a finished transport
  migration
- do not describe ex5 as fully on-grid while peer exchange is still
  bootstrap-only and CAS remains a sidecar read model

Source: `DI-vabek`; `DI-rovuz`; `DI-tivor`.
