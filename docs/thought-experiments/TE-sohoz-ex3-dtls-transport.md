# Ex3.5 DTLS transport architecture

**TE ID:** TE-sohoz  
**Status:** needs DF

## Decision under test

How should the experimental `ex3.5-grid-editor-dtls` preserve Ex3 Grid
Editor’s browser and Neovim capabilities while using DTLS for its live
transport?

## Scope and assumptions

- Ex3’s public behavior is the baseline: browser and Neovim live document
  sync, awareness, relay verification/persistence, and its HTTP-backed
  metadata and publish flows.
- The current `live-document` and `live-awareness` payload contracts remain
  distinct. This experiment tests carriage, not new PromiseGrid wire meaning.
- This is an implementation-local experiment. The bundled PromiseGrid
  Development Guide 0.1.0 treats app transport and kernel bindings as
  provisional while DR-tuhaz and DR-davod remain open; it does not prescribe
  DTLS for this app.
- Alice and Bob are cooperative participants. Mallory can inject, replay,
  reorder, drop, spoof, or flood UDP datagrams, and can operate a stale
  client.
- RFC 9147 specifies that DTLS protects datagram traffic but does not make
  application data reliable or ordered. DTLS handshake reliability,
  fragmentation, replay handling, return-routability, retransmission timers,
  and PMTU limits are transport concerns, not replacements for Ex3’s
  application-level recovery and replay behavior.

## Alternatives

### A. Raw DTLS/UDP for every client

The Go relay and every embodiment exchange live messages directly over DTLS
over UDP. Browser JavaScript would need a raw DTLS socket API.

### B. Browser WebRTC data channels; native DTLS/UDP for the Neovim sidecar

The browser uses `RTCDataChannel`, whose WebRTC stack uses
UDP/DTLS/SCTP. The relay terminates WebRTC for browser participants and raw
DTLS/UDP for the native sidecar. HTTP remains the signaling/bootstrap path.

### C. Loopback browser bridge plus raw DTLS/UDP relay transport

The browser retains a loopback WebSocket connection to a locally launched
bridge. The bridge and the Neovim sidecar both use raw DTLS/UDP to the relay.
The browser’s WebSocket is local embodiment plumbing; all remote live traffic
between the bridge/sidecar and relay uses DTLS.

### D. Native-only raw DTLS/UDP experiment

Copy the relay and sidecar, use raw DTLS/UDP, and omit or degrade the browser
experience.

## Scenario analysis

### Normal operation

Alice opens the browser while Bob joins through Neovim. A cannot satisfy this
scenario: standard browser JavaScript has no raw DTLS/UDP socket API. B keeps
the browser and uses DTLS below SCTP, but adds WebRTC peer-connection, ICE,
and signaling ownership to the relay. C preserves the present browser UI and
uses explicit raw DTLS on every remote live link. D fails the requested
feature parity.

### Loss, reordering, and incomplete writes

Mallory drops or reorders document datagrams while Alice edits. RFC 9147
deliberately leaves application records unreliable and unordered, so A, B,
C, and D must retain Ex3’s initial replay/recovery and make durable document
changes idempotent. Awareness can be latest-state and lossy. DTLS protects
the handshake with retransmission and fragmentation but does not acknowledge
the editor’s live payloads. B gains SCTP’s selectable data-channel reliability
semantics; C needs explicit bounded application recovery and datagram-size
limits for Automerge changes. D has the same burden without browser coverage.

### Concurrent and mixed-version participants

Alice’s new browser bridge, Bob’s new sidecar, and Carol’s older WebSocket
client may all contact the relay. A and D require a hard cutover because the
old browser cannot join. B requires the relay to multiplex WebRTC and native
DTLS, while C requires it to distinguish local bridge/sidecar DTLS sessions
from any deliberately retained legacy WebSocket endpoint. Both B and C need
a versioned transport capability negotiation; neither may reinterpret the
existing pCID-selected payloads.

### Trust boundary and adversarial UDP

Mallory spoofs Alice’s address or sends handshake floods. Direct UDP choices
(A, C, D) require return-routability protection, bounded pre-authentication
state, certificate/key continuity, rate limiting, and diagnostics that do not
log capabilities or plaintext. B delegates DTLS establishment to WebRTC but
introduces signaling, ICE candidates, and WebRTC certificate/fingerprint
verification as additional trust-boundary material. In all alternatives,
DTLS authentication is transport authentication; Ex3’s signed envelopes and
relay admission checks remain necessary.

### Scale, PMTU, and long-horizon maintenance

At many participants or across constrained paths, oversized Automerge changes
can exceed a safe UDP payload. A, C, and D must bound payloads below the
effective PMTU, fragment at the application layer only with an explicit
protocol rule, and avoid turning DTLS retransmission into application
reliability. B adds the operational weight of ICE, SCTP, and WebRTC
implementations, but gives browser reachability and standardized data-channel
semantics. C is smaller than B in relay protocol scope, but requires the
local bridge to be installed, supervised, and upgraded with the browser.

## Comparison

| Alternative | Makes easier | Makes harder / new obligations | Result |
| --- | --- | --- | --- |
| A | Single raw-DTLS model | Browser support is unavailable | Rejected |
| B | Full browser support with DTLS beneath the data channel | WebRTC/ICE/signaling and mixed transport implementation | Survives |
| C | Raw DTLS remote live links while retaining browser feature parity | Local bridge lifecycle and remote-payload recovery | Survives |
| D | Smallest native proof | Required browser parity | Rejected |

## Conclusions

Alternatives B and C survive. B is the better choice if “using DTLS” includes
the browser’s standards-based WebRTC DTLS layer and the experiment is willing
to own WebRTC complexity. C is the better choice if the experiment must expose
raw DTLS/UDP in its own Go implementation while preserving the existing
browser UI through a clearly local bridge.

The recommended conclusion is **C** for this named DTLS experiment: it makes
the remote live transport observable as raw DTLS/UDP without falsely claiming
that browsers can make raw DTLS sockets. It must explicitly retain HTTP for
metadata, publish, and bootstrap, retain signed envelopes, and add
application-level recovery for durable document changes.

## Output to DF

- **Rejected:** A (no browser API) and D (does not do everything Ex3 does).
- **Surviving:** B (WebRTC data channels) and C (loopback browser bridge plus
  raw DTLS/UDP).
- **Unresolved choice:** whether the experiment prioritizes standards-native
  browser DTLS (B) or raw DTLS/UDP on every remote link (C).
- **Naming/path/runtime decisions exposed:** a new module and executable
  naming scheme, the target directory copied from Ex3, relay UDP bind/address
  pattern, loopback bridge endpoint and lifecycle if C is selected, DTLS
  certificate/identity storage, and test-only UDP port allocation all require
  DF and path approval before implementation.

## Implications for open work

No Ex3 protocol pCID changes are implied. The chosen transport must be
documented as an implementation-local experimental carriage choice, not as a
frozen PromiseGrid requirement.

## Sources

- [RFC 9147: DTLS 1.3](https://datatracker.ietf.org/doc/html/rfc9147) —
  datagram delivery properties, handshake retransmission, return-routability,
  replay, and PMTU guidance.
- [W3C WebRTC](https://www.w3.org/TR/webrtc/) — browser data channels and
  their SCTP transport over an established DTLS connection.
