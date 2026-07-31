# Receiving-check workflow example

TE ID: TE-zovuk

## Status

decided

## Decision under test

How to represent ex5's receiving-check operational flow as an ex6 workflow
artifact without adding runtime behavior.

## Assumptions

- Workflows compose existing packages and do not execute themselves.
- Receiving records receipts and dispositions; runs retain performed-work
evidence and approval.

## Alternatives

1. A receiving-check artifact using `context`, `receiving`, and `runs`.
2. Treat inventory receipt as the receiving-check equivalent.

## Scenario analysis

Option 2 conflates two different operator questions: whether inbound goods pass
inspection, and how accepted goods are counted and reconciled in inventory.
Option 1 keeps the receiving decision explicit: inspect an incoming item,
record evidence, then accept, quarantine, or reject it. Failed inspection is
retained as a disposition instead of being silently folded into inventory.

## Conclusion

Choose option 1: add `receiving-check` as a distinct workflow artifact using
the existing context, receiving, and runs packages.

## Decision status

Locked by DI-zovuk.
