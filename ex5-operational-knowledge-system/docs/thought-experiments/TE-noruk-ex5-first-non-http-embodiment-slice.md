# ex5 first non-HTTP embodiment slice

TE ID: TE-noruk
## Status
decided

## Decision under test

What the first direct non-HTTP embodiment contract for
`ex5-operational-knowledge-system` should be now that the repo already ships:

- eight frozen signed families
- origin-aware ongoing peer exchange
- dedicated remote relay deployment
- websocket-preferred live drafting over the local adapter
- a clean current-scope statement that still says browser, CLI, and Neovim
  project through `local_http`

The question is not whether `ex5` should delete HTTP immediately. The question
is which first non-HTTP embodiment slice most honestly advances the runtime
without reopening the finished durable-family and relay work.

## Assumptions

- Browser, CLI, and Neovim currently share one local runtime and one local HTTP
  adapter.
- Durable PromiseGrid behavior is already anchored below that adapter in the
  frozen-family, CAS, peer-exchange, and relay layers.
- Mallory may observe local transport or replay local requests, but cannot
  forge signed family envelopes.
- Mixed-version rollout matters because current users and tests already depend
  on the local HTTP adapter.
- This TE is about direct embodiment contracts, not about replacing the remote
  relay protocol or changing the durable family boundaries.

## Alternatives

### Alternative A: terminal-first direct local contract

Introduce a direct non-HTTP local embodiment contract first for CLI and
Neovim, while the browser continues using the current local HTTP adapter.

The new contract would expose runtime operations over a local-only transport
such as a Unix domain socket or stdio-framed local RPC. Browser routes remain
as compatibility and embodiment surface for now.

### Alternative B: one shared non-HTTP local runtime socket for every embodiment

Introduce one new local runtime socket contract and move browser, CLI, and
Neovim toward it together. HTTP becomes compatibility only, and browser
communication is bridged through a helper process or sidecar rather than going
directly to the current local adapter.

### Alternative C: bypass the local runtime boundary and let embodiments speak
remote relay contracts directly

Treat the remote relay feed/blob protocols as the next embodiment contract and
have CLI, Neovim, or browser operations talk to relay-visible transport
directly for most work, minimizing the distinct local runtime adapter surface.

## Scenario analysis

### Scenario 1: normal local authoring and review

Alice uses the browser to review, Bob uses CLI to inspect queues, and Carol
uses Neovim for live drafting.

Alternative A improves the terminal embodiments directly without destabilizing
the browser. CLI and Neovim gain a non-HTTP contract closest to the runtime,
while browser workflows continue unchanged. This is a partial step, but it
keeps the first slice operationally narrow.

Alternative B is conceptually cleaner because all embodiments converge on one
new local contract. But it forces browser transport redesign immediately, even
though browser constraints are very different from CLI and Neovim.

Alternative C is the least coherent for normal authoring because relay
transport is about signed durable exchange and blob carriage, not about local
editing, search projections, or live draft state.

### Scenario 2: live drafting and presence

Dave and Ellen co-edit one knowledge item; Frank watches presence in Neovim.

Alternative A can move the terminal side closer to the runtime while leaving
browser live drafting on websocket-over-HTTP. That means mixed embodiment
carriage still exists, but only where the browser actually needs it.

Alternative B could unify live state around one non-HTTP local contract. That
is elegant if it works, but it creates a much larger first cut because browser
live collaboration must now tunnel through something other than the already
working websocket adapter.

Alternative C is a mismatch. The relay is durable transport, not local shared
draft presence.

### Scenario 3: durable operations versus embodiment projection

Grace records runs, uploads evidence, and approves work from both browser and
terminal surfaces.

Alternative A keeps the direct non-HTTP contract concentrated on the
embodiments that can use a local runtime-oriented interface most naturally.
The durable layers below stay unchanged, and the browser remains an adapter
projection.

Alternative B also preserves the durable layers, but it demands a much broader
adapter rewrite at the same time.

Alternative C risks conflating embodiment actions with relay-fed durable
history. That makes local workflow semantics depend too directly on transport
designed for peer exchange.

### Scenario 4: mixed-version rollout

Heidi upgrades the CLI first, while Ivan still uses the browser and Judy still
uses the current Neovim plugin.

Alternative A is easiest to roll out incrementally. The new contract can land
next to HTTP, and each terminal embodiment can adopt it independently.

Alternative B is much harder because browser migration becomes part of the
first rollout wave. The compatibility burden is higher across tests, docs, and
runtime launch flows.

Alternative C is awkward because relay transport is already a separate axis of
change. Tying embodiment rollout to relay semantics increases operational
coupling.

### Scenario 5: long-horizon PromiseGrid alignment

The repo later wants a stronger statement that embodiments are not defined by
HTTP route names.

Alternative A moves that statement forward concretely for the terminal
embodiments without forcing the browser into an unnatural first rewrite. It
creates a staged path: browser may remain an adapter longer, but CLI and
Neovim stop depending on HTTP route naming.

Alternative B is the purest local-contract end state if all embodiments can
converge there. It most strongly eliminates HTTP as the defining embodiment
surface.

Alternative C does not solve the right problem. It may reduce HTTP use, but it
does so by collapsing embodiment concerns into relay concerns.

### Scenario 6: implementation and testing scope

Karen needs the first `117` slice to be real, testable, and not another year
of infrastructure churn.

Alternative A is the narrowest viable cut. The runtime can expose a direct
local contract for terminal clients without redesigning browser integration and
without touching the finished durable-family and relay layers.

Alternative B is stronger architecturally, but it is a much larger first wave.
Browser support, local helper processes, and new integration harnesses all
arrive at once.

Alternative C is tempting only if the team wants to collapse local operations
into relay semantics immediately, which does not match the current runtime
shape.

## Conclusions

Rejected:

- Alternative C. It confuses embodiment transport with relay transport and
  does not fit live drafting or local projection behavior.

Surviving:

- Alternative A: terminal-first direct local contract
- Alternative B: one shared non-HTTP local runtime socket for every embodiment

Recommendation:

- Alternative A

Why:

- It is the clearest first non-HTTP embodiment slice that does real work
  without reopening the finished browser and relay layers all at once.
- It removes HTTP-route dependence from the embodiments that can most naturally
  consume a local runtime contract now: CLI and Neovim.
- It preserves a path toward later browser migration if that still proves
  worthwhile.

## Implications for TODOs and pending DIs

- TODO `117` should lock either:
  - `A`: terminal-first non-HTTP embodiment contract for CLI and Neovim, or
  - `B`: one shared non-HTTP local runtime socket for browser, CLI, and Neovim
- If `A` is chosen, the next DF should narrow:
  - exact local transport shape
  - whether CLI and Neovim both move in the same first slice or one leads
  - how browser compatibility remains described during the transition
- If `B` is chosen, the next DF should narrow:
  - browser-side helper/bridge shape
  - how websocket live drafting maps into the new contract
  - how existing browser automation migrates
