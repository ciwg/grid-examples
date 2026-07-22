# Ex5 Place/Resource Family Boundary

TE ID: `TE-puvok`
## Status
needs DF

## Decision under test

How `ex5` should make the remaining operational context references
peer-visible: as two separate frozen families for places and resources, or as
one combined context family.

Related TODO:

- `110` - `ex5-operational-knowledge-system/TODO/TODO-pivok-ex5-peer-visible-place-resource-families.md`

## Assumptions

- `ex5` already exchanges six frozen signed families:
  `knowledge-item`, `knowledge-approval`, `knowledge-evidence`,
  `operational-run`, `knowledge-link`, and `knowledge-responsibility`.
- Peer-visible entities already use create-envelope CIDs as durable IDs, while
  old short IDs remain aliases.
- Runs and links can already preserve `place` and `resource` references, but
  those references are still reported as unresolved across peers.
- The user wants the closest practical move toward a fully on-grid `ex5`, not
  a temporary compatibility-only patch.

## Scope and systems affected

- `ex5-operational-knowledge-system/protocols/*`
- service replay, signed envelope helpers, peer exchange, and CAS authority
- browser, CLI, and Neovim drilldowns that present place and resource context
- PromiseGrid claims, architecture, and peer-exchange staging docs

## Alternatives

### Alternative A

Freeze two separate families:

- `operational-place`
- `operational-resource`

Places keep hierarchy and naming context. Resources keep typed operational
things and continue to reference places.

### Alternative B

Freeze one combined family:

- `operational-context`

Both places and resources become one typed durable family with a shared
identity space and one pCID.

### Alternative C

Do not freeze either family yet. Keep place/resource references outside the
peer-visible slice and preserve them only as unresolved context.

## Scenario analysis

### Scenario 1: normal peer exchange between Alice and Bob

Alice exports a run recorded in a receiving bay with two resources: a dock
scanner and a pallet jack. Bob imports the run and later inspects its context.

Alternative A:

- Bob receives the run plus first-class place and resource artifacts.
- The run can point at a place family record and resource family records
  directly.
- The resource-to-place relationship remains explicit and typed.

Alternative B:

- Bob also receives the run plus first-class context artifacts.
- The distinction between places and resources lives inside one family rather
  than in separate pCIDs.
- Import can still work, but all context semantics now share one durable
  namespace and one contract.

Alternative C:

- Bob still gets an imported run with unresolved references.
- The context is visible only as opaque IDs or missing related lookups.
- The run remains portable, but the surrounding operational graph is not.

Result:

- A and B both solve the actual peer-visible gap.
- C does not.

### Scenario 2: mixed-version nodes and migration pressure

Alice upgrades first. Bob and Carol upgrade later. Historical runs and links
already contain place/resource references in current local forms.

Alternative A:

- Migration can be staged by family.
- Alice can freeze places first or resources first if later TE work shows an
  internal dependency, but the default still keeps the semantic split clear.
- Mixed-version reasoning stays understandable because each family has one job.

Alternative B:

- Migration happens as one large context-family step.
- That reduces the number of pCIDs, but it enlarges the migration blast
  radius.
- Any schema changes for one context type now perturb the shared family.

Alternative C:

- No migration work is needed immediately.
- But the unresolved-reference limitation persists and keeps blocking the
  stricter on-grid story.

Result:

- A is easier to stage safely.
- B is workable, but broader and less decomposable.
- C avoids work but preserves the gap.

### Scenario 3: long-horizon protocol evolution

Steve later wants richer place hierarchy semantics, resource lifecycle
attributes, or resource-specific operational metadata.

Alternative A:

- Place and resource semantics can evolve independently.
- New obligations for one family do not automatically expand the other.
- Link and run references remain typed at the family boundary.

Alternative B:

- A unified context family can still express type differences in payload.
- But semantic drift accumulates inside one shared contract.
- One family becomes responsible for two conceptually different models:
  physical context hierarchy and operational thing inventory.

Alternative C:

- There is no new family-level design burden because the gap remains open.
- But every later feature still has to work around missing peer-visible
  context.

Result:

- A creates cleaner long-horizon protocol boundaries.
- B centralizes context but risks a bloated family.

### Scenario 4: trust and verification boundaries

Mallory sends forged or malformed context histories to Bob.

Alternative A:

- Verification is split by family and easier to reason about.
- Bob can reject malformed resource history without conflating it with place
  history.
- CAS and replay authority can advance family by family.

Alternative B:

- Verification remains possible, but all context validation rules now sit in
  one family.
- This increases parser and verifier surface area per family.

Alternative C:

- The current unresolved-reference model avoids new verification logic.
- But it also avoids solving the missing peer-visible trust surface.

Result:

- A has the clearest verification boundary.
- B is broader with no compensating trust advantage.

### Scenario 5: embodiment impact

Alice uses the browser, Bob uses the CLI, and Carol uses Neovim. All three
inspect the same exchanged run.

Alternative A:

- Each embodiment can continue to say “place” and “resource” explicitly.
- Existing drilldowns map naturally onto separate durable families.
- UI/CLI copy remains straightforward.

Alternative B:

- The runtime can still project place/resource views.
- But the underlying family name and documentation become less aligned with
  what the embodiments actually present.

Alternative C:

- Embodiments keep showing unresolved or partially local-only context.

Result:

- A aligns best with current product language and operator expectations.

## Conclusions

Rejected alternative:

- Alternative C: it leaves the main strict PromiseGrid gap open.

Surviving alternatives:

- Alternative A: freeze separate `operational-place` and
  `operational-resource` families.
- Alternative B: freeze one combined `operational-context` family.

## Implications and future work

- If Alternative A is locked, the next follow-on question is likely ordering:
  whether `operational-place` or `operational-resource` should freeze first, or
  whether they land as one grouped batch with disjoint signed logs.
- If Alternative B is locked, the next follow-on question is the exact payload
  type system and how one family distinguishes hierarchy-bearing places from
  resource records cleanly.
- Either surviving alternative would then need a new DI plus the corresponding
  protocol docs, signed envelope helpers, peer-exchange coverage, and claims
  updates under TODO `110`.

## Decision status

Needs DF:

1. Alternative A: separate `operational-place` and `operational-resource`
   families.
2. Alternative B: one combined `operational-context` family.
