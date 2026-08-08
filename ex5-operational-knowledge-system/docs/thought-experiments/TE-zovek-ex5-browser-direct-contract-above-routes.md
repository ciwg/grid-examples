# Browser Direct Contract Above Routes

TE ID: TE-zovek
## Status
decided

## Decision under test

What the next typed browser direct-contract slice should be now that the
Chrome/Chromium native-messaging embodiment is shipped, browser inspect/search
reads have their first typed `operation` path, but many browser semantics still
cross the native host as generic `type:"request"` plus `method + path`.

This TE corresponds to TODO `rumav.1` / `rumav.2` / `rumav.3`.

## Assumptions

- Alice uses the browser as the main review, authoring, create, and operate
  surface.
- Bob and Carol already use direct local contracts in CLI and Neovim.
- Browser live drafting is already on the direct browser contract and is out of
  scope except where request/response workflows interact with it.
- The browser must stay fail-closed; no silent fallback to the older HTTP
  browser path is allowed.
- The current browser contract already has typed reads for:
  - `inspect_item`
  - `inspect_run`
  - `inspect_entity`
  - `search`
  - `problem_review`
- The remaining route-shaped browser traffic is concentrated in:
  - startup/catalog reads such as dashboard, places, resources,
    responsibilities, items, and runs
  - create/operate writes such as create place/resource/responsibility/item,
    record run, add evidence, approvals, revisions, and supersede

## Alternatives

### A. Keep the remaining browser request/response traffic generic for now

Do no further typed browser work yet. Keep browser list/dashboard reads and all
browser writes on generic `type:"request"` forwarding over the native host.

### B. Move the browser startup/catalog slice next

Add typed browser operations for the main non-detail read surfaces first:

- `dashboard`
- list places/resources/responsibilities/items/runs
- any supporting catalog/bootstrap reads needed for browser startup

Create/operate mutations would remain generic request forwarding temporarily.

### C. Move the browser create/operate slice next

Add typed browser operations for the main mutation workflows first:

- create place/resource/responsibility/item
- record run
- add evidence
- record approvals
- add revision
- supersede item

Catalog/bootstrap reads would remain generic request forwarding temporarily.

### D. Move both the browser catalog and create/operate request/response slice
together

Add typed operations for both the startup/catalog reads and the main
create/operate mutations, leaving only a much smaller compatibility residue in
generic request forwarding.

## Scenario analysis

### Scenario 1: normal operator startup and daily browser use

Alice opens the browser, reviews current work, creates or updates records, and
approves results.

#### A. Keep the remaining browser traffic generic

What it makes easier:
- no migration work
- current browser behavior stays stable

What it makes harder:
- browser still depends heavily on route-shaped semantics even after the direct
  embodiment move
- the browser contract remains only partly more runtime-shaped than the old
  HTTP adapter

New obligations:
- continue explaining why a “direct” browser embodiment still tunnels so much
  adapter vocabulary

#### B. Move the browser startup/catalog slice next

What it makes easier:
- browser startup and review-mode browsing become more runtime-shaped
- the browser’s first visible state comes from typed runtime operations
- list/detail semantics start to align better

What it makes harder:
- create and operate workflows still drop back into generic request forwarding
- the highest-value user mutations remain adapter-shaped

New obligations:
- define typed list/catalog operations clearly

#### C. Move the browser create/operate slice next

What it makes easier:
- the most important durable writes become runtime-shaped first
- PromiseGrid pressure moves off route-shaped mutation semantics where it
  matters most for durable history

What it makes harder:
- startup/catalog reads still look adapter-shaped
- browser directness remains visually mixed because loading and mutation use
  different semantic styles

New obligations:
- define typed mutation operations with clear payloads and validation

#### D. Move both browser catalog and create/operate together

What it makes easier:
- request/response behavior becomes much more uniformly runtime-shaped
- the browser starts and works inside one clearer contract family
- remaining generic forwarding becomes an explicit small residue rather than
  the bulk of the browser path

What it makes harder:
- the migration slice is broader
- more browser tests and service operations must move together

New obligations:
- define a coherent first browser operation set across read and write paths

### Scenario 2: failure handling and operator trust

Alice hits a malformed browser request, a validation failure, or a stale input
while using the direct browser embodiment.

#### A. Keep the remaining browser traffic generic

What it makes easier:
- current errors stay familiar

What it makes harder:
- failures are still often framed in route/handler terms
- the browser direct contract keeps inheriting HTTP-shaped validation seams

#### B. Move the browser startup/catalog slice next

What it makes easier:
- typed read failures become clearer and narrower

What it makes harder:
- write failures still speak adapter-shaped semantics

#### C. Move the browser create/operate slice next

What it makes easier:
- durable write validation happens at a cleaner operation boundary
- browser errors for creation and approvals can become more domain-shaped

What it makes harder:
- catalog/bootstrap errors still come through the adapter-shaped path

#### D. Move both together

What it makes easier:
- most request/response failures shift to runtime-operation semantics
- browser trust improves because startup and mutation both speak the same
  contract family

What it makes harder:
- larger regression surface in one wave

### Scenario 3: mixed-version rollout

Dave upgrades the runtime and extension while some local docs/tests still
assume the older generic request-forwarding path.

#### A. Keep the remaining browser traffic generic

What it makes easier:
- no rollout churn

What it makes harder:
- no architectural progress

#### B. Move the browser startup/catalog slice next

What it makes easier:
- rollout is additive and browser-visible quickly

What it makes harder:
- generic write forwarding still has to stay around, so the contract remains
  half-migrated

#### C. Move the browser create/operate slice next

What it makes easier:
- the highest-value durable writes move first

What it makes harder:
- rollout docs must explain why reads still look old while writes moved

#### D. Move both together

What it makes easier:
- fewer mixed contract stories remain afterward

What it makes harder:
- more code and tests change at once

### Scenario 4: long-horizon PromiseGrid alignment

Steve wants the browser embodiment to move closer to the strongest direct local
contract story practical in `ex5`.

#### A. Keep the remaining browser traffic generic

What it makes easier:
- short-term stability only

What it makes harder:
- stalls the main remaining browser contract impurity

#### B. Move the browser startup/catalog slice next

What it makes easier:
- improves the “shape” of browser review and startup

What it makes harder:
- still leaves durable writes adapter-shaped, which is a deeper PromiseGrid
  impurity than catalog reads

#### C. Move the browser create/operate slice next

What it makes easier:
- shifts the main durable action path above route-shaped semantics first
- better matches the PromiseGrid instinct that durable intent should not hide
  inside adapter route names

What it makes harder:
- startup/catalog still lag behind, so the browser is not yet broadly uniform

#### D. Move both together

What it makes easier:
- gives the browser one much clearer direct request/response contract
- aligns both startup review and durable writes in the same wave

What it makes harder:
- broader step than the minimum needed to improve PromiseGrid alignment

### Scenario 5: relationship to future substrate extraction

Frank wants later extraction work to have a cleaner boundary between reusable
runtime operations and adapter residue.

#### A. Keep the remaining browser traffic generic

What it makes easier:
- nothing

What it makes harder:
- browser still depends on route-shaped seams that are less extractable

#### B. Move the browser startup/catalog slice next

What it makes easier:
- read-model operations become more extractable

What it makes harder:
- write-side semantics still remain buried in route strings

#### C. Move the browser create/operate slice next

What it makes easier:
- durable mutation intents become explicit runtime operations earlier

What it makes harder:
- catalog-side extraction remains postponed

#### D. Move both together

What it makes easier:
- yields the strongest substrate-ready browser boundary short of a full rewrite

What it makes harder:
- increases the size of the immediate implementation wave

## Conclusions

Rejected:

- Alternative A: too weak; it leaves the strongest remaining browser impurity
  untouched.
- Alternative B: helpful, but it improves the catalog side before the more
  important durable mutation side.

Surviving:

- Alternative C: move the browser create/operate slice next
- Alternative D: move both browser catalog and create/operate together

Alternative C is the safest surviving option.

Alternative D is the most PromiseGrid-aligned surviving option because it
raises both startup/catalog and durable write behavior above route-shaped
forwarding in one coherent browser request/response wave instead of leaving one
half behind.

## Implications for open TODOs and pending DIs

- TODO `135` should lock whether this wave is mutation-first or a broader
  browser request/response slice.
- If `135` chooses the broader slice, the remaining generic browser request
  forwarding becomes a much smaller compatibility residue.
- The remaining DF question is whether to take only the create/operate slice or
  both the catalog and create/operate slice together.

The locked DF result is:

- `135C`
- typed browser write operations grouped one-per-workflow family
- all browser create/operate write operations move together in this slice
