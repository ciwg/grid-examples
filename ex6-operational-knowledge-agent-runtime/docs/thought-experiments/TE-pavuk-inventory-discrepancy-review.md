# Inventory discrepancy review workflow example

TE ID: TE-pavuk

## Status

decided

## Decision under test

How to represent ex5's inventory discrepancy review as an ex6 workflow artifact
without conflating it with receiving or routine inventory counting.

## Assumptions

- Inventory records counts and reconciliations; runs retain evidence and
approval.
- A workflow documents an operator sequence and does not execute it.

## Alternatives

1. An inventory-discrepancy-review artifact using `context`, `inventory`, and
   `runs`.
2. Treat inventory receipt or routine count as the discrepancy-review
   equivalent.

## Scenario analysis

Option 2 loses the question that matters after an unexpected count: why the
physical quantity and the expected record differ, and what local decision was
made. Option 1 records the count, retains supporting evidence, and records a
reconciliation decision such as adjust, investigate, or reject. It keeps a
receiving event separate from a later audit decision.

## Conclusion

Choose option 1: add `inventory-discrepancy-review`, composed from existing
context, inventory, and runs packages. It adds no worker, route, top-level
action, or runtime behavior.

## Decision status

Locked by DI-pavuk.
