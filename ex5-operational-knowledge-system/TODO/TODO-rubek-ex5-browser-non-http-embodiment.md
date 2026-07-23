# TODO rubek - ex5 browser non-http embodiment

## Decision Intent Log

ID: DI-rubek
Date: 2026-07-22 18:12:55 -0700
Author: jj@thesalleys.com (JJ)
Status: active
Decision: Track a future-scope pass to move the browser onto a direct non-HTTP embodiment contract.
Intent: Improve the biggest remaining embodiment-layer PromiseGrid impurity after the terminal surfaces have already moved onto the local socket contract.
Constraints: Treat this as a broader scope wave; it will likely require new transport and embodiment-boundary decisions rather than a small cleanup.
Affects: `ex5-operational-knowledge-system/web/*`, `ex5-operational-knowledge-system/service/*`, `ex5-operational-knowledge-system/docs/*`, `ex5-operational-knowledge-system/TODO/TODO.md`

## Goal

Define and stage a browser embodiment contract that no longer relies on the
current local HTTP adapter as its primary runtime surface.

## Tasks

- [ ] rubek.1 Define the first direct browser embodiment contract candidate.
- [ ] rubek.2 Compare it against the current local HTTP adapter behavior and migration cost.
- [ ] rubek.3 Lock the first browser non-HTTP slice and stage implementation.

## Status

- open
- future-scope PromiseGrid refinement
