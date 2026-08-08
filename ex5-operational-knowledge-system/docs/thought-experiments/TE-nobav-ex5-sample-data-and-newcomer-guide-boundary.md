## Title

ex5 sample data and newcomer guide boundary

## TE ID

TE-nobav

## Status

needs DF

## Decision under test

How `ex5-operational-knowledge-system` should add checked-in sample data and a
single newcomer-ready user guide without blurring operator-facing product truth
with implementation/reference material.

## Assumptions

- The user has already locked two boundary choices:
  - checked-in deterministic sample data under `ex5-operational-knowledge-system/sample-data/`
  - one primary comprehensive operator guide at
    `ex5-operational-knowledge-system/docs/user-guide.md`
- The current doc set is broad but split:
  - `user-guide.md` is workflow-oriented but short
  - `product-overview.md`, `browser-ui-guide.md`, and
    `terminal-capability-matrix.md` cover narrower slices
  - reference/technical docs stay separate
- The newcomer guide should be enough for a new operator to understand the
  system, but it should not turn into a duplicate of implementation claims,
  API reference, or architecture docs.
- Sample data should help a newcomer see the full operational-memory model
  end-to-end rather than only toy single-record examples.

## Alternatives

### Alternative A

Ship a minimal sample corpus and only lightly expand the current user guide.

- sample data covers a very small happy path
- guide remains mostly command examples plus cross-links
- newcomer still needs to hop quickly into the other docs

### Alternative B

Ship one rich checked-in sample corpus and one canonical newcomer guide.

- sample data shows a full context-to-review lifecycle:
  responsibilities, places, resources, items, runs, evidence, approvals, links
- the user guide becomes the single start-here operator document:
  concepts, setup, how to load the sample, how to inspect it, and how to use
  browser/CLI/Neovim against it
- deeper UI and reference docs remain linked, not duplicated

### Alternative C

Ship several sample corpora and a giant all-in-one guide.

- multiple scenarios such as receiving, inventory, maintenance, and training
- the primary guide absorbs most browser, terminal, and reference material
- newcomer gets one very large document and several data sets immediately

## Scenario analysis

### Scenario 1: a newcomer opens the repo cold

Alice wants to understand what ex5 is, what records matter, and how a real
operational thread moves through the system.

- Under A, she gets a shorter path, but still has to assemble the system model
  from several guides quickly.
- Under B, she gets one clear start-here guide plus one real sample corpus she
  can inspect across all embodiments.
- Under C, she gets maximum material, but the first-read burden becomes much
  heavier and the main guide risks becoming hard to navigate.

What gets easier:

- A: smallest doc/data change
- B: one canonical newcomer path
- C: broadest immediate coverage

What gets harder:

- A: newcomer still has to compose the mental model
- B: requires more deliberate editorial structure
- C: information overload and higher maintenance

### Scenario 2: an operator wants realistic practice data

Bob wants to click through browser review, run CLI triage, and inspect Neovim
surfaces without inventing his own records first.

- Under A, the corpus may be too thin to show the real lifecycle.
- Under B, one richer corpus can include draft items, approved items,
  superseded items, runs with evidence, grouped problems, and typed links.
- Under C, several corpora may show more breadth, but they also multiply the
  newcomer decision burden.

What gets easier:

- A: small fixture maintenance
- B: realistic practice on one known data set
- C: broader scenario coverage

What gets harder:

- A: sample data may not prove the product's actual value
- B: the corpus needs careful curation
- C: corpus sprawl and higher upkeep

### Scenario 3: PromiseGrid honesty versus demo polish

Carol asks whether the sample data should be generator-shaped, toy-shaped, or
durable-record-shaped.

- Under A, the sample can drift toward toy examples.
- Under B, the checked-in corpus can be treated like a stable demonstrator of
  the current product/runtime slice, with explicit IDs and inspectable
  relationships.
- Under C, the guide risks drifting toward a polished brochure rather than one
  durable, inspectable operational-memory corpus.

What gets easier:

- A: easy initial setup
- B: stable, inspectable product truth
- C: broad marketing-style examples

What gets harder:

- A: lower evidence value
- B: more editorial discipline needed
- C: more chance of duplicating or overstating behavior

### Scenario 4: keeping operator docs separate from technical/reference docs

Dave wants the user guide to be enough for use, but not a second API manual.

- Under A, the current split stays clearer, but the newcomer guide may remain
  too thin.
- Under B, the user guide can cover concepts, common tasks, embodiment choice,
  and sample-corpus walkthroughs while still linking out for browser-detail,
  terminal matrix, and HTTP/reference material.
- Under C, the main guide absorbs too much and starts duplicating the browser
  guide, terminal matrix, and technical claims docs.

What gets easier:

- A: less doc overlap
- B: strongest operator-facing boundary
- C: one-file convenience

What gets harder:

- A: newcomer still needs too many jumps
- B: requires strict section boundaries
- C: drift and contradiction risk rise

### Scenario 5: maintenance over time

Ellen imagines the guide and sample data after several more ex5 changes.

- Under A, maintenance is lighter but the newcomer value stays capped.
- Under B, one canonical guide plus one canonical sample corpus gives the repo
  a clear place to keep current.
- Under C, several corpora and one giant guide create many places for drift.

What gets easier:

- A: smallest maintenance burden
- B: one canonical operator narrative and one canonical sample corpus
- C: future expansions can slot into existing bulk

What gets harder:

- A: continued newcomer friction
- B: each major product change should be reflected in one important doc and one
  corpus
- C: maintenance cost and drift surface multiply

## Conclusions

Rejected:

- Alternative A: too weak for the stated newcomer goal
- Alternative C: too broad and too maintenance-heavy for the current need

Surviving:

- Alternative B: one rich checked-in sample corpus plus one canonical
  newcomer-ready user guide

Recommended:

- Alternative B is the strongest PromiseGrid-aligned surviving choice.

Why B is more aligned:

- it keeps the sample corpus durable, inspectable, and checked in
- it gives newcomers one canonical operator path
- it preserves a clean boundary between product-facing guidance and deeper
  technical/reference docs

## Implications for open TODOs and pending DIs

- A new ex5 TODO should track the sample-data and newcomer-guide wave.
- That TODO should likely lock a single rich sample corpus, not several corpora
  at once.
- The user-guide work should remain product-facing and cross-link to the
  browser guide, terminal matrix, and reference docs rather than absorbing
  them.
