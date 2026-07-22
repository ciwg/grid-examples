# TODO zidor - decide when ex5 embodiments should move beyond the local HTTP adapter contract

## Goal

Define the point at which browser, CLI, and Neovim should stop being described
primarily through the local HTTP adapter and instead bind more directly to the
shipped PromiseGrid runtime contract.

## Why this exists

The current docs correctly describe the local HTTP API as the embodiment
adapter. Getting ex5 fully on-grid will eventually require a cleaner statement
about how embodiments relate to the direct runtime contract.

## Tasks

- [ ] zidor.1 Run the required TE for embodiment-contract tightening timing and
  scope.
- [ ] zidor.2 Lock the staged boundary between local HTTP adapter behavior and
  direct runtime contract behavior.
- [ ] zidor.3 Define the first embodiment-facing migration slice, if any.
- [ ] zidor.4 Update the external and repo docs once that boundary is settled.

## Status

- open

