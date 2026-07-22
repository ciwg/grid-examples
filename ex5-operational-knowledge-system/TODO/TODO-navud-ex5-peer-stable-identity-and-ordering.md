# TODO navud - introduce peer-stable identity and ordering for ex5 peer exchange

## Goal

Define and implement peer-stable durable identity and ordering semantics so
`ex5` can accept non-bootstrap peer exchange into already-populated runtimes
honestly.

## Why this exists

The current compatibility event model still uses runtime-local event sequences
and runtime-local entity IDs, so arbitrary non-bootstrap import across peers
cannot be implemented honestly yet.

## Tasks

- [ ] navud.1 Run the required TE for peer-stable identity and ordering across
  multi-origin runtimes.
- [ ] navud.2 Lock the first durable identity layer for imported artifacts and
  compatibility replay.
- [ ] navud.3 Lock how duplicate delivery, origin tracking, and ordering work
  once multiple peers contribute history.
- [ ] navud.4 Implement the first peer-stable identity/order slice.
- [ ] navud.5 Re-open TODO `103` implementation on top of that settled model.

## Status

- open
- created because TODO `103` is now locked to solve peer-stable identity and
  ordering before non-bootstrap import

