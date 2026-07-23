## Title

ex5 browser demo setup and verification boundary

## TE ID

TE-tevok

## Status

needs DF

## Decision under test

How `ex5-operational-knowledge-system` should make the browser demo runnable
from one sheet without hidden prerequisite guessing.

## Assumptions

- The current browser path is not demo-safe from the sheet alone.
- The user expects:
  - one doc
  - one runnable setup path
  - exact demo steps that work as written
- `ex5` already ships:
  - the checked-in sample corpus
  - the browser extension under `chrome-extension/`
  - the native host binary source under `cmd/operational-browser-host/`
  - the native-host manifest template under
    `chrome-extension/native-host/operational_browser_host.json`
- The browser embodiment fails closed when the direct Chrome/native-host path is
  unavailable.
- A browser demo that still depends on unstated manual setup is not honest
  enough for the promised one-sheet demo path.

## Alternatives

### Alternative A

Keep the current guide shape and just clarify the prerequisites more loudly.

- the browser sheet stays mostly the same
- the reader is still responsible for extension load, native-host
  registration, and readiness verification

### Alternative B

Add a dedicated browser-demo setup and verification path.

- one setup script prepares the sample runtime and native-host registration
- one Chrome launch path loads the shipped extension from a dedicated demo
  profile
- one verification step proves the browser embodiment is ready before the demo
  sheet is considered usable
- the user guide can then point to one exact preflight and one exact demo path

### Alternative C

Abandon browser-first demo prep and officially make the demo CLI-first.

- browser becomes optional bonus material
- the one-sheet path is only for CLI

## Scenario analysis

### Scenario 1: the user follows one sheet literally

Alice opens one doc and expects to do exactly what it says.

- Under A, Alice still has to infer what “browser embodiment set up” means in
  practice.
- Under B, Alice gets a discrete setup/preflight step and the demo sheet only
  becomes valid after that step passes.
- Under C, Alice gets the safest path, but it no longer satisfies the
  browser-first demo goal.

What gets easier:

- A: smallest immediate scope
- B: one honest browser-demo contract
- C: least setup risk

What gets harder:

- A: hidden failure points remain
- B: more setup tooling and verification work
- C: loses the stronger visual demo path

### Scenario 2: Chrome extension and native host drift

Bob has Chrome open and the extension appears loaded, but the native host is
missing or misregistered.

- Under A, Bob only discovers the problem by hitting the browser page and
  getting the unavailable message.
- Under B, the dedicated verification path can fail closed earlier and point at
  the real missing boundary.
- Under C, the issue disappears only because the browser is no longer the main
  demo path.

What gets easier:

- A: no new tooling
- B: setup truth is explicit and testable
- C: no browser setup burden

What gets harder:

- A: live-demo fragility remains
- B: the repo must own a real preflight path
- C: browser-first communication value is lost

### Scenario 3: repeatable short video capture

Carol wants to capture short browser clips later today.

- Under A, each recording session starts with ambient uncertainty about whether
  the browser embodiment is actually ready.
- Under B, the same setup and verification path can serve both the live demo
  and later recording sessions.
- Under C, CLI capture is simpler, but the strongest visual story is gone.

What gets easier:

- A: no additional buildout
- B: repeatability for both live and recorded demos
- C: simpler capture tooling

What gets harder:

- A: every session needs troubleshooting
- B: some scripting and browser-launch discipline is required
- C: weaker newcomer visual impact

### Scenario 4: PromiseGrid honesty

Dave asks whether the browser is really first-class or still a hand-assembled
demo trick.

- Under A, the docs are more honest than before, but the browser demo still
  depends on hidden operator skill.
- Under B, the repo owns the browser-demo boundary explicitly: setup,
  registration, launch, verification, then demo.
- Under C, the repo avoids overclaiming by moving to CLI-first, but also
  retreats from the shipped browser embodiment.

What gets easier:

- A: minimal documentation adjustment
- B: embodiment truth becomes operational, not just textual
- C: fewer moving parts

What gets harder:

- A: the one-sheet promise is still not met
- B: more repo-owned demo surface
- C: browser-first positioning weakens

### Scenario 5: maintenance over time

Ellen imagines future browser or host changes.

- Under A, documentation can drift away from real setup state again.
- Under B, one demo setup path and one verification path give the repo one
  place to keep current.
- Under C, browser maintenance is deferred, but the browser-demo goal remains
  unresolved.

What gets easier:

- A: smallest immediate maintenance
- B: one maintained browser-demo contract
- C: short-term simplification

What gets harder:

- A: repeated support cost
- B: setup tooling must be maintained
- C: the original user need stays unmet

## Conclusions

Rejected:

- Alternative A: too weak for the one-sheet browser-demo requirement
- Alternative C: does not satisfy the user’s explicit browser-demo goal

Surviving:

- Alternative B: dedicated browser-demo setup and verification path

Recommended:

- Alternative B is the most PromiseGrid-aligned surviving choice.

Why B is more aligned:

- it makes the browser embodiment operationally honest instead of merely
  documented
- it turns hidden prerequisites into an explicit preflight boundary
- it preserves the browser-first demo goal while failing closed if readiness is
  missing

## Implications for open TODOs and pending DIs

- A new ex5 TODO should track a browser-demo setup and verification wave.
- Final DF questions should lock:
  - whether setup includes a dedicated Chrome launch path
  - whether the extension is loaded via launch flags/profile instead of assumed
    manual state
  - what exact verification command or script must pass before the one-sheet
    demo is considered ready
