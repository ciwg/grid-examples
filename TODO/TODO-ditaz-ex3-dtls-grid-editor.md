# TODO ditaz - Ex3.5 DTLS Grid Editor

## Decision Intent Log

### DI-ditaz

- Date: 2026-11-18 17:06:17
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Build the experimental `ex3.5-grid-editor-dtls` as a copy of Ex3
  in which the browser uses a loopback WebSocket bridge and the bridge plus
  Neovim sidecar carry remote live document and awareness traffic to the relay
  with raw DTLS/UDP.
- Intent: Exercise observable raw DTLS/UDP on every remote live link while
  preserving Ex3 browser and Neovim behavior without claiming that browser
  JavaScript can open a raw DTLS socket.
- Constraints: This is implementation-local experimental carriage, not a
  PromiseGrid transport standard. HTTP retains bootstrap, metadata, and
  publish responsibilities. Existing signed envelopes and pCID-selected
  `live-document`/`live-awareness` semantics remain unchanged. Durable
  document recovery must not rely on DTLS application-data reliability.
- Affects: `ex3.5-grid-editor-dtls/`, its runtime transport configuration,
  tests, documentation, and this TODO record.

### DI-logob

- Date: 2026-11-18 17:10:43
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Use a pinned, patched `github.com/pion/dtls/v3` dependency as the
  in-process Go implementation of raw DTLS/UDP for the experimental relay,
  bridge, and Neovim sidecar.
- Intent: Keep DTLS lifecycle, certificate verification, MTU configuration,
  error handling, and deterministic test control within the copied Go
  application instead of introducing an external-process and IPC boundary.
- Constraints: The selected version must be v3.1.4 or later. This dependency
  is an implementation choice, not a PromiseGrid transport standard. DTLS
  must not be configured to skip certificate verification outside tests.
- Affects: `ex3.5-grid-editor-dtls/go.mod`, `go.sum`, live-transport code,
  tests, documentation, and this TODO record.

### DI-higom

- Date: 2026-11-18 17:11:40
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Deliver `live-document` over raw DTLS/UDP with application
  acknowledgements, bounded retransmission, deduplication, and replay/recovery;
  carry `live-awareness` as latest-state traffic that may be dropped when
  superseded.
- Intent: Preserve Ex3’s eventually shared document state while keeping
  ephemeral cursor/presence traffic lightweight and compatible with DTLS’s
  unordered, unreliable application-data behavior.
- Constraints: The relay’s signed envelope and admission checks remain in
  force. Retransmission must be bounded, idempotent, and observable. A stale
  awareness update must not overwrite a newer local view.
- Affects: `ex3.5-grid-editor-dtls/` live transport, recovery tests,
  diagnostics, documentation, and this TODO record.

### DI-puzoj

- Date: 2026-11-18 17:12:10
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Persist a self-signed DTLS relay certificate under the relay data
  root and distribute its fingerprint through the existing authenticated
  bootstrap response; the browser bridge and Neovim sidecar pin that
  fingerprint for their DTLS sessions.
- Intent: Keep the experimental transport self-contained while making relay
  transport identity verifiable, without relying on an operator-managed CA or
  disabling certificate verification.
- Constraints: The certificate private key remains local to the relay data
  root. Bootstrap must not expose it. Certificate verification may be skipped
  only in explicitly isolated tests. Signed envelope identity remains the
  application-level identity.
- Affects: `ex3.5-grid-editor-dtls/` bootstrap, DTLS configuration, runtime
  data, tests, documentation, and this TODO record.

### DI-tizal

- Date: 2026-11-18 17:12:49
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Name the local browser executable `grid-dtls-bridge`; place new
  Go transport code in package `transport`; prefix transport-specific exported
  identifiers with `DTLS`; and use descriptive role-based local identifiers
  such as `relayCertificate`, `deliveryReceipt`, and `retransmitDeadline`.
- Intent: Make the experimental DTLS scope evident at executable and API
  boundaries while keeping names legible beside copied Ex3 code.
- Constraints: Existing copied identifiers remain unchanged unless a transport
  behavior change requires a new identifier. New identifiers must follow this
  convention.
- Affects: `ex3.5-grid-editor-dtls/` Go sources, scripts, tests,
  documentation, and this TODO record.

### DI-halij

- Date: 2026-11-18 17:16:48
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Make `ex3.5-grid-editor-dtls` a genuinely runnable PromiseGrid
  development example: browser and Neovim participants on separate hosts must
  collaborate through the bridge, raw DTLS/UDP relay transport, persistence,
  delivery recovery, and verified proof tests.
- Intent: Ensure that experimental DTLS carriage means evaluated working
  behavior rather than a mock, documentation-only prototype, or incomplete
  feature slice.
- Constraints: The example must retain Ex3 feature parity within the selected
  transport architecture and document its deployment/verification procedure.
  This does not itself claim a production service or a frozen PromiseGrid
  transport standard.
- Affects: `ex3.5-grid-editor-dtls/`, tests, runnable scripts,
  documentation, and this TODO record.

### DI-badig

- Date: 2026-11-18 17:17:28
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Bind the Ex3.5 raw DTLS/UDP relay listener to `:7026` by default
  so the documented normal run supports collaboration from separate hosts.
- Intent: Make the PromiseGrid development example immediately runnable for
  the multi-host DTLS collaboration scenario it demonstrates.
- Constraints: Public UDP exposure does not remove the requirement for
  certificate pinning, authenticated bootstrap, bounded connection attempts,
  rate limiting, signed-envelope admission, and recovery tests. Operators may
  override the bind address with `--dtls-listen`.
- Affects: `ex3.5-grid-editor-dtls/` relay configuration, bridge/sidecar
  configuration, tests, documentation, and this TODO record.

### DI-pirov

- Date: 2026-11-18 17:21:33
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Replace the relay-oriented Ex3.5 topology with direct peer nodes.
  Each node owns its identity, durable state, and DTLS sessions; participants
  exchange signed peer-invite links containing direct endpoints and DTLS
  certificate fingerprints through out-of-band channels such as email, text,
  QR code, or copied links.
- Intent: Deliver a genuinely decentralized, working PromiseGrid development
  example without hiding peer discovery or collaboration behind a central
  service.
- Constraints: No central relay, discovery service, or mandatory hosted
  rendezvous component is permitted in the normal collaboration path. Nodes
  must persist outbound document changes and retry direct delivery to known
  peers. The documented initial deployment requires reachable peer UDP
  endpoints or an operator-provided private network.
- Affects: `ex3.5-grid-editor-dtls/` topology, persistence, transport,
  invitation flow, scripts, tests, documentation, and this TODO record.
- Supersedes: DI-ditaz, DI-puzoj, DI-badig

## Work items

- [ ] ditaz.1 Copy Ex3 into the approved experimental directory.
- [ ] ditaz.2 Replace remote WebSocket live carriage with raw DTLS/UDP while
  preserving loopback browser bridge behavior.
- [ ] ditaz.3 Add deterministic DTLS transport and full-feature parity tests.
- [ ] ditaz.4 Document configuration, trust boundaries, recovery behavior, and
  verification commands.
