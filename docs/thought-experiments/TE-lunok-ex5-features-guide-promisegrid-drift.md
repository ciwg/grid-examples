# ex5 features guide PromiseGrid drift

TE ID: TE-lunok
## Status
decided

## Decision under test

How `docs/features-guide.md` should be rewritten so it no longer understates
the shipped signed-envelope, relay, and current Neovim embodiment scope.

This is not a new product-scope question. The runtime and the other summary
docs already say what ships. The problem is that the feature guide still
describes an older “browser and CLI over the local runtime” shape and still
uses a pre-relay PromiseGrid note.

## Assumptions

- README, product overview, relay guide, and PromiseGrid claims doc are now
  the more accurate sources for current shipped scope.
- The feature guide should match those documents, not define a different
  boundary.
- No runtime behavior changes are needed here.

## Alternatives

### Alternative A: fully align the feature guide with the shipped scope now

Rewrite the opening scope framing and any stale embodiment wording so the
feature guide describes:

- browser, CLI, and Neovim as current embodiments
- signed-envelope and relay behavior as shipped within the current scope
- the remaining future-scope items only as explicit boundaries

### Alternative B: keep the older opening framing and append a correction note

Preserve the old wording near the top, then add a later correction note that
the runtime now also ships signed-envelope and relay layers.

## Scenario analysis

### Scenario 1: reader starts at the feature guide

Alice reads only `docs/features-guide.md`.

Alternative A gives her the same current-scope story the rest of the repo now
gives.

Alternative B keeps the first impression wrong and depends on later careful
reading to undo it.

### Scenario 2: PromiseGrid review

Bob compares the README, claims doc, and feature guide.

Alternative A makes them agree.

Alternative B preserves an avoidable contradiction that future reviewers will
flag again.

### Scenario 3: future maintenance

Carol later updates one embodiment or relay behavior.

Alternative A leaves one consistent baseline for further edits.

Alternative B keeps stale framing around as legacy prose debt.

## Conclusions

Rejected:

- Alternative B. It keeps the known contradiction in place.

Surviving:

- Alternative A: fully align the feature guide with the shipped scope now

Recommendation:

- Alternative A

Why:

- It is the cleanest PromiseGrid doc-alignment path.
- It removes the last remaining summary doc contradiction found in the review.

## Implications for TODOs and pending DIs

- TODO `126` is locked to Alternative `A`.
- The implementation should rewrite the feature-guide opening scope statements
  and confirm no summary contradiction remains afterward.
