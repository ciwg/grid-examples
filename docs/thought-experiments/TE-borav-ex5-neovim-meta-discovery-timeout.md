# ex5 Neovim meta discovery timeout

TE ID: TE-borav
## Status
decided

## Decision under test

How Neovim should bound runtime-first `/api/meta` socket discovery so that a
dead or blackholed `OKS_BASE_URL` does not stall editor startup.

`121` already locked runtime-first discovery, so this TE is no longer about
whether Neovim should ask `/api/meta` first. The question is how that request
should fail fast while preserving the same direct-socket preference.

## Assumptions

- Neovim still prefers the direct local Unix socket for request/response and
  live-draft carriage.
- `/api/meta` remains the first discovery step when no explicit socket path is
  configured.
- Repo-root fallback should remain available when discovery fails.
- Mixed environments matter because some users will have a live local runtime,
  while others will point `OKS_BASE_URL` at something unreachable.

## Alternatives

### Alternative A: bounded HTTP discovery timeout with immediate local fallback

Keep the current startup shape, but add an explicit short timeout to the
`/api/meta` discovery request. If the timeout or any network error occurs,
Neovim falls back immediately to its local repo-root socket inference.

### Alternative B: asynchronous discovery after startup

Start Neovim immediately with local fallback assumptions, then query
`/api/meta` asynchronously and switch to the runtime-advertised socket path
later if needed.

### Alternative C: launcher-preseeded canonical socket path only

Undo runtime-first discovery in the editor itself and require the launcher or
environment to pre-populate `OKS_SOCKET_PATH`.

## Scenario analysis

### Scenario 1: normal local runtime

Alice starts Neovim while the local ex5 runtime is already up.

Alternative A works well: the discovery request returns quickly and Neovim
uses the canonical socket path.

Alternative B also works, but it complicates startup semantics because the
transport target may change after startup.

Alternative C works only when the launcher or environment was configured
correctly in advance.

### Scenario 2: dead `OKS_BASE_URL`

Bob has `OKS_BASE_URL` pointing at a host that blackholes connections.

Alternative A fails fast and falls back predictably.

Alternative B also avoids blocking if implemented carefully, but it introduces
asynchronous state changes early in editor startup.

Alternative C avoids the HTTP stall only by abandoning the runtime-first
discovery model that `121` already chose.

### Scenario 3: custom runtime root discovered through HTTP

Carol uses a custom runtime root that Neovim cannot infer locally.

Alternative A still reaches the right answer whenever the runtime is actually
available, and only falls back when discovery fails.

Alternative B can also reach the right answer, but now transport selection may
change after the editor already began using the local fallback path.

Alternative C gives up the self-locating behavior entirely unless the launcher
or environment was updated perfectly.

### Scenario 4: PromiseGrid layering and embodiment clarity

Dave wants runtime-first discovery to stay the truth source, but not at the
cost of brittle startup behavior.

Alternative A keeps the layering clean and operationally simple: the runtime
is authoritative when reachable, and Neovim reuses the already-approved local
fallback only when discovery fails.

Alternative B is more dynamic, but adds unnecessary complexity to the
embodiment path.

Alternative C shifts too much responsibility back to wrappers and operator
setup.

## Conclusions

Rejected:

- Alternative C. It reopens the already-locked discovery model.

Surviving:

- Alternative A: bounded HTTP discovery timeout with immediate local fallback
- Alternative B: asynchronous discovery after startup

Recommendation:

- Alternative A

Why:

- It preserves the runtime-first truth source.
- It is the narrowest fix for the startup stall risk.
- It avoids adding asynchronous transport churn to editor startup.

## Implications for TODOs and pending DIs

- TODO `125` is locked to Alternative `A`.
- The implementation should add a short bounded discovery timeout with
  immediate repo-root socket fallback and cover that behavior in regression
  tests.
