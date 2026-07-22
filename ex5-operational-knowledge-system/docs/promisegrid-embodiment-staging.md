# ex5 PromiseGrid embodiment staging

This note records when `ex5` should tighten embodiment language beyond the
current local HTTP adapter. Source: `DI-vabek`.

## Current state

Browser, CLI, and Neovim still share one runtime through the current local HTTP
adapter. That remains the honest shipped embodiment contract today. Source:
`DI-sobek`; `DI-vabek`.

## Tightening trigger

Embodiment-contract tightening should wait until both of these runtime layers
exist in implemented form:

- the first relay-visible exchange layer
- the additive CAS-backed storage layer

Only after those exist should `ex5` restate embodiments as binding more
directly to the shipped PromiseGrid runtime contract. Source: `DI-vabek`.

## What not to do yet

- do not rewrite current browser, CLI, or Neovim docs as if they already speak
  a direct peer/runtime contract
- do not couple embodiment wording changes to planning-only peer/storage docs

Source: `DI-vabek`.
