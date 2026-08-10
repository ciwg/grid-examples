# Ex2 testing guide

TE ID: TE-sapih

## Status

decided

## Decision status

Locked by DI-bubab: choose Alternative A. Create `docs/testing.md`, link it
from the README, document the full verification command set and distinct test
layers, and explain temporary test-root boundaries. This TE does not add or
change runtime behavior.

## Decision under test

How should Ex2 document its test layers so readers can run verification,
understand which claims each layer proves, and avoid mistaking relay-local
evidence or a one-host test run for centralized or globally authoritative
PromiseGrid proof?

## Assumptions and trust model

- Ex2 now has focused tests for pCID publication consistency, CAS/observation
  evidence, relay validation, and CRDT/awareness behavior.
- Existing service interoperability tests cover browser/Neovim behavior and
  relay-to-relay ingestion; they are the current decentralized cross-node
  coverage. Source: DI-guros.
- The canonical architecture is multi-relay peer collaboration. One host can
  run a genuine one-node session or simulate several relay nodes by using
  separate processes, data roots, and keys. Source: DI-nilas.
- The current profiles are local drafts, not frozen upstream specs or blanket
  peer-interoperability claims. Source: DI-ralit.
- Alice is a contributor running tests. Mallory may read a passing local test
  as a claim that a central relay, arbitrary peer, or physical machine is
  globally trusted.

## Alternatives

### A. One concise `docs/testing.md` guide linked from README

Document the commands, focused test layers, existing interoperability layer,
temporary data roots, and artifact/evidence boundaries in one exercise-local
guide. Keep architecture and protocol details in their existing documents.

### B. Expand the README with all testing detail

Place commands, test matrix, evidence boundaries, and topology explanation in
the main README.

### C. Leave test discovery to package names and CI commands

Keep tests runnable but do not create a dedicated guide.

## Scenario analysis

### Contributor changes a protocol document

Alice changes `live-document.md`. Under A, the guide tells her that the
profile-inventory test derives pCIDs from exact source bytes and checks the
published inventory/scope records. B can say the same but overloads the entry
README. C lets her discover the relationship only after a failure.

### Relay receives invalid or unsupported traffic

Mallory sends malformed bytes, an invalid proof, or a future pCID. A explains
that focused service tests prove local CAS/observation handling and exclusion
from accepted replay; the observations are not proof of Mallory's intent. B
risks burying the distinction in a long README. C leaves the evidence policy
hard to find.

### Browser, Neovim, and peer relay

Alice uses browser while Bob uses Neovim and their relays exchange traffic. A
points at existing interoperability tests as the cross-node layer and states
that they do not turn one relay key or display label into global authority. B
can describe it but makes first-entry prose long. C encourages readers to
assume package tests cover topology implicitly.

### One-host versus distributed simulation

Alice runs all test processes on one host. A states that Go tests use isolated
temporary data roots and that multiple roots/keys model independent relay
nodes; this is test topology, not a claim about physical-machine identity. B
adds more README complexity. C leaves the distinction undocumented.

### Maintenance and scale

A keeps a stable testing guide that every later alignment pass can update.
B makes the README grow with implementation detail. C has lowest immediate
cost but violates the repository alignment requirement and loses reader-facing
test provenance.

## Conclusions

C is rejected because the repository policy requires an exercise-local testing
guide during PromiseGrid alignment. B is rejected because Ex2's README already
serves as application orientation and should not duplicate detailed test
semantics.

A survives and is recommended: create `docs/testing.md`, link it from the
README, list `go vet ./...`, `go test ./...`, and `errcheck ./...`; explain
focused pCID/evidence tests, existing interoperability coverage, temporary
data roots, and the relay-local/decentralized evidence boundary.

## Decisions still requiring DF

1. **Guide layout:** create one concise `docs/testing.md` linked from README
   (recommended), expand README only, or no dedicated guide?
2. **Command set:** document `go vet ./...`, `go test ./...`, and `errcheck
   ./...` together (recommended), or only `go test ./...`?
3. **Layer explanation:** name focused pCID/evidence tests and existing
   interoperability tests with their different claims (recommended), or list
   packages without explaining their boundaries?
4. **Runtime note:** state that `t.TempDir()` roots isolate relay nodes and are
   automatically removed, not shared production evidence (recommended), or
   omit test-runtime storage guidance?

## Implications for open work

- `sojot.4` can proceed after these DFs are locked in `TODO-sojot`.
- `sojot.5` remains the final broader reader-guide pass and should link the
  testing guide rather than duplicate its detail.
