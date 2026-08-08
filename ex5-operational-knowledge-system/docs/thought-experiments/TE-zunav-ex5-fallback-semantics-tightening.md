## Title

ex5 fallback semantics tightening

## TE ID

TE-zunav

## Status

decided

## Decision under test

How `ex5` should tighten the remaining compatibility fallback semantics for the
browser and Neovim embodiments after the CLI already moved to fail-closed
local-socket behavior with explicit HTTP opt-in.

## Assumptions

- The current shipped `ex5` scope is already PromiseGrid-aligned enough that
  this is a refinement wave, not a repair of broken durable semantics.
- Browser live drafting currently prefers websocket carriage under the local
  HTTP adapter and falls back to the HTTP live route.
- Neovim currently prefers the local Unix socket, then falls back to websocket
  over the local HTTP adapter, then finally to the HTTP live route.
- CLI transport policy is already locked: direct local socket by default,
  explicit HTTP compatibility only through operator opt-in.
- Durable signed families, relay feed behavior, and CAS authority are out of
  scope for this decision.

## Alternatives

### Alternative A

Keep the current remaining implicit fallback behavior:

- browser: websocket over local HTTP -> HTTP live route
- Neovim: local Unix socket -> websocket over local HTTP -> HTTP live route

### Alternative B

Tighten only the cross-adapter fallback:

- browser keeps implicit websocket -> HTTP fallback because both lanes remain
  inside the same local HTTP embodiment adapter
- Neovim stops silently crossing from the local Unix-socket embodiment contract
  into HTTP-adapter compatibility transport; websocket/HTTP compatibility
  becomes explicit opt-in

### Alternative C

Fail closed everywhere unless compatibility transport is explicitly selected:

- browser no longer falls back from websocket to HTTP live route automatically
- Neovim no longer falls back from local socket to websocket or HTTP

## Scenario analysis

### Scenario 1: normal operation on one local runtime

Alice runs the `ex5` browser and Neovim against a healthy local runtime. The
browser upgrades live drafting to websocket immediately. Neovim resolves the
runtime socket from `/api/meta` and uses the direct local socket.

- Alternative A is easy here because nothing changes.
- Alternative B is also easy here. Normal operation keeps the currently
  preferred lanes and only changes failure behavior.
- Alternative C is also operationally clean during the happy path, but it
  creates no additional PromiseGrid value over B in this scenario because the
  browser is already using its primary adapter.

What gets easier:

- A: no migration work
- B: normal traffic stays unchanged while cross-adapter fallback becomes more
  honest
- C: maximum strictness

What gets harder:

- A: the actual embodiment contract remains partly ambiguous
- B: operators must learn one explicit compatibility switch for Neovim
- C: browser sessions become brittle even when only websocket upgrade fails

### Scenario 2: runtime restart or temporary transport failure

Bob is editing in Neovim while the runtime restarts. The local socket drops.
The browser websocket also drops during the restart window.

- Under A, Neovim silently demotes to websocket or HTTP. Browser silently drops
  to HTTP polling. Users keep working, but the embodiment contract truth is
  blurry: one terminal embodiment has crossed into adapter compatibility without
  explicit operator intent.
- Under B, browser behavior stays pragmatic because the browser still lives
  inside the local HTTP adapter either way. Neovim instead surfaces that the
  direct contract is unavailable unless the operator explicitly allows
  compatibility transport.
- Under C, both browser and Neovim stop live collaboration immediately unless
  compatibility was preselected. This is the purest interpretation, but it
  turns a local websocket upgrade failure into a visible browser outage even
  though the same embodiment adapter is still present.

What gets easier:

- A: resilience through silent demotion
- B: browser resilience without weakening the terminal contract
- C: no ambiguity anywhere

What gets harder:

- A: operators cannot easily tell when they lost the stronger embodiment path
- B: Neovim interruptions become explicit and need operator action
- C: browser experience degrades sharply during benign websocket churn

### Scenario 3: mixed-version nodes and staged rollout

Carol upgrades the runtime first, while Dave still has an older Neovim plugin
or browser tab. The browser websocket route is available, but some sessions may
still need compatibility behavior during rollout.

- A tolerates rollout easily, but it keeps long-term ambiguity.
- B preserves browser rollout tolerance while making terminal compatibility a
  conscious operator choice. This matches the repo's current direction: keep
  compatibility where it is still structurally useful, but stop pretending it
  is the primary contract.
- C forces every live path to be explicitly aligned during rollout, which is
  clean in principle but heavy for the browser because browser upgrade and local
  runtime restart windows are common.

What gets easier:

- A: maximum staged compatibility
- B: staged compatibility where it still matches the embodiment boundary
- C: strict deployment discipline

What gets harder:

- A: compatibility drift can persist indefinitely
- B: terminal rollout docs and UX need one more explicit mode
- C: rollout friction rises for browser sessions with little architectural gain

### Scenario 4: trust-boundary clarity

Ellen asks what embodiment contract each client is actually using so she can
reason about PromiseGrid alignment and operator expectations.

- A makes this hardest because Neovim may move from local socket to HTTP
  compatibility without explicit user intent.
- B is clearer: browser remains one local HTTP embodiment with an internal live
  transport fallback, while Neovim either stays on the local socket or is
  explicitly told to use compatibility transport.
- C is clearest in a narrow sense, but it treats the browser websocket upgrade
  as if it were a different embodiment boundary rather than a different live
  transport inside the same adapter.

What gets easier:

- A: nothing new to document
- B: embodiment boundaries become explicit and machine-readable
- C: absolute transport strictness

What gets harder:

- A: auditing and reasoning about the active contract
- B: one more explicit operator control for Neovim
- C: over-couples browser availability to websocket upgrade success

### Scenario 5: long-horizon evolution

Frank wants `ex5` to keep moving toward more direct embodiment contracts,
eventually including a future non-HTTP browser path.

- A slows that trajectory because compatibility remains implicitly normal.
- B matches the likely migration order: tighten cross-adapter fallback first,
  then later decide whether the browser should gain a more direct embodiment.
- C front-loads strictness in a place that may be discarded later if the browser
  gets a direct non-HTTP embodiment anyway.

What gets easier:

- A: short-term stability
- B: staged PromiseGrid tightening with low conceptual debt
- C: absolute compatibility minimization

What gets harder:

- A: future cleanup remains necessary
- B: still leaves browser compatibility logic to revisit later
- C: introduces browser strictness before the browser embodiment is redesigned

## Conclusions

Rejected:

- Alternative A: too much implicit compatibility remains, especially for
  Neovim's cross-adapter demotion.
- Alternative C: stricter than necessary for the browser because websocket and
  HTTP are still both inside the same local HTTP embodiment adapter.

Surviving:

- Alternative B: tighten only the cross-adapter fallback.

Alternative B is the most PromiseGrid-aligned remaining path because it draws
the boundary at the embodiment contract rather than at every transport detail.
Browser websocket-to-HTTP fallback stays inside the already-declared browser
adapter. Neovim no longer silently crosses from the direct local socket
contract into compatibility transport.

The locked DF result is:

- browser keeps implicit websocket-to-HTTP fallback inside `local_http`
- Neovim compatibility transport is explicit opt-in through
  `oks-nvim --socket=off`

## Implications for open TODOs and pending DIs

- TODO `130` should lock a per-embodiment fallback policy rather than a single
  global fallback rule.
- The remaining DF question is how Neovim should expose explicit compatibility
  mode once silent cross-adapter fallback is removed.
- TODO `131` becomes cleaner after this decision because runtime contract work
  will sit on top of a more honest embodiment policy.
