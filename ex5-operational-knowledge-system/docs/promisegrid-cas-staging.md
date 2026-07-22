# ex5 PromiseGrid CAS staging

This note records how `ex5` should introduce CAS-backed storage without
breaking the current local durable runtime. Source: `DI-ribek`.

## First CAS step

The first CAS step is additive, not replacing:

- dual-write signed family envelopes into CAS
- dual-write copied evidence blobs into CAS
- keep `events.jsonl`, the current signed family logs, and the current
  `attachments/` tree as compatibility/manifests during migration

Source: `DI-ribek`.

## Why the first step is additive

`ex5` already has stable local replay and verification over the current family
logs. The first CAS step should improve portability and content-addressability
without forcing an immediate storage cutover. Source: `DI-ribek`.

## What CAS unblocks

- later peer-visible evidence exchange over portable blob identities
- later movement of read paths from log-only replay toward CAS-backed reads
- clearer long-horizon integrity and migration semantics for signed envelopes

Source: `DI-ribek`.

## What stays out of the first CAS step

- embodiment-contract tightening
- browser/CLI/Neovim migration away from the current local HTTP adapter
- immediate deletion of the current family logs or copied attachment tree

Source: `DI-ribek`.
