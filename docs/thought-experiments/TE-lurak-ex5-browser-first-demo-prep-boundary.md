## Title

ex5 browser-first demo prep boundary

## TE ID

TE-lurak

## Status

needs DF

## Decision under test

How `ex5-operational-knowledge-system` should prepare a browser-first demo,
CLI proof slice, and optional short-video recording helpers without creating a
second product surface that drifts away from the shipped newcomer path.

## Assumptions

- The date is July 23, 2026.
- `ex5` already ships one checked-in newcomer corpus under
  `ex5-operational-knowledge-system/sample-data/newcomer-runtime/`.
- `docs/user-guide.md` is already the canonical newcomer-facing operator path.
- The user wants three artifacts:
  - a demo TODO
  - a demo script
  - optional recording helpers for short videos
- The browser should be the main demo embodiment.
- The CLI should be shown as proof that the same operational world is usable
  from a shell-first embodiment too.
- The recording-helper boundary should stay small unless a broader helper layer
  is clearly justified by the shipped surfaces.

## Alternatives

### Alternative A

Write only lightweight presentation notes and leave the rest informal.

- no canonical demo path file
- no checked-in demo sequence
- no recording helpers

### Alternative B

Ship one browser-first demo pack anchored to the checked-in newcomer corpus.

- one canonical demo TODO that tracks the wave
- one canonical demo script that follows the same sample world as the user
  guide
- one small CLI proof segment embedded in that same story
- optional thin recording helpers only where they remove repeatability risk

### Alternative C

Ship a broad multi-embodiment demo toolkit.

- separate primary demo flows for browser, CLI, and Neovim
- a larger recording-helper surface
- more automation around capture, resets, and possibly clip generation

## Scenario analysis

### Scenario 1: a first live demo for newcomers

Alice is showing `ex5` to a newcomer audience. She needs one clear narrative
 that starts from the checked-in sample data and does not require improvising
 record IDs or workflow steps in real time.

- Under A, Alice still has to improvise the sequence and may drift from the
  newcomer guide.
- Under B, Alice gets one stable browser-first path that can reuse the same
  storyline, IDs, and review flow already taught in the user guide.
- Under C, Alice gets more coverage, but also has to choose between several
  competing top-level demo flows.

What gets easier:

- A: smallest immediate documentation change
- B: one canonical newcomer-aligned demo path
- C: widest embodiment coverage

What gets harder:

- A: live-demo drift and higher presenter error risk
- B: requires curation of one strong path
- C: scope, maintenance, and rehearsal burden rise quickly

### Scenario 2: proving the browser is primary while the CLI is real

Bob wants the main demo to live in the browser, but he also wants to prove
 this is not a browser-only toy.

- Under A, Bob may mention the CLI, but he has no canonical proof slice.
- Under B, the browser remains the main story while one short CLI segment
  proves the same runtime and sample world are usable directly from the shell.
- Under C, CLI and maybe Neovim can get full equal-weight demo tracks, but the
  browser-first promise becomes less clear.

What gets easier:

- A: no extra demo material to maintain
- B: browser-first story with one honest terminal proof step
- C: broad embodiment parity storytelling

What gets harder:

- A: embodiment breadth is under-proven
- B: CLI slice must stay deliberately short and evidence-based
- C: the main story can fragment across embodiments

### Scenario 3: preparing short videos from the live demo flow

Carol wants short 30-to-90-second clips after the live demo is stable.

- Under A, there is no repeatable capture scaffolding, so clip capture becomes
  ad hoc.
- Under B, small helpers can focus on the repeatability bottlenecks: loading a
  clean sample runtime, a canonical command list, and maybe one thin capture
  wrapper or checklist.
- Under C, a more ambitious toolset can automate more, but it risks becoming a
  demo-production subsystem rather than a thin support layer.

What gets easier:

- A: no helper maintenance
- B: reliable repetition without overbuilding
- C: richer production support

What gets harder:

- A: each recording session becomes more fragile
- B: helper scope must stay disciplined
- C: tool sprawl and non-product-facing complexity

### Scenario 4: alignment with the shipped newcomer path

Dave compares the proposed demo flow with `docs/user-guide.md`,
`sample-data/newcomer-runtime/`, and the browser’s review-first UI.

- Under A, the live demo can drift into a custom narrative that newcomers
  cannot replay.
- Under B, the demo can start from the same sample load, the same storylines,
  and the same browser review queue order already recommended to newcomers.
- Under C, several demo tracks may tempt the repo into maintaining separate
  “showcase” paths that no longer match the newcomer guide.

What gets easier:

- A: minimal repo churn
- B: one shared truth between newcomer onboarding and demo prep
- C: wider event/demo options

What gets harder:

- A: weaker replay value after the demo
- B: demo editing must stay disciplined and product-honest
- C: higher drift risk between onboarding and demo collateral

### Scenario 5: maintenance over time

Ellen imagines future `ex5` UI or workflow changes.

- Under A, there is little direct maintenance, but every future demo depends on
  human memory.
- Under B, one canonical demo script and a small helper layer give the repo one
  place to update when the browser-first flow changes.
- Under C, multiple demo tracks and larger helpers create several places for
  drift and breakage.

What gets easier:

- A: smallest checked-in surface
- B: one maintainable demo path
- C: broader future event support

What gets harder:

- A: repeatability remains weak
- B: demo assets still need to be kept current
- C: maintenance surface grows much faster than user value

## Conclusions

Rejected:

- Alternative A: too weak for the user’s stated need
- Alternative C: too broad for the first honest demo-prep wave

Surviving:

- Alternative B: one browser-first demo pack anchored to the checked-in sample
  corpus, with a short CLI proof slice and only thin recording helpers

Recommended:

- Alternative B is the most PromiseGrid-aligned surviving choice.

Why B is more aligned:

- it keeps the demo path anchored to the same durable sample world and operator
  guide the repo already ships
- it proves embodiment breadth without pretending all embodiments need equal
  demo weight
- it treats recording helpers as support for repeatability, not as a second
  product/tooling project

## Implications for open TODOs and pending DIs

- A new ex5 TODO should track the browser-first demo-prep wave.
- The implementation should likely stay centered on:
  - one canonical browser-first demo script
  - one short CLI proof segment
  - one small set of recording helpers only where repeatability would otherwise
    be fragile
- Final DF questions should lock:
  - the exact browser story arc
  - the exact CLI proof slice
  - the exact helper boundary and file paths
