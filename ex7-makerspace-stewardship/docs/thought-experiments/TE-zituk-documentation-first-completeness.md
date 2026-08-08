# Documentation-first Ex7 completeness

TE ID: TE-zituk

## Status

needs DF

## Decision under test

What evidence must exist before Ex7 can be called broadly complete, with
documentation treated as a product requirement rather than a by-product of the
Go service implementation.

## Assumptions

- Ex7 is a single-process local makerspace demonstration with append-only local
  evidence, no authentication, no replication, and no signatures.
- The repaired persistence scope already has passing Go service and HTTP tests.
- A new reader should be able to run, understand, verify, and safely limit the
  example without reading its implementation.
- Alice records an observation and loan, Carol clears a safety hold, and Mallory
  represents the local filesystem and untrusted-input boundary described by the
  existing evidence TEs.

## Alternatives

1. **Documentation-first completion standard.** Require an operator guide,
   workflow/evidence guide, explicit limits, reproducible demo, and targeted
   browser/API verification before broad completion.
2. **Code-and-tests-only standard.** Treat passing Go tests and a short README
   as sufficient.
3. **UI-first standard.** Build browser automation and polish before defining
   the reader-facing guide and limits.

## Scenario analysis

### New volunteer onboarding

Alice starts the demo for the first time. She needs to know that in-space use
is not checkout, what a safety hold does, who may clear it, how a loan records
accepted terms, and where evidence persists.

Alternative 1 gives Alice one guided path and makes the local-only trust limit
visible. Alternative 2 leaves those facts scattered across source, UI labels,
and TEs. Alternative 3 can make controls look polished while their operational
meaning remains undocumented.

### Evidence failure and recovery

Mallory corrupts the local evidence log. Carol restarts the service and needs
to understand why startup fails, that the log is preserved, and how a demo can
be reset without claiming repair.

Alternative 1 requires this failure behavior and reset boundary in the guide
and demo verification. Alternative 2 proves it only through tests. Alternative
3 may automate a happy-path browser flow while omitting the safety boundary.

### Long-horizon review

Months later, a maintainer revisits a legacy loan with incomplete accepted
terms. They need to distinguish known facts from absent evidence and verify the
browser presents that distinction.

Alternative 1 requires a workflow/evidence reference plus targeted UI/API
evidence. Alternative 2 relies on implementation reading. Alternative 3 risks
testing only appearance rather than provenance.

### Scope and deployment trust

A group considers using the demo beyond one trusted volunteer machine.

Alternative 1 puts the no-authentication, no-replication, and local-filesystem
limits next to run instructions, reducing accidental overclaim. Alternatives 2
and 3 make the limitations easier to miss.

## Conclusions

Alternative 2 is rejected: tests alone do not make an example understandable
or safe to evaluate. Alternative 3 is rejected as an ordering: UI automation
must validate a documented operational claim, not define it afterward.

Alternative 1 survives and is recommended.

## Output to decision framing

The proposed broad-completion gate contains:

1. An operator guide with setup, run, reset, and verification steps.
2. A workflow/evidence guide covering observation/hold, hold clearance, loan,
   return, accepted policy terms, and incomplete legacy terms.
3. An explicit scope and trust-limits section.
4. A reproducible browser/API demo path.
5. Targeted tests that prove the documented happy path and evidence-failure
   boundary.

Remaining DF choices are the guide structure, demo medium, browser-test depth,
exact file paths, and whether the standard is Ex7-only or reusable across
examples.

## Decision status

needs DF
