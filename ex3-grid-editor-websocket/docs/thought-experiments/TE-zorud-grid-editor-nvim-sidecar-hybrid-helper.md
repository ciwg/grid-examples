# grid-editor Neovim sidecar hybrid helper

TE ID: TE-zorud
## Status
decided

## Decision under test

How should `grid-editor` deliver a real Neovim CRDT sidecar now that the relay/browser path is live, but the workspace still lacks a native Go Automerge replica?

## Assumptions

- The relay and browser CRDT path are already implemented.
- The workspace contains a working Node-based Automerge helper in `viduct`, but not a native Go Automerge engine.
- The public command path `cmd/grid-nvim-sidecar` is already locked.

## Alternatives

1. Keep Neovim on the old snapshot polling client.
2. Stop until a native Go Automerge replica exists.
3. Use a Go launcher that runs a bundled Node Automerge helper behind the locked `grid-nvim-sidecar` command.
4. Replace the Go command path with a direct Node helper surface.

## Scenario analysis

### Real browser to Neovim convergence now

Alternative 1 leaves the main gap untouched.

Alternative 2 preserves purity but ships no working Neovim CRDT path.

Alternative 3 gives Neovim a real Automerge replica immediately while preserving the repo’s command vocabulary and sidecar boundary.

Alternative 4 also works technically, but it discards the locked command path and spreads implementation details into the user-facing surface.

### Command and packaging stability

Alternative 1 keeps stability by doing nothing.

Alternative 2 also keeps stability, but only by deferring the problem.

Alternative 3 keeps `grid-nvim-sidecar` as the stable entrypoint and makes the Node helper internal.

Alternative 4 makes Node itself part of the public command path.

### Long-horizon replacement with a native Go replica

Alternative 1 offers no migration path.

Alternative 2 waits for the ideal future but blocks present progress.

Alternative 3 isolates the Automerge engine behind a sidecar protocol so a future native Go replica can replace the helper without changing the plugin contract.

Alternative 4 makes later replacement harder because the helper becomes the surface.

## Conclusions

- Reject alternative 1 because it keeps Neovim on the broken collaboration model.
- Reject alternative 2 because it delays a needed working embodiment.
- Reject alternative 4 because it breaks the already-locked command shape.
- Keep alternative 3.

## Implications for TODOs and DIs

- Add a Neovim-side TODO for the hybrid sidecar slice.
- Keep the Go command path while treating the Node helper as an internal engine.
- Move the Lua plugin to a stdio sidecar protocol and remove direct relay snapshot polling from the hot path.

## Decision status

locked via DI-sulod and DI-gafit
