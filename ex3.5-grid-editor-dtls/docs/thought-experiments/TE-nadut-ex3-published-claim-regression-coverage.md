# Ex3 published-claim regression coverage

TE ID: TE-nadut

## Status

decided

## Decision status

Locked by DI-dilav: add focused source-derived publication/boundary coverage
and retain the current decentralized interoperability suites. This TE does not
change remote-admission, evidence, WebSocket, or browser behavior.

## Decision under test

How should Ex3 prove that its published local-draft pCIDs, provisional
remote-admission boundary, relay-local evidence policy, and decentralized
browser/Neovim/relay paths remain true without duplicating slow integration
tests for every low-level claim?

## Assumptions and trust model

- `CHANGELOG.md` and `docs/architecture.md` publish four local-draft pCIDs
  derived from exact `protocols/*.md` source bytes. Source: DI-hadil.
- Accepted replay contains only fully supported envelopes; rejected peer
  envelopes create local observations, and capability denials create separate
  non-secret diagnostics. Source: DI-pazis; DI-darif; DI-lozut.
- Existing service interoperability tests cover browser/Neovim behavior,
  peer-relay exchange, late join, and remote capability use; browser JavaScript
  and sidecar tests cover their embodiment-local behavior.
- Alice changes a profile document, Bob joins from a second relay or embodiment,
  and Mallory sends malformed ingress or a denied remote request.

## Alternatives

### A. Focused publication and boundary tests plus existing interoperability

Add a source-derived pCID inventory test; strengthen focused ingress/admission
tests to prove durable local evidence, secret exclusion, and accepted replay
exclusion; retain the existing browser, Neovim, peer-relay, WebSocket, and
late-join suites as decentralized cross-node coverage.

### B. One comprehensive end-to-end topology for every claim

Run remote browser, Neovim, two relays, WebSocket, every rejected envelope
class, admission denial, and published pCID comparison in one scenario.

### C. Documentation-only publication review

Rely on manual inspection of pCIDs and retain only current behavior tests.

## Scenario analysis

### Profile-source revision

Under A, a direct test derives the pCID from each exact profile source and
compares both published records, producing a small actionable failure. B adds
topology setup without improving a documentation-drift diagnosis. C permits
reader-facing values to become stale.

### Rejected peer envelope

A asserts raw CAS retention, one observation per receipt, and absence from
accepted feed/restart input. B adds WebSocket and embodiment machinery to a
relay-local fact. C leaves the DI-pazis contract vulnerable to regression.

### Denied remote admission

A verifies that a denial adds a non-secret diagnostic but no accepted message
or bearer/bootstrap material. B makes that assertion timing- and transport-
heavy. C relies on code review to protect the sensitive boundary.

### Browser, Neovim, and multi-relay collaboration

A retains existing interoperability, headless late-join, browser-JS, and
sidecar suites as the appropriate cross-node and embodiment proof. B repeats
already-tested paths for every scope claim. C does not clarify which current
tests carry the decentralized behavior claim.

### Long-horizon maintenance

A keeps source publication, evidence policy, admission diagnostics, and
interoperability independently testable as profile drafts and mechanisms
evolve. B grows a slow flaky test that obscures causes. C concentrates trust
in undocumented manual review.

## Conclusions

C is rejected because published pCIDs and evidence boundaries are reader-facing
claims that can drift. B is rejected as the sole strategy because local
publication and evidence assertions do not need a full distributed topology to
be meaningful.

A survives and is recommended: add focused source-derived inventory coverage;
prove local evidence/admission boundaries directly; retain—not duplicate—the
existing browser/Neovim, peer-relay, WebSocket, and late-join coverage.

## Decisions still requiring DF

1. **Coverage strategy:** use focused publication/boundary tests plus existing
   interoperability (recommended), one full topology for every claim, or
   documentation-only publication review?
2. **Inventory authority:** derive every pCID from `protocols/*.md` and compare
   both `docs/architecture.md` and `CHANGELOG.md` (recommended), compare only
   runtime constants, or compare just one document?
3. **Boundary assertion:** test raw CAS retention, one observation per receipt,
   accepted replay/feed exclusion, and non-secret admission diagnostics
   (recommended), test only log-file presence, or test accepted errors only?
4. **Cross-node scope:** retain the existing interoperability, headless,
   browser-JS, sidecar, and WebSocket suites as decentralized coverage
   (recommended), add duplicate topology coverage, or remove those suites from
   the alignment claim?

## Implications for open work

- `fozoz.3` can add focused profile inventory and boundary tests after DF.
- `fozoz.4` must explain the layers and commands in `docs/testing.md`.
- TODO tamuk remains open; its manual private-browser verification is not
  replaced by these deterministic suites.

## Refinements

### 2026-08-10 — Private-browser verification resolved

TODO tamuk is now complete. Isolated normal and incognito Chrome sessions
converged document text in both directions through one isolated relay. The
browser-level check used local DevTools and native browser input, not a
human-driven usability review; see `DI-sodoj` in
`TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md` for the exact
evidence and constraint.
