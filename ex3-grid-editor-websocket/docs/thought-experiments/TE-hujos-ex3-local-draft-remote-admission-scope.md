# Ex3 local-draft and remote-admission scope

TE ID: TE-hujos

## Status

decided

## Decision status

Locked by DI-hadil: use the scope declaration plus architecture inventory,
classify capabilities as local admission mechanics, state the identity boundary,
and disclose the still-open private-browser manual verification. This TE does
not change profiles, capability issuance, remote access, or browser behavior.

## Decision under test

How should Ex3 publish its four source-derived local-draft pCIDs and its
remote mutation-admission mechanism so readers can distinguish peer-visible
profile meaning from local WebSocket/HTTP carriage and relay-issued capability
mechanics?

## Assumptions and trust model

- The four local profile documents under `protocols/` select the peer-visible
  meanings for live document, awareness, metadata, and publish. Their pCIDs
  must be derived from exact source bytes. Source: DI-figak; DI-povip.
- WebSocket carries live document and awareness traffic; it does not create a
  fifth protocol or define protocol meaning. Source: DI-povip.
- A relay's operator-configured bootstrap secret mints short-lived, signed,
  document- and protocol-scoped mutation capabilities for remote HTTP and
  WebSocket mutation. The capability verifies issuer, audience, document,
  pCID, action, and expiry. Source: DI-povip.
- The relay signing key is the current app-node identity. Capability audience
  and browser/Neovim participant data are scoped routing/admission inputs, not
  a general person-identity, delegation, or role-continuity system.
- Alice operates a relay and configures its bootstrap secret. Bob requests a
  remote session. Mallory obtains a token or runs an incompatible profile.

## Alternatives

### A. Scope declaration plus architecture inventory

Create `CHANGELOG.md` as the concise implementation-scope and explicit
non-claim declaration, and create the primary source-derived pCID inventory in
`docs/architecture.md`, with README links to both. Describe capability issuance
as implementation-local admission mechanics beside the scope, not a profile.

### B. Architecture-only publication

Put both pCIDs and remote-admission caveats only in `docs/architecture.md`.

### C. Treat capabilities as a fifth public profile

Publish a pCID for the capability token and present it as a peer-visible
PromiseGrid protocol alongside the four current profiles.

## Scenario analysis

### Alice changes a profile document

Under A, the architecture inventory gives a precise source-to-pCID table while
the scope declaration tells readers that each value is a local draft rather
than an upstream frozen contract. B can hold the data but buries the
implementation claim and non-claims together. C confuses a local admission
format with the existing peer-visible profile boundary.

### Bob joins remotely

A explains that Alice's relay uses a local bootstrap secret once to issue
short-lived capabilities scoped to Bob, one document, one pCID, and `mutate`.
The live message still has the selected profile meaning regardless of HTTP or
WebSocket carriage. B makes that crucial operational boundary harder to find.
C suggests that Bob can infer a general interoperable authorization protocol
from an Ex3-local implementation.

### Mallory reuses a capability

A can state the actual verification constraints without promising revocation,
delegation, cross-relay recognition, user identity, or a universal role model.
B risks an incomplete reader picture. C turns those absent semantics into
implied requirements of a fictitious fifth profile.

### Mixed-version relays

A lets a reader compare each profile pCID before assuming shared meaning; a
capability issued by one relay remains local admission evidence rather than a
cross-relay trust assertion. B provides less navigable scope. C creates a
false expectation that capability bytes alone define portable peer behavior.

### Long-horizon evolution

A permits a later frozen-spec implementation promise or a separately designed
identity/delegation protocol without rewriting Ex3 history. B loses the clear
claim/non-claim boundary. C makes a future migration harder by naming local
mechanics as settled wire semantics.

## Conclusions

C is rejected: the present capability is a relay-local admission mechanism,
not a declared peer-visible profile. B is rejected because the reader needs
both a concise scope/non-claim declaration and a navigable technical inventory.

A survives and is recommended: publish the four exact source-derived local
draft pCIDs in `docs/architecture.md`; publish an Ex3 `CHANGELOG.md` scope
declaration; link both from the README; and clearly limit remote admission to
the current relay-local bootstrap/capability implementation.

## Decisions still requiring DF

1. **Publication layout:** use a `CHANGELOG.md` scope declaration plus primary
   architecture inventory and README links (recommended), architecture only,
   or capabilities as a fifth profile?
2. **Capability classification:** describe capabilities as relay-local,
   implementation-level admission mechanics (recommended), a fifth public
   profile, or omit capability scope from published claims?
3. **Identity wording:** state that the relay key is current app-node identity
   and capability audiences/embodiment fields are not general person identity,
   delegation, or role continuity (recommended), call them user identities, or
   omit identity boundaries?
4. **Private-browser non-claim:** state that automated hardening exists but
   TODO tamuk's manual private/incognito verification remains open
   (recommended), omit the condition, or defer all scope publication until it
   closes?

## Implications for open work

- `fozoz.1` can create the scope declaration, architecture inventory, and
  README links after DF is locked.
- `fozoz.2` must separately evaluate rejected-ingress evidence policy.
- TODO tamuk remains open; this scope record must not complete it.

## Refinements

### 2026-08-10 — Private-browser verification resolved

TODO tamuk is now complete. Isolated normal and incognito Chrome sessions
converged document text in both directions through one isolated relay. The
browser-level check used local DevTools and native browser input, not a
human-driven usability review; see `DI-sodoj` in
`TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md` for the exact
evidence and constraint.
