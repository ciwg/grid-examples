# TODO nobek - decide and stage CAS-backed envelope storage for ex5 PromiseGrid families

## Goal

Define when and how ex5 should move from local append-only family logs toward
CAS-backed envelope storage as part of the shipped PromiseGrid runtime.

## Why this exists

The current implementation persists signed family logs locally, but CAS-backed
envelope storage is still outside the shipped ex5 operational workflow.

## Tasks

- [ ] nobek.1 Run the required TE for CAS-backed storage scope and migration
  order.
- [ ] nobek.2 Lock what stays in compatibility logs and what moves to CAS.
- [ ] nobek.3 Define the first staged implementation slice.
- [ ] nobek.4 Add storage-boundary docs and migration notes.

## Status

- open

