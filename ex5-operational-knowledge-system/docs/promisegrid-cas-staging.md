# ex5 PromiseGrid CAS staging

This note records the first shipped CAS-backed storage step in `ex5`. Source:
`DI-ribek`; `DI-lavuz`.

## First CAS step

The first CAS step is additive, not replacing:

- dual-write signed family envelopes into CAS
- dual-write copied evidence blobs into CAS
- keep `events.jsonl`, the current signed family logs, and the current
  `attachments/` tree as compatibility/manifests during migration

Source: `DI-ribek`.

This now ships in the runtime:

- signed family envelopes are dual-written into CAS by envelope CID
- copied evidence blobs are dual-written into CAS by blob CID
- the six frozen family envelopes now replay/export authoritatively from CAS
- compatibility event replay and attachment paths remain active for still-unfrozen
  runtime state
- imported evidence blobs now re-materialize into the local compatibility
  attachment tree from CAS when needed

Source: `DI-lavuz`; `DI-rovud`.

## Why the first step is additive

`ex5` already has stable local replay and verification over the current family
logs. The first CAS step should improve portability and content-addressability
without forcing an immediate storage cutover. Source: `DI-ribek`.

## What CAS unblocks

- later peer-visible evidence exchange over portable blob identities
- later movement of replay and read paths for still-unfrozen runtime state
  toward authoritative CAS-backed reads
- clearer long-horizon integrity and migration semantics for signed envelopes

Source: `DI-ribek`; `DI-tivor`; `DI-rovud`.

## What stays out of the first CAS step

- embodiment-contract tightening
- browser/CLI/Neovim migration away from the current local HTTP adapter
- immediate deletion of the current family logs or copied attachment tree

Source: `DI-ribek`.
