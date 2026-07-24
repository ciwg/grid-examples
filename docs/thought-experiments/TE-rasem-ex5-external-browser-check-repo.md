## Title

ex5 external browser-check repo boundary

## TE ID

TE-rasem

## Status

needs DF

## Decision under test

How `ex5` should validate real browser-demo behavior outside the main
`grid-examples` repo, including the external repo shape and whether the first
version should use Playwright or a lighter browser-record/replay approach.

## Assumptions

- The user does not want browser-demo validation artifacts in the main repo.
- The current in-repo tests are strong on correctness, but they are still not
  enough to guarantee an obvious live-demo path from a presenter's point of
  view.
- The user wants something broader than `ex5` alone, so the external repo
  should be able to grow to `ex1`, `ex2`, `ex3`, `ex4`, and `ex5`.
- The preferred external repo name is `~/lab/cswg/grid-examples-browser-checks/`.
- `ex5` already has:
  - one checked-in sample corpus
  - one scripted browser setup / launch / verify path
  - one canonical browser demo sheet in `docs/user-guide.md`
- The goal is not to replace in-repo tests. The goal is to add presenter-grade
  browser validation for real visible interactions.

## Alternatives

### Alternative A

Manual rehearsal and checklist only.

- no external code repo
- no repeatable browser automation
- all demo validation stays human-driven

### Alternative B

External browser-check repo with Playwright as the primary harness.

- one separate repo: `~/lab/cswg/grid-examples-browser-checks/`
- one subdirectory per example: `ex1/`, `ex2/`, `ex3/`, `ex4/`, `ex5/`
- Playwright drives a real browser, clicks exact UI controls, checks exact
  visible text, and can capture screenshots or traces when something fails

### Alternative C

External browser-check repo with lighter record/replay tooling first.

- still separate from the main repo
- rely on browser recorder exports, shell wrappers, or ad hoc scripts first
- possibly smaller first step than Playwright

## Scenario analysis

### Scenario 1: today's ex5 live demo prep

Alice needs to know whether the exact browser path will work before she shows
it to anyone.

- Under A, Alice still rehearses by hand and can still miss ambiguous visual
  handoffs until the moment she presents.
- Under B, Alice can run one repeatable browser flow that checks visible
  browser facts like `Current Record`, `Problem hotspots`, and the actual
  result text after a click.
- Under C, Alice gets some automation, but it is more likely to be brittle,
  opaque, or tied to one local browser capture.

What gets easier:

- A: no new tooling
- B: exact click path and exact visible text assertions
- C: smaller first automation step

What gets harder:

- A: no repeatable proof beyond rehearsal
- B: introduces a real test harness
- C: recorder artifacts tend to age poorly and explain failures badly

### Scenario 2: browser-demo regressions that are visible but not semantic

Bob changes the UI and the data is still correct, but the presenter can no
longer tell what changed on screen.

- Under A, the issue is only found during manual rehearsal.
- Under B, the external checks can assert visible anchors such as `Current
  Record`, the literal lane text, and the exact record that appears after
  `Inspect`.
- Under C, some regressions can be caught, but recorder-style tools tend to be
  weaker at maintaining readable assertions over time.

What gets easier:

- A: nothing new to maintain
- B: checks the specific class of failures the user hit today
- C: lower barrier than a full harness

What gets harder:

- A: demo trust stays fragile
- B: external assertions must be curated
- C: lower long-term clarity and debuggability

### Scenario 3: future growth across multiple examples

Carol wants browser checks for `ex3`, `ex4`, and `ex5`, and wants future-her to
remember exactly which product repo they align to.

- Under A, there is no reusable structure.
- Under B, `grid-examples-browser-checks/` with `ex1/` through `ex5/` keeps the
  alignment obvious and gives one place for cross-example browser validation
  habits.
- Under C, the same repo shape can exist, but the per-example browser harness
  is more likely to diverge into one-off scripts.

What gets easier:

- A: nothing to organize
- B: one stable multi-example browser-check home
- C: still keeps the checks outside the product repo

What gets harder:

- A: no durable testing discipline
- B: needs one consistent harness choice
- C: likely fragmentation between examples

### Scenario 4: failure investigation

Dave hits a failing browser-demo check and wants to know what happened.

- Under A, there is no machine-readable failure artifact.
- Under B, Playwright can give a clear failing step, visible assertion, and
  optional screenshot or trace.
- Under C, recorder-style output usually tells less coherent stories when a
  UI flow drifts.

What gets easier:

- A: no harness debugging
- B: clearer failure evidence
- C: lightweight first capture

What gets harder:

- A: human-only diagnosis
- B: setup cost for the harness
- C: harder to keep failures readable

### Scenario 5: PromiseGrid honesty

Ellen asks whether these browser checks belong in the product repo or are a
separate operational confidence layer.

- Under A, the product repo stays clean, but there is no durable confidence
  layer.
- Under B, the boundary is explicit: correctness and contract tests remain in
  the product repo, while presenter-grade browser-path checks live outside it.
- Under C, the boundary is still external, but the tooling choice is weaker as
  a long-term confidence layer.

What gets easier:

- A: strongest simplicity
- B: clean separation of concerns
- C: external boundary still holds

What gets harder:

- A: the user still lacks adequate browser confidence
- B: one more repo to maintain
- C: the external repo may become a bag of ad hoc scripts

## Conclusions

Rejected:

- Alternative A: too weak for the user's stated need because it leaves browser
  demo validation mostly manual.

Surviving:

- Alternative B: external browser-check repo with Playwright first
- Alternative C: external browser-check repo with lighter record/replay tooling

Recommended:

- Alternative B is the most PromiseGrid-aligned surviving choice.

Why B is more aligned:

- it preserves the cleanliness of the main repo
- it gives one durable, multi-example browser-check home aligned explicitly to
  `grid-examples`
- it tests the exact visible presenter path the current in-repo tests do not
  fully cover
- it produces clearer failure evidence than recorder-style or ad hoc scripts

## Implications for open TODOs and pending DIs

- A new ex5 TODO should track the external browser-check repo decision and the
  first ex5 slice.
- Final DF questions should lock:
  - whether `grid-examples-browser-checks/` is the external repo root
  - whether Playwright is the first harness
  - what the first `ex5` checks are
  - whether the first version should include screenshots/traces or just visible
    assertion checks
