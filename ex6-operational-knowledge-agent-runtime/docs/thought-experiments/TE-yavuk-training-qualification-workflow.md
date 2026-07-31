# Training-qualification workflow example

TE ID: TE-yavuk

## Status

decided

## Decision under test

How to represent ex5's training qualification flow as an ex6 workflow artifact
without adding runtime behavior.

## Assumptions

- Training sessions and certifications are already supported by the training
  package, with runs retaining performed-session history.
- A workflow documents the operator sequence but does not execute it.

## Alternatives

1. A training-qualification artifact using `training` and `runs`.
2. Treat procedure execution as sufficient for training.

## Scenario analysis

Option 2 omits the trainee, instructor, completion, and certification boundary
that distinguishes qualification from merely performing a procedure. Option 1
records a session, retains evidence, and makes a certification decision
explicit. A failed or incomplete session remains a durable training outcome,
not an automatic qualification.

## Conclusion

Choose option 1: add `training-qualification`, composed from the existing
training and runs packages. It adds no worker, route, top-level action, or
runtime behavior.

## Decision status

Locked by DI-yavuk.
