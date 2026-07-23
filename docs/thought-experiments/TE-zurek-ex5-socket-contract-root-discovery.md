# ex5 socket contract root discovery

TE ID: TE-zurek
## Status
decided

## Decision under test

How terminal embodiments should discover the intended direct local Unix-socket
contract for `ex5-operational-knowledge-system` when the runtime is started
with a non-default `-data-root`.

The immediate problem is that the runtime can bind its socket under an
explicitly chosen data root, but the CLI and Neovim defaults still infer the
socket mostly from repo-root or nearest-ancestor heuristics. That is good
enough for the default dev path, but it is not a self-locating PromiseGrid
contract once operators move the runtime root.

## Assumptions

- One local runtime still owns the durable state for browser, CLI, and
  Neovim.
- Browser remains on the local HTTP adapter for the current shipped scope.
- CLI and Neovim should keep preferring the local Unix socket over HTTP when
  that contract is available.
- Mallory may observe local traffic or induce operator confusion, but cannot
  forge signed family envelopes.
- Mixed-version and mixed-launch environments matter because users may start
  the runtime with custom flags, shell wrappers, or non-repo working
  directories.

## Alternatives

### Alternative A: keep filesystem heuristics only

Continue discovering the socket by nearest known runtime root, repo root, or
cwd-relative fallback. Improve the heuristics if needed, but do not add a
runtime-advertised discovery path.

### Alternative B: add HTTP-assisted discovery for terminal clients

Keep the direct Unix-socket contract as the preferred terminal transport, but
make the socket self-locating through the already-shipped local HTTP adapter.
The runtime advertises the canonical socket path in `/api/meta`, and terminal
embodiments use that path when their local heuristics do not already point at
the right runtime root.

### Alternative C: add a separate local manifest file for socket discovery

Write a small well-known manifest file that records the active runtime root
and socket path. CLI and Neovim read that manifest first, then connect to the
socket from there.

## Scenario analysis

### Scenario 1: normal default repo-root usage

Alice runs the default runtime from the repo root and uses CLI plus Neovim in
the same working tree.

Alternative A works because the current heuristics already find the socket in
the default `.operational-knowledge-system/` root.

Alternative B also works. The direct socket remains primary, and the HTTP
discovery path stays mostly dormant.

Alternative C also works, but it introduces one more persistent local file
just to solve a case that the defaults already cover.

### Scenario 2: custom runtime root

Bob starts `operational-knowledge -data-root /srv/ex5/site-a`.

Alternative A becomes fragile. Terminal clients started from the repo root or
some other cwd will infer the wrong socket and either miss it entirely or
demote to HTTP compatibility.

Alternative B keeps the direct contract but makes it discoverable from the
already-running runtime itself. CLI and Neovim can ask the local adapter which
socket path the runtime actually owns, then switch to that direct path.

Alternative C also works, but it creates a second discovery artifact whose
lifecycle must stay in sync with the runtime root and socket file.

### Scenario 3: stale or conflicting local state

Carol has an old repo-root socket path lying around, but the real runtime now
uses a different data root.

Alternative A can mislead the embodiment into trying the stale path first,
then silently falling back.

Alternative B lets the runtime advertise the currently owned socket path.
Heuristics may still be tried first, but the authoritative recovery path comes
from the live runtime, not from stale filesystem guesses.

Alternative C depends on the manifest staying fresh. If the manifest is stale,
it becomes another source of confusion.

### Scenario 4: runtime absent or not yet started

Dave starts CLI or Neovim before the runtime exists.

Alternative A still gives a local best-effort guess and may work later once
the runtime appears.

Alternative B still needs those local guesses or explicit overrides, because
the HTTP discovery path does not exist until the runtime is up. But once the
runtime comes up, discovery becomes authoritative again.

Alternative C is similar: before runtime startup there may be no manifest yet.

### Scenario 5: PromiseGrid layering and contract honesty

Ellen wants the embodiment contract to stay direct and not collapse back into
HTTP as the real terminal transport.

Alternative A preserves that purity but leaves the direct contract under-
specified once runtime roots vary.

Alternative B preserves the direct contract while using HTTP only as discovery
metadata, not as the terminal embodiment transport itself. That matches the
current ex5 boundary well: HTTP remains an adapter surface and capability lane,
not the durable or preferred terminal contract.

Alternative C also preserves the direct contract, but the manifest becomes a
second local coordination protocol that the current runtime does not otherwise
need.

### Scenario 6: long-horizon evolution

Frank later wants remote relay and local embodiment work to stay clearly
separate.

Alternative A leaves the current ambiguity in place.

Alternative B composes cleanly with the existing `/api/meta` capability story.
The runtime already advertises what local and relay features are present, so
socket discovery becomes one more honest capability datum instead of a hidden
filesystem assumption.

Alternative C adds a parallel discovery mechanism that future docs and tools
must continue carrying.

## Conclusions

Rejected:

- Alternative A. It is not reliable enough once custom runtime roots matter.
- Alternative C. It solves the problem, but with a second local manifest layer
  that the current runtime model does not otherwise need.

Surviving:

- Alternative B: HTTP-assisted discovery for terminal clients

Recommendation:

- Alternative B

Why:

- It keeps the Unix socket as the real preferred terminal embodiment contract.
- It uses the already-shipped local adapter capability surface instead of
  inventing a second discovery protocol.
- It recovers cleanly from non-default runtime roots and stale local guesses.

## Implications for TODOs and pending DIs

- TODO `121` is locked to Alternative `B`.
- The follow-on DF under `121` is also locked: terminal clients query
  `/api/meta` first for the canonical socket path and use filesystem heuristics
  only as fallback when the runtime is unavailable.
