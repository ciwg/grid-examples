# ex5 absolute socket path advertisement

TE ID: TE-zavuk
## Status
decided

## Decision under test

Where ex5 should canonicalize the runtime root and local Unix-socket path so
that `/api/meta` advertises one correct absolute socket path even when
`operational-knowledge` starts with a relative `-data-root` and clients run
from different working directories.

This is not a transport-choice question anymore. `121` already locked
runtime-first discovery through `/api/meta`. The remaining question is where
the runtime establishes the canonical absolute path that discovery should
publish.

## Assumptions

- The Unix socket remains the preferred terminal embodiment contract.
- `/api/meta` remains the discovery surface for the canonical local socket
  path.
- Mixed working directories matter because the server, CLI, and Neovim may all
  start from different shells or launcher contexts.
- Mallory is not central here; the risk is incorrect path advertisement rather
  than signed-artifact forgery.

## Alternatives

### Alternative A: canonicalize `data-root` once at startup

Resolve the runtime `data-root` to an absolute path before constructing the
App, socket server, and any runtime metadata. The rest of the process then
works only from that canonical root.

### Alternative B: keep `data-root` as given, but absolutize only the socket
advertisement

Let the runtime keep its current `data-root` string internally, but make
`Meta()` or `EmbodimentSocketPath()` convert the socket path to absolute form
right before it is advertised or used for discovery.

### Alternative C: keep the runtime unchanged and teach clients to resolve a
relative advertised socket path against the server’s cwd somehow

Keep the current relative path publication and shift the reconciliation burden
to the CLI and Neovim clients.

## Scenario analysis

### Scenario 1: default repo-root startup

Alice starts the runtime from the repo root with the default relative
`.operational-knowledge-system` root.

Alternative A works cleanly. The runtime canonicalizes once and every derived
path becomes stable.

Alternative B also works for the specific socket path being advertised.

Alternative C still leaves the client needing extra interpretation logic for
no real benefit.

### Scenario 2: custom relative runtime root

Bob starts `operational-knowledge -data-root ../state/ex5`.

Alternative A makes the entire runtime internally consistent. The stored root,
socket path, and metadata all agree on one absolute location.

Alternative B can make the advertised socket path correct, but now the runtime
still carries one relative root internally while exposing an absolute socket.
That split is workable, but less clean.

Alternative C is the weakest because clients would need some way to know what
the server’s cwd was when it resolved `../state/ex5`.

### Scenario 3: future code paths that derive more runtime-local paths

Carol later adds another local runtime artifact that should be advertised or
opened relative to the data root.

Alternative A scales best because the runtime root is already canonical at the
source.

Alternative B risks repeating “absolutize this one field too” in multiple
places.

Alternative C compounds the problem by spreading path interpretation across
clients.

### Scenario 4: PromiseGrid layering

Dave wants the runtime to be authoritative about its own local embodiment
contract.

Alternative A is strongest because the runtime owns the canonical root before
publishing any capability metadata.

Alternative B is still acceptable, but it feels more like patching one exposed
symptom than establishing one canonical runtime view.

Alternative C is least aligned because it makes clients reinterpret a runtime
contract instead of letting the runtime state it correctly.

## Conclusions

Rejected:

- Alternative C. It pushes a runtime-truth problem outward into clients.

Surviving:

- Alternative A: canonicalize `data-root` once at startup
- Alternative B: absolutize only the socket advertisement

Recommendation:

- Alternative A

Why:

- It is the cleanest PromiseGrid-aligned ownership boundary.
- It fixes the socket advertisement problem at the root instead of at one
  output field.
- It creates a better base for any future runtime-local path publication.

## Implications for TODOs and pending DIs

- TODO `124` is locked to Alternative `A`.
- The implementation should canonicalize the runtime root once at startup and
  derive the advertised socket path from that canonical root.
