# Ex3 PromiseGrid alignment plan

TE ID: TE-hovij

## Status

decided

## Decision status

Locked by DI-figak: use Alternative A, describe remote admission as
provisional, retain TODO tamuk's manual verification requirement, and create a
comprehensive exercise-local testing guide. This TE does not change WebSocket
transport, remote admission, browser startup, or protocol behavior.

## Decision under test

How should Ex3 be aligned with the PromiseGrid development guide while keeping
its remote mutation-capability demo clearly provisional, preserving
decentralized relay boundaries, and honestly carrying forward the remaining
manual private/incognito-browser verification gap?

## Assumptions and trust model

- Ex3's peer-visible semantics remain the four repo-local pCID-selected
  profiles. WebSocket is carriage, not a new protocol meaning. Source:
  DI-povip.
- For non-loopback mutation, an operator-configured bootstrap secret mints
  short-lived relay-signed, document-scoped mutation capabilities. This is
  local provisional admission, not a frozen universal PromiseGrid app-auth
  API. Source: DI-povip; DI-talih.
- Relay signing identity is currently the app-node identity; browser and
  Neovim labels are embodiments rather than independent cryptographic agents.
- Private/incognito browser storage and bootstrap have automated hardening and
  regression coverage, but real manual private-browser verification remains
  open in TODO tamuk. Source: DI-ribaf; DI-bonuv; DI-sulor; DI-vobek.
- Alice operates one relay, Bob operates another, and Mallory may obtain a
  bootstrap secret, replay a capability, mistake WebSocket carriage for
  protocol meaning, or read an automated browser test as proof of every
  private-browser environment.

## Alternatives

### A. Documentation-first alignment with targeted assurance tests

First publish Ex3's exact source-derived profile and provisional remote-admission
scope; then audit remote-capability and rejected-ingress evidence boundaries;
add focused regression coverage; create the required testing guide; and finish
with a README reader pass. Keep TODO tamuk open until its manual scenario is
actually reproduced and verified.

### B. Absorb the private-browser TODO into alignment and declare completion

Treat existing automated coverage as sufficient to close the manual
private/incognito verification gap while publishing the alignment documents.

### C. Redesign the remote admission model before documentation

Replace the current bootstrap-token/capability approach with a new generalized
identity, delegation, or authorization system before describing the existing
example.

## Scenario analysis

### Alice runs a local one-relay session

Under A, the scope declaration distinguishes loopback convenience from the
remote capability path and identifies local draft profiles without overstating
them. B risks presenting the browser hardening as a completed cross-environment
claim. C delays useful documentation while inventing an unneeded authority
system for the example.

### Bob joins from another machine

A documents that the operator's bootstrap secret is only an admission input;
the relay mints short-lived, document-scoped capabilities for remote mutation,
and the peer-visible profile meaning remains pCID-selected. B can retain those
facts but mixes them with an unrelated manual browser assertion. C may improve
future capabilities but changes the teaching example before its present
contract is auditable.

### Mallory reuses a capability or sends an unsupported envelope

A makes tests and docs prove the current bounded claim: remote mutation needs
the relay's scoped capability, and local rejection/evidence behavior belongs
to the observing relay. It does not claim globally recognized roles or make a
relay's observation universal proof. B creates pressure to call broad
browser/remote behavior complete. C introduces more security semantics that
would need independent protocol and trust decisions.

### Private/incognito browser joins an existing document

A retains the honest split: automated late-join and fallback coverage has
hardened likely failure paths, but manual private-browser proof remains open.
B falsely turns tests into a complete environmental guarantee. C does not
address whether the existing private-browser scenario was actually observed.

### Long-horizon evolution

A keeps source-derived pCID documentation, remote-admission scope, evidence
policy, tests, testing guide, and README navigation independently maintainable.
B hides an unresolved operational condition. C entangles this example with a
future identity/delegation design and makes migration harder.

## Conclusions

B is rejected because an automated regression suite cannot honestly replace
the explicitly outstanding manual private/incognito verification. C is rejected
for this alignment pass because Ex3 can document and test its current local,
provisional remote-admission model without claiming it solves generalized
PromiseGrid identity or delegation.

A survives and is recommended. The alignment plan should:

1. publish Ex3's source-derived local-draft pCID inventory and provisional
   remote-admission scope/non-claims;
2. audit and document remote capability, WebSocket-carriage, and relay-local
   rejected-ingress evidence boundaries;
3. add focused source-derived and admission/evidence regression coverage while
   retaining existing cross-embodiment and multi-relay tests;
4. create `docs/testing.md` linked from the README; and
5. finish the README reader pass, explicitly linking the still-open TODO tamuk
   condition rather than declaring private-browser interoperability complete.

## Decisions still requiring DF

1. **Alignment strategy:** use the documentation-first, targeted-assurance
   plan above (recommended), absorb/close the private-browser TODO, or redesign
   remote admission first?
2. **Remote-admission claim:** publish the bootstrap secret and short-lived
   document capability as a repo-local provisional mechanism with explicit
   non-claims (recommended), describe it as a general PromiseGrid auth API, or
   omit it from the alignment scope?
3. **Private-browser status:** retain TODO tamuk as a separately open manual
   verification requirement and state that boundary in Ex3 docs (recommended),
   close it from automated coverage, or defer all Ex3 alignment until it is
   manually verified?
4. **Testing-guide scope:** create one `docs/testing.md` covering Go, browser
   JavaScript, sidecar checks, focused admission/evidence coverage, and
   existing cross-node tests (recommended), put the detail in README, or only
   list `go test ./...`?

## Implications for open work

- A new Ex3 alignment TODO should carry the five plan stages after DF is
  locked.
- TODO tamuk remains open and is not silently completed by the alignment pass.
- Any future generalized identity, delegation, or remote-authorization change
  needs its own TE and DI before code edits.

## Refinements

### 2026-08-10 — Private-browser verification resolved

TODO tamuk is now complete. Isolated normal and incognito Chrome sessions
converged document text in both directions through one isolated relay. The
browser-level check used local DevTools and native browser input, not a
human-driven usability review; see `DI-sodoj` in
`TODO/TODO-tamuk-grid-editor-private-browser-document-sync.md` for the exact
evidence and constraint.
