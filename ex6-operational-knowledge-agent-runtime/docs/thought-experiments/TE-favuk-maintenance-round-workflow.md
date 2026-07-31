# Maintenance round workflow example

TE ID: TE-favuk

## Status

decided

## Decision under test

How to add a third workflow example that demonstrates a maintenance concern
without adding new runtime behavior.

## Assumptions

- A workflow is a non-executing artifact that composes existing packages.
- Maintenance services and findings are already represented by the maintenance
  package, with run records for performed work.

## Alternatives

1. A maintenance-round workflow using `context`, `maintenance`, and `runs`.
2. A training-qualification workflow using `training` and `runs`.
3. A knowledge-review workflow using `knowledge` and `runs`.

## Scenario analysis

Option 1 starts with a known resource, records performed service, and retains a
finding when repair, escalation, or follow-up is required. This is visibly
different from inventory receipt and procedure execution while using existing
commands. Options 2 and 3 are valid future examples but do not show the
physical-resource and finding boundary provided by maintenance.

On a failed inspection, the workflow preserves the finding rather than
inventing automatic remediation. At scale, each round remains a small
content-addressed artifact while the maintenance and runs packages retain the
operational history.

## Conclusion

Choose option 1: add `maintenance-round`, composed from context, maintenance,
and runs. It adds no worker, route, top-level action, or runtime behavior.

## Decision status

Locked by DI-favuk.
