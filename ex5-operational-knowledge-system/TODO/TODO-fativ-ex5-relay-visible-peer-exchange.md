# TODO fativ - decide and stage relay-visible ex5 PromiseGrid peer exchange

## Goal

Define the first relay-visible peer-exchange slice for ex5 after the remaining
core durable families are frozen.

## Why this exists

The current shipped ex5 runtime is PromiseGrid-native only at the local durable
family layer. Relay-visible peer exchange is still explicitly not implemented.

## Tasks

- [ ] fativ.1 Run the required TE for the first relay-visible exchange scope.
- [ ] fativ.2 Lock what is exchanged first, by whom, and under what trust
  assumptions.
- [ ] fativ.3 Define the first staged runtime and storage changes.
- [ ] fativ.4 Add tracking docs for what becomes peer-visible.

## Status

- open

