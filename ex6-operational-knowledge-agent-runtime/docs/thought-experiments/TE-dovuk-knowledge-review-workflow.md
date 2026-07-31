# Knowledge review workflow example

TE ID: TE-dovuk

## Status

decided

## Decision under test

How to represent ex5's knowledge-item review lifecycle as an ex6 workflow
artifact without adding a separate review runtime.

## Assumptions

- Knowledge already supports item creation, revision snapshots, approval, and
  supersedence.
- A workflow artifact documents review behavior and does not execute it.

## Alternatives

1. A knowledge-review artifact using the existing `knowledge` package.
2. Treat procedure execution as the review equivalent.

## Scenario analysis

Option 2 omits the review lifecycle around reusable operational knowledge.
Option 1 retains the distinction between drafting a revision, approving the
current item, and superseding an older item. A rejected or incomplete review
does not automatically change the item lifecycle.

## Conclusion

Choose option 1: add `knowledge-review`, composed from the existing knowledge
package. It adds no worker, route, top-level action, or runtime behavior.

## Decision status

Locked by DI-dovuk.
