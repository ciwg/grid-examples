# Ex3 final README alignment

TE ID: TE-gotan

## Status

decided

## Decision status

Locked by DI-tajin: add concise orientation, explicit isolated relay commands,
demo-token guidance, local-evidence navigation, and testing-guide navigation.
This TE does not alter profiles, relay behavior, capability issuance, WebSocket
carriage, or browser startup.

## Decision under test

How should Ex3's README make its local-draft scope, multi-relay topology,
provisional remote-admission demo, relay-local evidence, verification guide,
and outstanding private-browser manual condition clear without duplicating the
canonical scope, architecture, and testing documents?

## Assumptions and trust model

- The four pCIDs are repo-local draft profiles. WebSocket is carriage, not a
  fifth profile. Source: DI-hadil.
- A bootstrap secret is relay-local operator configuration used to mint
  short-lived scoped capabilities for remote mutation; it is not a universal
  auth API or person identity system. Source: DI-hadil; DI-povip.
- Accepted replay differs from relay-local observations and admission
  diagnostics; neither local evidence record is shared proof. Source:
  DI-pazis; DI-darif; DI-lozut.
- The testing guide is canonical for commands and test layers. TODO tamuk's
  real normal/private browser verification remains open. Source: DI-hosit.
- Alice runs one relay, Bob runs another, and Mallory may read a shared host,
  a bootstrap token, or a local diagnostic as a global trust claim.

## Alternatives

### A. Add a concise orientation section and reproducible commands

Use explicit distinct `--data-root` values for local and two-relay examples;
add a concise Scope, topology, admission, and evidence section linked to the
canonical documents; label `ex3-demo-access` as a checked-in demo-only token
and tell operators not to reuse it; replace the stale final test snippet with
the testing-guide link.

### B. Keep README as a feature/demo walkthrough

Leave scope, evidence, roots, and complete verification in linked documents;
correct only the final test snippet if needed.

### C. Copy detailed scope, architecture, and testing content into README

Place pCID tables, capability claim fields, evidence schemas, test matrix, and
private-browser procedure directly in the README.

## Scenario analysis

### Alice starts an isolated local relay

Under A, Alice uses an explicit root and knows where that relay's artifacts
live. B leaves the implicit default-root behavior as the primary command. C
can explain it but overloads the application entry point.

### Bob joins through another relay

A shows separate relay roots and the `--peer` relationship, making a
same-host two-process run a decentralized-node simulation rather than a shared
authoritative server. B relies on readers finding topology detail elsewhere.
C duplicates the architecture and risks drift.

### Remote demo token

A makes the checked-in `ex3-demo-access` token visibly demo-only and explains
that real operators choose their own local bootstrap secret. It retains the
current example flow while avoiding a reusable-secret implication. B leaves
the sample value easily mistaken for normal deployment guidance. C creates a
large authorization tutorial around an intentionally provisional mechanism.

### Mallory sends rejected input or is denied remote admission

A directs Alice to the selected relay root and the architecture/testing guides
for `observations.jsonl` and `admission-diagnostics.jsonl`, while saying they
are local facts rather than shared proof. B makes this boundary less
discoverable. C duplicates sensitive evidence details and creates a second
source of truth.

### Private/incognito browser session

A links the testing guide and TODO tamuk, clearly separating automated
hardening from still-pending manual verification. B can leave the existing
caveat buried later in README. C repeats an evolving troubleshooting record.

### Long-horizon maintenance

A keeps README as concise orientation and leaves exact pCIDs, storage, test
layers, and non-claims canonical elsewhere. B under-documents the primary
entry point. C creates duplicated contracts that can diverge.

## Conclusions

C is rejected because it duplicates canonical scope, architecture, and testing
material. B is rejected because Ex3's primary entry point currently leaves
root isolation, remote-demo safety, evidence locality, and complete
verification too implicit.

A survives and is recommended: add concise orientation and reproducible
commands, label the demo token correctly, link canonical detail, and replace
the stale one-command Tests section.

## Decisions still requiring DF

1. **README structure:** add one concise Scope, topology, admission, and
   evidence section linked to canonical guides (recommended), correct only the
   test text, or duplicate detailed guides?
2. **Operation examples:** show one explicit local `--data-root` and a
   separate-root two-relay `--peer` example (recommended), keep implicit
   default-root commands, or point only to Docker?
3. **Remote demo guidance:** retain `ex3-demo-access` only as clearly labeled
   checked-in demo value and tell operators to choose their own local bootstrap
   secret (recommended), present it as normal deployment configuration, or
   remove remote guidance from README?
4. **Evidence/verification navigation:** point readers to their selected relay
   root and the architecture/testing guides, link TODO tamuk, and replace the
   stale `go test ./...` snippet (recommended), list internal details in README,
   or omit the boundary?

## Implications for open work

- `fozoz.5` can perform the README pass after DF is locked.
- TODO tamuk remains open and cannot be completed by this documentation change.
- Any generalized identity, delegation, or remote authorization redesign needs
  a new TE and DI rather than a README edit.

## Refinements

### 2026-08-10 — Private-browser verification resolved

TODO tamuk is now complete. Isolated normal and incognito Chrome sessions
converged document text in both directions through one isolated relay. The
browser-level check used local DevTools and native browser input, not a
human-driven usability review; see `DI-sodoj` in
`TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md` for the exact
evidence and constraint.
