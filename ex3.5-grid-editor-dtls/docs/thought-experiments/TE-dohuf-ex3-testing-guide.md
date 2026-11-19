# Ex3 testing guide

TE ID: TE-dohuf

## Status

decided

## Decision status

Locked by DI-hosit: create the concise guide, document grouped command and
test layers, and retain the temporary-root and manual private-browser
boundaries. This TE does not add or change relay, capability, WebSocket,
browser, or sidecar behavior.

## Decision under test

How should Ex3 document its verification commands and test layers so readers
can distinguish source-derived publication, relay-local evidence, remote
admission, browser/Neovim/relay interoperability, and the still-open manual
private-browser condition?

## Assumptions and trust model

- Go coverage verifies pCID inventory, relay ingress/evidence, remote
  admission, WebSocket paths, browser/Neovim interoperability, and headless
  late join. Source: DI-dilav.
- Browser JavaScript has `npm test` and `npm run build`; the Neovim sidecar has
  `npm run build`. Those are embodiment-local checks, not protocol authority.
- `t.TempDir()` gives Go tests isolated relay identities and data roots that
  are automatically cleaned up; they are not shared production evidence.
- Automated hardening covers likely private/incognito startup failures, but
  TODO tamuk still requires a real manual private-window verification. Source:
  DI-figak; DI-hadil.
- Alice runs the commands. Bob joins through another relay or embodiment.
  Mallory may treat one passing local suite as global trust proof.

## Alternatives

### A. One concise `docs/testing.md` linked from README

Document the complete Go, browser, and sidecar command groups; explain the
focused publication/evidence/admission layer separately from existing
decentralized interoperability suites; explain temporary data roots; and state
the manual private-browser boundary.

### B. Expand README with all testing detail

Put commands, layer matrix, temporary-root semantics, and private-browser
caveat directly in the application README.

### C. Leave tests discoverable only through package scripts

Keep commands runnable but provide no dedicated reader-facing guide.

## Scenario analysis

### Alice changes a profile document

Under A, the guide identifies the source-derived inventory test and its two
published documents. B can say the same but overloads the entry README. C
forces Alice to infer the relationship from failures or package names.

### Mallory sends malformed ingress or lacks a capability

A explains that focused service/store tests prove relay-local observations,
accepted replay exclusion, and non-secret admission diagnostics—not global
sender validity. B buries that distinction in application orientation. C makes
the evidence boundary hard to discover.

### Bob joins through browser, Neovim, or a second relay

A identifies the existing interoperability, headless late-join, browser-JS,
sidecar, and WebSocket layers as decentralized behavior evidence. B makes the
README longer. C invites a reader to mistake any single unit suite for the
whole collaboration claim.

### Private/incognito browser session

A clearly separates automated hardening from the open manual verification in
TODO tamuk. B can do so but duplicates feature-guide material. C permits a
passing test run to be read as completion.

### Maintenance and scale

A gives each document one job: README orientation, architecture semantics,
testing verification. B creates repeated detail. C violates the alignment rule
requiring an exercise-local testing guide.

## Conclusions

C is rejected because the repository alignment protocol requires a testing
guide. B is rejected because Ex3's README is already extensive and should
remain an application entry point.

A survives and is recommended: create `docs/testing.md`, link it from README,
document all command groups and layer boundaries, explain temporary roots, and
explicitly retain TODO tamuk's manual private-browser status.

## Decisions still requiring DF

1. **Guide layout:** create one concise `docs/testing.md` linked from README
   (recommended), expand README only, or no dedicated guide?
2. **Command set:** include Go `go vet ./...`, `go test ./...`, `errcheck
   ./...`; browser `npm test` and `npm run build`; and sidecar `npm run build`
   (recommended), list Go only, or list all commands without grouping?
3. **Layer explanation:** name focused pCID/evidence/admission tests separately
   from existing interoperability, headless, browser-JS, sidecar, and WebSocket
   suites (recommended), list packages only, or describe all tests as one
   generic suite?
4. **Runtime/private-browser note:** explain `t.TempDir()` isolation and state
   TODO tamuk's still-open manual private-window verification (recommended),
   explain temporary roots only, or omit both boundaries?

## Implications for open work

- `fozoz.4` can create the guide and README link after DF is locked.
- `fozoz.5` should link the guide rather than duplicate it.
- TODO tamuk remains open and cannot be completed by this documentation pass.

## Refinements

### 2026-08-10 — Private-browser verification resolved

TODO tamuk is now complete. Isolated normal and incognito Chrome sessions
converged document text in both directions through one isolated relay. The
browser-level check used local DevTools and native browser input, not a
human-driven usability review; see `DI-sodoj` in
`TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md` for the exact
evidence and constraint.
