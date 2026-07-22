# TODO lurog - tighten the ex5 embodiment contract beyond the current HTTP adapter

## Goal

Decide whether later `ex5` embodiments should keep routing strictly through
the local HTTP adapter or expose a more direct grid-native runtime contract.

## Why this exists

The current embodiment story is now honest and stable, but it is still
explicitly HTTP-adapter-first. If `ex5` is to become more fully on-grid, that
boundary may need a later tightening pass.

## Tasks

- [ ] lurog.1 Run the required TE for later embodiment-contract tightening.
- [ ] lurog.2 Lock whether browser, CLI, and Neovim stay adapter-first or gain
  a more direct runtime-facing contract.
- [ ] lurog.3 Implement the next concrete tightening step, if any.
- [ ] lurog.4 Update docs, claims, and operator guidance to match the chosen
  embodiment contract.

## Status

- open
- created from the post-109 PromiseGrid review

