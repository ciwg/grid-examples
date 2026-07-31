# Inventory receipt workflow example

TE ID: TE-voruk

## Status

decided

## Decision under test

Which second example best demonstrates that the workflow loader is reusable
without adding new runtime behavior.

## Assumptions

- Workflows are immutable, non-executing artifacts.
- Existing packages provide the operational commands; a workflow only composes
  and documents them.
- The example must differ materially from procedure execution while remaining
  understandable in a receiving operation.

## Alternatives

1. A maintenance-round workflow using `maintenance` and `runs`.
2. An inventory-receipt workflow using `context`, `receiving`, `inventory`,
   and `runs`.
3. A second procedure variant using `procedures` and `runs`.

## Scenario analysis

In normal use, option 3 repeats the current package shape and does not prove
that the loader can carry a distinct operational concern. Option 1 is useful
but has less visible handoff between existing packages. Option 2 begins with a
physical receipt, records its disposition, and then records an inventory count
and reconciliation. It therefore demonstrates a distinct artifact that spans
four existing packages.

On incomplete receipt or a rejected count, the artifact remains inspectable and
the packages retain evidence and a disposition; no workflow engine invents an
automatic repair action. At larger scale, the workflow remains a small static
artifact and the existing record packages retain the operational history.

## Conclusion

Choose option 2: add an `inventory-receipt` workflow artifact. It documents
receipt, disposition, count, reconciliation, and evidence while relying only on
existing package commands. It adds no top-level action, route, worker, or
runtime behavior.

## Decision status

Locked by DI-voruk.
