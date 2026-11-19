# Ex3 implementation scope

## 2026-08-10 — Local-draft profiles and remote-admission scope declaration

This is a scope declaration for the Ex3 Grid Editor example, not a formal
PromiseGrid implementation-promise claim against a frozen upstream
specification doc-CID. Source: `DI-hadil`.

### Local profiles implemented

| Local draft profile | Current pCID | Implemented components |
| --- | --- | --- |
| `live-document` | `bafkreidoqtzj76jlzzbytfrtni655k6zk4ymo2zuhgwnzqfiz6wtrgzwje` | Relay-signed Automerge change storage/relay, browser replica, Neovim sidecar replica, and WebSocket carriage. |
| `live-awareness` | `bafkreidowbo76hmjrcqfa6l5ahmftlejqtjuzqgzusqjrxjnkngksntsl4` | Relay-signed latest-state presence feed, browser/Neovim presentation, and WebSocket carriage. |
| `document-metadata` | `bafkreih7sut3exri37qperlyvohk74vmreppzmcoxenjzozj6zkblanday` | Relay-signed latest-state document metadata and relay-backed catalog search. |
| `publish-document` | `bafkreibhvpjr5uddw5z5qmkkowfucawuvzr7hafvsc5djetkax3s6amt3a` | Relay-signed publish manifests, CAS-backed text/replica objects, and local import/exchange resolution. |
| `restore-published-version` | `bafkreicugmkq4edygey752uexzxrte3rcvhldepptc7fzwf44qpfxijlgm` | Capability-gated, relay-signed current-time restore promises that name a published manifest and carry one exact continuing Automerge change. |

The scope covers Ex3's pCID-selected signed envelopes, relay CAS and
append-only log, relay-to-relay ingestion, browser and Neovim local replicas,
and WebSocket carriage for live document and awareness traffic. WebSocket is
not a fifth profile and does not define peer-visible message meaning. Source:
`DI-hadil`; `DI-povip`.

### Provisional remote admission

For non-loopback mutation, an operator-configured bootstrap secret can mint
short-lived relay-signed capabilities scoped to a participant audience,
document ID, profile pCID, and `mutate` action. The relay verifies those claims
and expiry before allowing remote HTTP or WebSocket mutation. Loopback mutation
remains a local no-secret fast path.

This is Ex3-local implementation admission, not a fifth public profile or a
frozen PromiseGrid application-auth API. The relay signing key is the current
app-node identity. Capability audiences and browser/Neovim participant IDs,
display names, colors, and embodiment labels do not establish general person
identity, delegation, role continuity, or cross-relay recognition. Source:
`DI-hadil`; `DI-povip`.

### Explicit non-claims

Ex3 does not claim:

- conformance to a frozen upstream PromiseGrid spec doc-CID;
- interoperability with an independently implemented peer solely because it
  uses a similar profile name, WebSocket transport, or capability format;
- a general key-rotation, delegation, revocation, recognized-role, or
  role-continuity policy;
- that relay-local capability issuance or a relay observation is shared proof
  or settles another participant's intent;
- that the bootstrap secret, local HTTP adapter, browser polling, sidecar
  transport, Docker topology, or local filesystem layout is a portable
  protocol contract; or
- a human-driven usability review of private/incognito browser operation.

The restore profile is a repo-local-draft, append-only provenance claim. It
does not establish owner/admin authority, destructive rollback, byte-identical
results under concurrent editing, or frozen-spec interoperability. Source:
`DI-hihok`; `DI-tibum`; `DI-nihiz`.
  Browser-level verification has passed with isolated normal and incognito
  Chrome sessions connected to one isolated relay and converging document text
  in both directions; it used local DevTools and native browser input. Source:
  `DI-sodoj`.

When Ex3 adopts a frozen upstream spec, a later entry will use the guide's
formal implementation-promise fields (`claim`, `spec`, `scope`,
`breaking-change`, and `notes`) and name that exact frozen spec doc-CID.
