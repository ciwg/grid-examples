# Ex2 final README alignment

TE ID: TE-mahot

## Status

decided

## Decision status

Locked by DI-bafar: choose Alternative A. Add a concise README orientation
section, provide isolated one- and two-relay commands, point readers to the
selected relay root and canonical detail, and link the testing guide. This TE
does not change Ex2 runtime behavior or protocol semantics.

## Decision under test

How should the Ex2 README make its local-draft scope, decentralized relay
topology, reproducible local operation, and inspectable relay evidence clear
without duplicating the architecture and testing guides or promoting a local
demo into a global PromiseGrid trust claim?

## Assumptions and trust model

- Ex2 carries four source-derived, repo-local draft pCIDs. They are not frozen
  upstream specifications or automatic independent-peer interoperability
  claims. Source: DI-ralit.
- The canonical topology is multi-relay peer collaboration. Browser and
  Neovim are local embodiments; a single host can run one logical node or
  simulate several nodes with separate processes, keys, and roots. Source:
  DI-nilas.
- The relay records accepted and exception artifacts locally. Exception
  observations are not accepted replay input and do not settle another
  participant's intent. Source: DI-todav.
- `docs/architecture.md` owns the detailed architecture and storage model;
  `docs/testing.md` owns test commands and their exact guarantees. Source:
  DI-bubab.
- Alice is a contributor running the local demo. Bob runs another relay.
  Mallory may mistake a passing local test, one relay's observation, or a
  display label for broader trust evidence.

## Alternatives

### A. Add a concise README orientation section with links to canonical detail

Add one compact "Scope, topology, and evidence" section near Quick Start,
replace the stale final test snippet with a link to the testing guide, and add
a copyable isolated-run/two-relay example. Link to the scope declaration,
architecture, and testing guide for detail.

### B. Keep the README as an application walkthrough only

Leave scope, topology, evidence, and complete verification in the linked
documents. Correct only the final test command if needed.

### C. Copy the architecture and testing guidance into the README

Put pCID inventory, non-claims, detailed storage semantics, every test layer,
and the full Docker discussion directly in the README.

## Scenario analysis

### Alice starts a local browser demo

Under A, Alice gets an explicit `--data-root` example and learns that the
resulting artifacts belong to her relay node. B starts the demo but leaves the
runtime root and evidence boundary implicit. C works initially but buries the
simple startup path in repeated reference material.

### Bob joins through a second relay

A shows separate relay roots and the `--peer` relationship, then links to the
topology detail. It makes clear that two relay processes represent distinct
logical nodes even on one host. B relies on a reader finding that distinction
elsewhere. C repeats the architecture diagram and risks later divergence.

### Mallory sends malformed or unsupported input

A gives the reader one durable pointer: inspect the selected relay data root,
then follow the architecture/testing links for CAS and `observations.jsonl`
semantics. B makes inspection discoverability poor. C copies nuanced evidence
rules into two locations and risks treating an observation as global proof
when one copy evolves.

### Local profile revision

Alice changes a local profile. A points directly to the scope declaration and
testing guide, which explain source-derived pCIDs and verification. B makes
the scope declaration less discoverable from the primary entry point. C makes
the README a second inventory that can drift from the canonical one.

### Long-horizon maintenance

A preserves a short orientation layer while each detailed guide retains its
single purpose. B minimizes edits but leaves the README's final `go test`
snippet inconsistent with the complete verification guide. C increases
duplication and the chance of a stale reader-facing contract.

## Conclusions

C is rejected because copying detailed architecture and testing guidance into
the README creates competing sources of truth. B is rejected because the
README is the primary entry point and currently understates both verification
and evidence/topology boundaries.

A survives and is recommended: add a concise orientation section, link the
canonical detailed documents, provide a reproducible isolated local and
two-relay command example, and replace the final one-command test snippet
with the testing-guide link.

## Decisions still requiring DF

1. **README structure:** add one concise "Scope, topology, and evidence"
   section linked to the canonical guides (recommended), make only the stale
   test correction, or duplicate detailed guidance in the README?
2. **Operation example:** show `--data-root` for an isolated relay plus a
   second relay with a separate root and `--peer` (recommended), show only the
   existing default-root quick start, or refer only to Docker?
3. **Evidence guidance:** state that a user can inspect the selected relay
   root and link to the architecture/testing guides for artifact semantics
   (recommended), list exact internal artifact files in the README, or omit
   evidence inspection from the README?
4. **Verification navigation:** replace the final `go test ./...` snippet with
   a link to `docs/testing.md` (recommended), retain the snippet and add the
   guide link beside it, or remove the Tests section entirely?

## Implications for open work

- `sojot.5` can proceed once the four DFs are recorded in the Ex2 alignment
  TODO with a new DI.
- The resulting README changes are documentation-only and must preserve the
  detailed architecture, scope, and testing documents as canonical sources.
