## Title

ex5 runtime contract above adapters

## TE ID

TE-zoruk

## Status

needs DF

## Decision under test

What the next `ex5` runtime-contract slice above adapter-shaped seams should
be, now that CLI and Neovim already have a direct local socket transport but
that transport still forwards generic `method + path` requests through the HTTP
handler.

## Assumptions

- The current shipped `ex5` scope is already strong enough that this is a
  contract-shape refinement, not a recovery from broken durable semantics.
- Browser remains on the local HTTP adapter in this wave.
- CLI and Neovim already have a local Unix-socket transport, but their
  request/response path still uses route-shaped semantics instead of a typed
  runtime contract.
- Live drafting over the local socket is already more runtime-shaped than the
  generic request/response path because it carries item-level live state rather
  than raw HTTP methods and paths.
- The user asked for the most PromiseGrid-aligned path, but also wants the
  system to keep working and stay well tested.

## Alternatives

### Alternative A

Keep the current generic socket request envelope:

- local socket stays `type=request`, `method`, `path`, `headers`, `body`
- socket server continues forwarding into the HTTP handler
- improvements happen only in docs and metadata

### Alternative B

Introduce a first typed local runtime contract for terminal embodiments while
keeping HTTP as an adapter:

- add typed local operations for the highest-value non-live workflows
- CLI and Neovim call those operations over the local socket instead of
  constructing route-shaped requests
- HTTP server remains as the browser adapter and compatibility surface
- generic route-shaped socket forwarding can remain temporarily for the
  untouched operations

### Alternative C

Replace both the local socket and HTTP handlers with one shared typed runtime
operation layer immediately:

- define one internal operation table
- rewrite socket and HTTP adapters together to call it for all workflows
- move the whole runtime surface above route-shaped seams in one wave

## Scenario analysis

### Scenario 1: normal operator use

Alice uses browser, CLI, and Neovim against one local runtime. Browser creates
records, CLI performs review work, and Neovim edits and approves items.

- Alternative A keeps behavior stable but leaves the direct terminal contract
  semantically thin. The “direct socket” still behaves like tunnelled HTTP.
- Alternative B improves the terminal contract where PromiseGrid pressure is
  already highest: the non-browser embodiments stop depending on route strings
  for the first chosen workflows.
- Alternative C yields the cleanest architecture eventually, but it couples too
  many surfaces into one immediate rewrite.

What gets easier:

- A: no migration work
- B: terminal contract becomes meaningfully more runtime-shaped without moving
  the browser yet
- C: one unified internal model

What gets harder:

- A: PromiseGrid alignment stalls at the seam that now matters most
- B: mixed contract surface during transition
- C: much larger regression surface

### Scenario 2: failure and corruption handling

Bob hits a malformed request or a partially upgraded embodiment during local
socket use.

- Under A, errors still surface in HTTP-route terms. That is practical, but it
  means the local direct contract still inherits route parsing and handler
  assumptions from the browser adapter.
- Under B, typed operations can validate arguments at the runtime-operation
  boundary, reducing dependence on route parsing for the adopted workflows while
  leaving compatibility behavior intact for untouched ones.
- Under C, error handling can be most uniform eventually, but getting there in
  one wave means more paths change at once, which increases the chance of new
  edge regressions.

What gets easier:

- A: current failures stay familiar
- B: adopted operations get cleaner validation and clearer error semantics
- C: eventual uniformity

What gets harder:

- A: adapter-shaped failures remain part of the “direct” contract
- B: two styles of local socket operation coexist during rollout
- C: large correctness burden in one step

### Scenario 3: mixed-version rollout

Carol upgrades the runtime first. Dave is still on an older CLI or Neovim
client.

- Alternative A is easiest for mixed versions because nothing changes.
- Alternative B allows additive rollout if the new typed operations live beside
  the older route-shaped socket requests temporarily. New clients can adopt the
  stronger contract while old ones still function.
- Alternative C is hardest because all embodiments and both adapters need to
  move together or accept a bigger compatibility layer anyway.

What gets easier:

- A: zero rollout friction
- B: staged adoption with explicit compatibility
- C: eventual consistency if rollout succeeds perfectly

What gets harder:

- A: no architectural gain
- B: temporary dual path inside the local socket service
- C: rollout coordination becomes heavy

### Scenario 4: trust-boundary clarity

Ellen asks which surface is the true runtime contract versus which surfaces are
adapters.

- Under A, the answer remains muddy: the browser adapter is explicit, but the
  “direct” terminal socket still tunnels route names and HTTP-shaped semantics.
- Under B, the answer becomes clearer: adopted socket operations are true local
  runtime contract, while HTTP remains the browser adapter and untouched socket
  requests remain transitional compatibility.
- Under C, the answer is clearest after the rewrite, but the cost of getting
  there is much higher.

What gets easier:

- A: no new docs needed beyond honesty
- B: a real runtime-vs-adapter distinction appears in code
- C: strongest conceptual cleanliness

What gets harder:

- A: PromiseGrid claims about terminal directness remain partially overstated
- B: docs must explain the transition line
- C: implementation scope expands sharply

### Scenario 5: long-horizon evolution

Frank wants this wave to prepare for a later browser non-HTTP embodiment and a
possible generalized runtime substrate.

- Alternative A contributes little to those later steps.
- Alternative B creates a reusable typed runtime-operation layer in pieces,
  which is a strong staging ground for both browser non-HTTP work and any later
  substrate extraction.
- Alternative C may get to the end-state faster in theory, but only if the repo
  can absorb a much broader rewrite now.

What gets easier:

- A: immediate stability only
- B: incremental path toward `132` and `133`
- C: maximum end-state purity if successful

What gets harder:

- A: future waves have to reopen the same seam
- B: requires deciding which operations move first
- C: risks turning `131` into an oversized umbrella rewrite

## Conclusions

Rejected:

- Alternative A: too weak; it leaves the main adapter-shaped seam untouched.
- Alternative C: too large for the intended “safe next step” after `129` and
  `130`.

Surviving:

- Alternative B: add a first typed local runtime contract for selected terminal
  workflows while keeping HTTP as an adapter and preserving backward
  compatibility during rollout.

Alternative B is the most PromiseGrid-aligned remaining path because it moves
real behavior above route-shaped seams without pretending the browser adapter
has already been replaced. It is also the strongest staged move toward `132`
and `133`.

## Implications for open TODOs and pending DIs

- TODO `131` should lock a first typed local runtime-operation slice rather
  than a full adapter rewrite.
- TODO `132` becomes easier after `131` if the browser later adopts the same
  typed runtime-operation layer instead of jumping directly from HTTP routes.
- TODO `133` also becomes easier because the runtime-operation layer would be
  more extractable than today's route-shaped socket forwarding.
- The remaining DF question is which first workflow slice should move from
  route-shaped socket forwarding into typed local runtime operations.

