# Browser And Neovim Transport Doc Drift

TE ID: TE-lurem
## Status
decided

## Decision under test

How broad the `136` doc-alignment pass should be now that the shipped browser
contract is Chrome/Chromium native messaging, Neovim is socket-first with
explicit compatibility mode, and some long-form docs still describe older
websocket-first or HTTP-primary transport language.

This TE corresponds to TODO `zunek.1` / `zunek.2` / `zunek.3`.

## Assumptions

- Alice reads the README and high-level guides first when deciding what `ex5`
  actually ships.
- Bob uses the CLI and Carol uses Neovim; both depend on the docs being honest
  about direct-contract versus compatibility transport.
- The code and tests are already the source of truth for the current shipped
  embodiment behavior; this wave is about documentation honesty, not runtime
  repair.
- The remaining drift is concentrated in high-visibility long-form docs rather
  than the already-updated HTTP API guide and recent TODO/TE records.

## Alternatives

### A. Minimal stale-paragraph patch

Only rewrite the exact paragraphs already known to be stale in the README,
features guide, architecture, and practical implementation notes.

### B. High-visibility transport sweep

Rewrite the known stale paragraphs and re-sweep the main high-visibility doc
surfaces for nearby transport framing so each major surface tells the same
browser and Neovim story in one pass.

### C. Full doc corpus sweep

Re-audit every ex5 doc for transport wording, including lower-priority notes
and adjacent protocol/backlog commentary, before closing `136`.

## Scenario analysis

### Scenario 1: a new reader decides what ex5 ships today

Alice reads the README, then the features guide, then one implementation guide.

#### A. Minimal stale-paragraph patch

What it makes easier:
- smallest change set
- low risk of unnecessary doc churn

What it makes harder:
- surrounding nearby wording may still imply the older transport story
- readers can still encounter mixed framing across the major docs

New obligations:
- accept that another doc sweep may be needed soon

#### B. High-visibility transport sweep

What it makes easier:
- the main reader journey becomes consistent
- browser and Neovim embodiment claims line up across the top-level docs
- PromiseGrid alignment is clearer without touching every niche note

What it makes harder:
- broader doc editing than the minimum patch

New obligations:
- define which surfaces count as the high-visibility set

#### C. Full doc corpus sweep

What it makes easier:
- strongest consistency if done perfectly

What it makes harder:
- much larger review/edit surface
- more chances to churn low-signal text without changing meaning

### Scenario 2: operator trust in embodiment claims

Bob wants to know whether browser still depends on websocket-over-HTTP and
whether Neovim still silently falls back across adapters.

#### A. Minimal stale-paragraph patch

What it makes easier:
- fixes the most obvious contradictions

What it makes harder:
- neighboring text can still soften or blur the actual boundary

#### B. High-visibility transport sweep

What it makes easier:
- the primary operator-facing docs all tell the same contract story
- easier to trust that browser directness and Neovim explicit compatibility are
  not accidental wording artifacts

What it makes harder:
- requires slightly more editorial discipline

#### C. Full doc corpus sweep

What it makes easier:
- can remove almost all remaining wording drift

What it makes harder:
- scope expands well beyond the concrete operator-trust problem

### Scenario 3: maintenance cost

Dave wants a doc pass that is large enough to be worth doing, but not so large
that it becomes a sprawling prose rewrite.

#### A. Minimal stale-paragraph patch

What it makes easier:
- quick closure

What it makes harder:
- raises the chance of another near-term follow-on for missed neighboring text

#### B. High-visibility transport sweep

What it makes easier:
- closes the meaningful reader-facing drift in one pass
- stays bounded to the docs people actually read first

What it makes harder:
- requires a slightly wider manual consistency check

#### C. Full doc corpus sweep

What it makes easier:
- likely pushes residual drift furthest down

What it makes harder:
- high cost for diminishing returns

### Scenario 4: PromiseGrid alignment

Steve wants the docs to reflect the shipped embodiment truth as directly as the
code does.

#### A. Minimal stale-paragraph patch

What it makes easier:
- removes obvious inaccuracies

What it makes harder:
- still treats the doc issue as isolated lines instead of one cross-surface
  transport story

#### B. High-visibility transport sweep

What it makes easier:
- makes the core surfaces consistently say:
  - browser: Chrome/Chromium native messaging
  - CLI: local socket with explicit HTTP opt-in
  - Neovim: local socket first, explicit compatibility mode
- best balance of honesty and bounded scope

What it makes harder:
- modestly larger patch

#### C. Full doc corpus sweep

What it makes easier:
- maximum consistency ambition

What it makes harder:
- broader than the actual PromiseGrid problem requires right now

## Conclusions

Rejected:

- Alternative A: too narrow; it risks leaving the main reader journey partly
  inconsistent.
- Alternative C: too broad for the concrete drift that remains.

Surviving:

- Alternative B: high-visibility transport sweep

Alternative B is the most PromiseGrid-aligned remaining path because it makes
the main doc surfaces speak with one honest embodiment/transport story without
turning `136` into a full-corpus editorial rewrite.

## Implications for open TODOs and pending DIs

- TODO `136` should lock a bounded but cross-surface doc sweep, not just a
  one-line patch and not a whole-corpus rewrite.
- Locked result: `136B.1`, meaning the high-visibility transport sweep is
  limited to the already-approved six files.
