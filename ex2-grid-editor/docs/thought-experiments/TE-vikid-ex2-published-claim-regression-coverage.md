# Ex2 published-claim regression coverage

TE ID: TE-vikid

## Status

decided

## Decision status

Locked by DI-guros: use focused source-derived inventory and evidence tests,
and retain the existing interoperability suite as decentralized topology
coverage.

## Decision under test

How should Ex2 regression coverage prove that its four published local-draft
pCIDs remain derived from their exact source specs, that relay exception
evidence stays separate from accepted replay, and that decentralized relay and
cross-embodiment behavior remain covered without duplicating tests needlessly?

## Assumptions and trust model

- The architecture inventory and scope declaration publish four pCIDs derived
  from `protocols/*.md`. Source: DI-ralit.
- The relay stores bounded exception bytes in CAS, appends local observations,
  and excludes those records from accepted-message replay. Source: DI-todav.
- Ex2's canonical topology is multi-relay peer collaboration; browser and
  Neovim are embodiments, and one-host multi-relay runs simulate nodes. Source:
  DI-nilas.
- Existing service tests cover browser/Neovim interoperability, accepted
  message replay, and peer ingestion. New evidence tests cover malformed,
  invalid-proof, and no-supported-handler records.
- Alice updates a local spec, Bob runs a second relay, and Mallory sends an
  unsupported or invalid envelope.

## Alternatives

### A. Layered focused coverage plus existing interoperability tests

Add a focused source-derived pCID inventory test beside protocol/service code;
extend focused evidence tests only where they prove CAS/observation/replay
separation; retain existing browser/Neovim and peer-relay interoperability
tests as the cross-node proof.

### B. One full end-to-end topology test for every claim

Run a multi-relay/browser/Neovim topology for inventory consistency, each
exception class, replay, and interoperability.

### C. Documentation-only pCID inventory

Rely on manual review of published pCIDs and leave current behavior tests as
the only automated evidence.

## Scenario analysis

### Local spec revision

Alice changes one local protocol document. Under A, a direct test derives its
pCID and checks the architecture inventory and scope declaration, failing at
the published-claim boundary. B starts a topology that cannot make a stale
Markdown inventory failure clearer. C permits documentation drift.

### Malformed and unsupported ingress

Mallory sends invalid bytes or a valid future pCID. A keeps assertions focused:
raw CAS retention, one local observation, and no accepted feed/replay entry.
B adds ports, polling, and timing to an observer-local assertion. C leaves
DI-todav's behavioral claim without regression protection.

### Browser/Neovim and two-relay collaboration

Bob uses Neovim while Alice uses a browser and their relays exchange traffic.
Existing interoperability tests already prove the relevant embodiment and
peer-ingestion path. A preserves them as the cross-node layer rather than
reimplementing those scenarios for every low-level assertion. B duplicates
slow, failure-prone topology setup. C gives no inventory protection.

### Restart and long-horizon evolution

An accepted relay restarts and rebuilds projections from `message-log.jsonl`.
Exception observations must not become replay input. A adds direct tests that
keep those paths separate while existing accepted-message replay coverage
remains intact. B can expose a restart problem but makes the cause hard to
localize. C risks future accidental coupling.

### Scale and maintenance

A gives exact failures close to source documents and storage behavior while
paying topology cost only for existing interoperability evidence. B raises suite
cost and flake risk. C is cheapest but does not make published claims durable.

## Conclusions

C is rejected because source-derived pCID publication is a concrete reader
claim that can drift. B is rejected as the sole strategy because inventory and
local evidence failures do not require a full topology to be meaningful.

A survives and is recommended: add focused source-derived inventory coverage;
retain focused CAS/observation/replay-separation assertions; and rely on the
existing browser/Neovim plus peer-relay interoperability tests for cross-node
behavior. This keeps relay-local evidence from becoming a fake global test
authority.

## Decisions still requiring DF

1. **Coverage strategy:** use layered focused coverage plus existing
   interoperability tests (recommended), full topology for every claim, or
   documentation-only inventory?
2. **Inventory authority:** derive each pCID from `protocols/*.md` and compare
   it with both `docs/architecture.md` and `CHANGELOG.md` (recommended), or
   compare only runtime protocol constants?
3. **Exception assertion:** test raw CAS retention, one observation per receipt,
   and exclusion from accepted peer/replay paths (recommended), or test only
   observation-file contents?
4. **Cross-node scope:** retain existing interoperability tests as the
   decentralized topology coverage (recommended), or add a new duplicate
   multi-relay/browser/Neovim scenario?

## Implications for open work

- `sojot.3` can proceed after these DFs are locked in `TODO-sojot`.
- The required `docs/testing.md` guide remains `sojot.4` work and should link
  these focused and interoperability layers after they exist.
