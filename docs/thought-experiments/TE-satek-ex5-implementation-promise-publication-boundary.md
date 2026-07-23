## Title

ex5 implementation promise publication boundary

## TE ID

TE-satek

## Status

needs DF

## Decision under test

How `ex5-operational-knowledge-system` should publish its implementation
promise claims so they stay aligned with the PromiseGrid dev guide: whether the
repo should keep one family-only `CHANGELOG.md`, expand that `CHANGELOG.md`
into a component-aware claim surface, or split claims across multiple files.

## Assumptions

- The PromiseGrid dev guide says App Devs should use B-side implementation
  promise claims and that `CHANGELOG.md` should make compatibility scope
  auditable by exact spec doc-CID and explicit scope.
- The same guide also says a multi-component app should not hide behind one
  vague app-wide label if different pieces speak different parts of the
  contract.
- `ex5` now has materially different components:
  - local runtime under `service/` and `cmd/operational-knowledge`
  - browser embodiment through `web/`, `chrome-extension/`, and
    `cmd/operational-browser-host`
  - CLI under `cmd/oks-cli`
  - Neovim embodiment under `nvim/`
  - remote relay under `cmd/operational-relay`
- The current `CHANGELOG.md` still describes browser, CLI, and Neovim as if
  they were using the old local HTTP adapter directly, so the current
  publication layer is stale even though the higher-level boundary docs are up
  to date.

## Alternatives

### Alternative A

Keep one family-only `CHANGELOG.md`, but update its stale embodiment wording.

- Refresh the current family claim entries so they describe the shipped direct
  embodiment contracts accurately
- Do not add component-aware sections
- Keep `CHANGELOG.md` as a short family list only

### Alternative B

Keep `CHANGELOG.md` as the single publication surface, but make it
component-aware.

- Preserve the family claim entries for the frozen durable families
- Add explicit component-level implementation promise sections for the local
  runtime, remote relay, browser embodiment helper path, and terminal
  embodiment surfaces
- Let one file answer both questions:
  - which frozen families are claimed
  - which shipped components implement or delegate which contract surfaces

### Alternative C

Split implementation claims across multiple files now.

- Keep family claims in `CHANGELOG.md`
- Move component claims to one or more new docs/files
- Use the higher-level docs to cross-link the pieces

## Scenario analysis

### Scenario 1: an app developer checks ex5 compatibility quickly

Alice wants to know whether `ex5` speaks the frozen specs she intends to use
and whether the current browser/CLI/Neovim surfaces are direct or compatibility
paths.

- Under A, she gets current family claims again, but still lacks one place that
  states which component actually implements which surface.
- Under B, she can read one file and see both the family claims and the
  per-component implementation/delegation line.
- Under C, she may get the same truth eventually, but only by reading and
  reconciling multiple files.

What gets easier:

- A: smallest repair to stale wording
- B: one auditable publication surface
- C: cleaner separation between family and component claims

What gets harder:

- A: component-level honesty remains weaker than the dev guide expects
- B: one file becomes broader
- C: readers have to compose the answer from several files

### Scenario 2: the browser embodiment changes again

Bob changes the browser extension/native-host path or adds one more direct
contract slice.

- Under A, the family claims may still stay current, but the component promise
  boundary remains under-described.
- Under B, the changed component can be updated directly in the same
  implementation promise surface App Devs are already told to consult.
- Under C, the repo risks drift between family claims, component claim files,
  and summary docs.

What gets easier:

- A: narrow edits
- B: component changes stay attached to the implementation-promise source
- C: family and component stories can evolve independently

What gets harder:

- A: the browser helper path is still not claimed explicitly enough
- B: the publication file needs careful organization
- C: cross-file drift becomes more likely

### Scenario 3: relay and local runtime evolve at different speeds

Carol adds a relay-only feature while the local runtime stays the same.

- Under A, `CHANGELOG.md` still mainly describes families, so the relay
  component boundary remains implicit.
- Under B, the relay can carry its own scoped claim section without needing a
  separate file layout.
- Under C, separate files can express that cleanly too, but at the cost of a
  more complex publication surface before there is evidence that several files
  are actually needed.

What gets easier:

- A: no file-structure change
- B: per-component scope stays explicit in one place
- C: stronger formal separation if the repo grows much larger

What gets harder:

- A: local runtime and relay still blur together for readers
- B: the single file needs disciplined sectioning
- C: packaging/publication complexity rises early

### Scenario 4: PromiseGrid dev-guide compliance as the primary yardstick

Dave compares `ex5` directly to the dev guide.

- Under A, the stale wording problem is fixed, but the multi-component claim
  problem is only partially addressed.
- Under B, the guide's two strongest expectations are both met:
  exact frozen family claims remain auditable and different shipped components
  no longer hide behind one vague label.
- Under C, the repo might also satisfy the guide, but only if the reader can
  easily discover and join the claim files.

What gets easier:

- A: quick partial alignment
- B: strongest direct guide alignment
- C: possible long-term publication structure

What gets harder:

- A: still weaker than the guide's component-honesty standard
- B: more editorial work now
- C: more navigational burden now

### Scenario 5: long-horizon maintenance and repo honesty

Ellen wants the claim publication layer to stay honest over time without
turning into a mini-framework of its own.

- Under A, the repo keeps the smallest surface but leaves important component
  truth outside the implementation promise claims.
- Under B, the repo keeps one canonical publication file while making its
  scope honest enough for the current system shape.
- Under C, the repo may be future-proofing too early by multiplying files
  before the publication surface is proven too large for one document.

What gets easier:

- A: minimal editing discipline
- B: one canonical, honest claim surface
- C: future expansion room

What gets harder:

- A: under-claims component truth
- B: requires a stronger file structure inside `CHANGELOG.md`
- C: risks complexity before necessity

## Conclusions

Rejected:

- Alternative A: it fixes the stale transport wording, but it does not meet the
  dev guide's stronger multi-component implementation-promise discipline.

Surviving:

- Alternative B: one component-aware `CHANGELOG.md`
- Alternative C: split publication across multiple files

Recommended:

- Alternative B is the strongest PromiseGrid-aligned surviving choice.

Why B is more aligned than C:

- The dev guide already points App Devs to `CHANGELOG.md` as the B-side
  implementation promise surface.
- `ex5` needs stronger honesty, not yet a more fragmented publication system.
- One file can still hold the current family claims plus explicit component
  claim sections without forcing readers to stitch together the answer from
  multiple locations.

Alternative C remains viable later if the current file truly becomes too large
or if the component claim surface becomes independently versioned enough to
justify separate publication artifacts.

## Implications for open TODOs and pending DIs

- A new ex5 TODO should track the dev-guide publication-alignment fix wave.
- That TODO should likely lock Alternative B unless the user explicitly prefers
  a more fragmented claim-publication structure.
- The next implementation pass should update `CHANGELOG.md` first, then align
  any summary docs that currently point to it as the implementation-promise
  source of truth.
