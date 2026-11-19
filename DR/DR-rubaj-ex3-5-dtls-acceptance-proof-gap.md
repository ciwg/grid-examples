# DR-rubaj — Ex3.5 DTLS acceptance-proof gap

- **DR-ID:** DR-rubaj
- **Date:** 2026-11-18 19:16:14
- **Asked by:** jj@thesalleys.com (JJ)
- **State:** open
- **Question:** What automated acceptance tests must pass before
  `ex3.5-grid-editor-dtls` can be represented as working or considered for
  merge?
- **Why this blocks progress:** The feature branch was incorrectly represented
  as merge-ready after lower-level DTLS and local-process checks. It lacks the
  required automated proofs for browser-to-browser direct-node collaboration,
  Neovim-to-browser direct-node collaboration, restart/loss recovery, and a
  real multi-host/private-network execution boundary.
- **Affects:** `ex3.5-grid-editor-dtls/`, its transport, service integration,
  browser/Neovim harnesses, documentation, CI evidence, and merge readiness.
- **Unblocks:** TODO ditaz.2, TODO ditaz.3, and any future merge request for
  the Ex3.5 branch.
- **Waiting on:** jj@thesalleys.com (JJ)
- **Decision:** Pending.
- **Linked DI:** DI-pirov; DI-halij
- **Related commits:** 14ee434; 3fc4ac6; 903955a
- **Last updated:** 2026-11-18 19:16:14

## Required acceptance evidence

1. Two independent browser harnesses attached to separate local nodes converge
   after their nodes replicate only through direct DTLS.
2. A Neovim sidecar attached to one node and a browser attached to another
   node converge through direct DTLS.
3. Restart and injected delivery-loss tests prove persisted signed envelopes
   are retried and deduplicated through the full node path.
4. A reproducible two-host/private-network test proves the documented UDP
   reachability deployment path.

Until all four are automated and passing, the branch must be described as
**incomplete** and must not be merged.

## Events

### 2026-11-18 — Real-editor local-node evidence added; DR remains open

- `TestHeadlessBrowserEditorsConvergeThroughDirectDTLS` now drives real
  headless Chrome CodeMirror surfaces on two separate local app roots in both
  directions, while the inter-node path is direct DTLS/UDP.
- `TestHeadlessBrowserAndNeovimConvergeThroughDirectDTLS` now drives real
  Chrome CodeMirror and the real Neovim sidecar in both directions over the
  same direct-node path.
- `TestHeadlessBrowserDirectDTLSRetriesAfterReceiverRestart` now injects a
  dropped UDP datagram, starts the receiver after the first persisted browser
  edit, restarts it, and verifies the later persisted edit converges.
- `errcheck ./...` and `go test ./...` pass after these tests were added.
- This is not closure evidence: the browser's local embodiment carriage is
  still WebSocket rather than raw browser DTLS, and no actual second
  host/private-network execution has been available to run the two-host
  harness. The state remains `open` and merge remains blocked.
