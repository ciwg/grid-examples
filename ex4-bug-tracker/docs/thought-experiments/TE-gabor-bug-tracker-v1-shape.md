# Ex4 Bug Tracker V1 Shape

TE ID: TE-gabor
## Status
decided

## Decision under test

What is the smallest useful first implementation shape for a non-editor example
application that still exercises durable workflow, multiple roles, attachments,
and more than one embodiment?

## Assumptions

- The app should be usable on its own, not only as a throwaway demo.
- V1 should stay small enough to implement and verify in one pass.
- Browser usability matters more than live collaboration.
- The repo already has stronger examples for distributed routing (`ex1`) and
  richer browser embodiment (`ex2`), so `ex4` does not need to repeat either
  one exactly.

## Alternatives

1. Realtime-first support desk with websocket collaboration.
2. Generic CRUD ticket app with mutable rows.
3. Durable-first bug tracker with append-only issue history and a simple CLI.

## Scenario analysis

### Normal operation

Alternative 1 adds transport and presence complexity before the issue workflow
is proven useful. Alternative 2 is easy to ship but loses the durable history
shape that makes the example more instructive. Alternative 3 keeps the
workflow legible: create, triage, assign, work, resolve, and reopen all show up
as explicit timeline events.

### Failure and incomplete writes

Alternative 1 creates more failure surfaces because live sessions and durable
storage interact. Alternative 2 can overwrite issue state without preserving a
reasonably inspectable trail. Alternative 3 can rebuild current state from an
append-only log and keeps attachments as copied artifacts under app-managed
storage.

### Multiple embodiments

Alternative 1 strongly favors the browser. Alternative 2 can add a CLI but
usually treats it as a second-class management tool. Alternative 3 naturally
supports a browser queue/detail UI and a small engineer CLI against the same
issue history.

### Future expansion

Alternative 1 risks over-committing to live transport. Alternative 2 risks
needing a storage rewrite when later features demand better history. Alternative
3 leaves seams for team scoping, richer permissions, and more automation
without changing the core storage shape.

## Conclusions

The surviving alternative is the durable-first bug tracker:

- browser-first queue + detail UI
- simple built-in identities
- append-only issue events
- copied file attachments
- small engineer CLI
- hidden built-in team field for future multi-team work

## Implications for TODOs and DIs

- The foundational TODO is `TODO-valop-bug-tracker-foundation.md`.
- The locked implementation decisions are recorded as `DI-dajak`,
  `DI-nunit`, `DI-ninuf`, and `DI-gofub`.

## Decision status

locked
