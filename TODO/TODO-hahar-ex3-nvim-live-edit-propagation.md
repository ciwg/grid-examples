# TODO hahar - Fix Ex3 live Neovim edit propagation

## Goal

Make two independently launched Ex3 Neovim embodiments of the same document
propagate a user edit from either window to the other through the relay.

## Incident

Date: 2026-10-01

Observed: Two live Neovim windows opened as `grid-editor://demo` against the
same `http://127.0.0.1:7038` relay. A line entered in the first window did not
appear in the second window.

Expected: A local text change in either `demo` embodiment is sent through its
sidecar and relayed to the other embodiment without a manual refresh.

Scope: Ex3 Neovim plugin, Neovim sidecar, and relay interaction only. Browser
behavior is out of scope.

## Tasks

- [ ] hahar.1 Reproduce the live edit-propagation failure with a deterministic two-Neovim test.
- [ ] hahar.2 Identify and fix the failing plugin, sidecar, or relay handoff.
- [ ] hahar.3 Verify an edit from each live embodiment appears in the other and update the demo coverage.
