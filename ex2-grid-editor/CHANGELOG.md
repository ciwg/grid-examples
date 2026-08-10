# Ex2 implementation scope

## 2026-08-10 — Local draft profile scope declaration

This is a scope declaration for the Ex2 Grid Editor example, not a formal
PromiseGrid implementation-promise claim against a frozen upstream
specification doc-CID. Source: DI-ralit.

### Local profiles implemented

| Local draft profile | Current pCID | Implemented components |
| --- | --- | --- |
| `live-document` | `bafkreidoqtzj76jlzzbytfrtni655k6zk4ymo2zuhgwnzqfiz6wtrgzwje` | Relay-signed Automerge change storage/relay, browser replica, and Neovim sidecar replica. |
| `live-awareness` | `bafkreidowbo76hmjrcqfa6l5ahmftlejqtjuzqgzusqjrxjnkngksntsl4` | Relay-signed latest-state presence feed plus browser and Neovim awareness presentation. |
| `document-metadata` | `bafkreih7sut3exri37qperlyvohk74vmreppzmcoxenjzozj6zkblanday` | Relay-signed latest-state document metadata and relay-backed catalog search. |
| `publish-document` | `bafkreibhvpjr5uddw5z5qmkkowfucawuvzr7hafvsc5djetkax3s6amt3a` | Relay-signed publish manifests, CAS-backed text/replica objects, and local import/exchange resolution. |

The scope covers Ex2's pCID-selected signed envelopes, relay CAS and
append-only log, relay-to-relay ingestion, and embodiment-local browser and
Neovim CRDT replicas. The relay signing key is the current app identity for
these flows. Browser and Neovim participant IDs, display names, colors, and
embodiment labels are local presentation or routing data; they do not establish
independent peer identity or signing-key continuity. Source: DI-ralit.

Each pCID is derived from the named local draft profile bytes. Changing those
bytes creates a distinct local protocol profile.

### Explicit non-claims

Ex2 does not claim:

- conformance to a frozen upstream PromiseGrid spec doc-CID;
- interoperability with an independently implemented peer solely because it
  uses a similar profile name or UI behavior;
- a general key-rotation, delegation, or recognized role-continuity policy
  beyond the current relay signing identity;
- that browser or Neovim presentation fields are cryptographic identities;
- that the local HTTP adapter, browser polling, sidecar transport, Docker
  topology, or local filesystem layout is a portable protocol contract; or
- that retained relay artifacts constitute shared proof or settle another
  participant's intent.

When Ex2 adopts a frozen upstream spec, a later entry will use the guide's
formal implementation-promise fields (`claim`, `spec`, `scope`,
`breaking-change`, and `notes`) and name that exact frozen spec doc-CID.
