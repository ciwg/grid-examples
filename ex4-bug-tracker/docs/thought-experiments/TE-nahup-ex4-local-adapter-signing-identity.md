# Ex4 local-adapter signing identity

TE ID: TE-nahup

## Status

decided

## Decision status

Locked by DI-muzal: use embodiment-local browser and CLI keys, public-key-derived
`agent_id`, explicit local enrollment, and server-side public bindings only.
This TE refines `DI-ninul`'s local signing-key choice and does not change Ex4
behavior.

## Decision under test

How should Ex4's browser and CLI hold private signing keys and establish local
adapter identity so they submit their own signed issue promises, while the
server retains only local public-key/role admission policy and never signs on a
user's behalf?

## Assumptions and trust model

- `DI-ninul` locks distinct local signing keys and explicitly rejects
  server-signing on behalf of users.
- Current `reporter`, `triage`, and `engineer` labels are local application
  policy. They do not establish general agent identity, delegation, or role
  continuity. Source: `DI-nibuh`.
- Alice runs the local server. Bob uses the browser, Carol uses the CLI, and
  Mallory may copy a browser profile, steal a CLI key file, or submit an
  unknown public key.
- The first promise slice has no remote exchange. Local enrollment is only
  bootstrap/admission for this Ex4 service and does not assert grid-wide trust.

## Alternatives

### A. Embodiment-local keys with explicit local enrollment

Browser generates and persists an extractable-or-nonextractable signing key in
its own browser-profile storage; CLI loads or creates a key at an explicit
operator-selected `--agent-key` path. Each embodiment signs its own promise.
On first use, it enrolls the public key with the local service under a selected
local role label; the service stores only the public-key binding and applies
that local policy to later signed promises.

### B. Server keyring selected by role header

The server keeps a private key per built-in role and signs after receiving the
existing `X-Bug-User` header.

### C. Checked-in or deterministic development keys

Browser and CLI derive fixed private keys from each role label or use keys
checked into the repository.

### D. One shared private key copied between browser and CLI

Generate one private key per role and distribute the same key to every local
embodiment that presents that role.

## Scenario analysis

### S1 — Bob starts the browser for the first time

A creates a browser-local key and records a public-key-to-local-role admission
binding. Bob's later promise has a signer proof that the server can verify
without holding the private key. B merely creates a server statement. C makes
every clone impersonate Bob. D makes all same-role embodiments indistinguishable
in key provenance.

### S2 — Carol uses the CLI on the same or another machine

A requires an explicit CLI key path, making key custody visible and portable
only by deliberate operator action. Carol may have a different key while the
local service maps both keys to the same local role policy. B again attributes
the promise to the service. C leaks repeatable credentials. D silently shares
one credential across machines.

### S3 — Browser private/incognito mode or profile loss

A treats a fresh or private profile as a fresh local embodiment that must enroll
another public key; it does not fabricate continuity from a presentation label.
B hides the distinction under the server key. C/D create accidental continuity
or credential reuse. A's local service may choose whether its role policy
admits the new key, but that is not a universal authorization claim.

### S4 — Mallory steals or replays material

With A, stolen private material is an embodiment-local compromise and can be
removed from Alice's local enrollment list; an unknown key is a local admission
failure. B concentrates all compromise in the server. C intentionally exposes
keys. D expands every key compromise to all shared embodiments. None of these
local actions claim general revocation semantics.

### S5 — Server restart and projection rebuild

A persists the local public-key bindings alongside accepted artifacts and can
rebuild local admission context without recovering a client's private key. B
requires server secret continuity. C/D preserve insecure key material or
ambiguous provenance. Existing `events.jsonl` remains separate local history.

### S6 — Long-horizon evolution and cost

A supports a later explicit key-rotation or recognized-role protocol without
mislabeling current local bindings as that protocol. It adds browser storage,
CLI key-file, enrollment, and test obligations. B/C/D appear simpler now but
erase the signer/projection boundary that the issue-promise layer was chosen to
teach.

## Conclusions

B is rejected because it violates `DI-ninul`: the server would create promises
on users' behalf. C is rejected because checked-in or derivable keys are not
credible local custody. D is rejected because shared private keys erase
embodiment provenance and magnify compromise.

A survives and is recommended. It is the smallest design that allows browser
and CLI to make their own signed promises while describing enrollment as local
service policy rather than a portable identity or authorization system.

## Decisions still requiring DF

1. **Key custody:** use A's browser-profile key plus explicit CLI
   `--agent-key` path (recommended), a server keyring, deterministic dev keys,
   or a shared role key?
2. **Local identity wording:** name each enrolled public key an `agent_id`
   derived from its public key and treat the selected role as local admission
   metadata (recommended), use the role label as the agent identity, or call
   the local binding a recognized role?
3. **Enrollment:** require explicit first-use enrollment through the local
   service with a public-key proof and selected local role (recommended),
   silently enroll any submitted key, or ship pre-enrolled private keys?
4. **Persistence:** store only public enrollment bindings server-side and keep
   browser/CLI private keys in their embodiment-local stores (recommended),
   store all private keys under the server runtime root, or retain no binding?

## Implications for open TODOs and pending DIs

- `kakon.2` needs these choices before it selects envelope fields, enrollment
  endpoints, key-store paths, or browser/CLI adapter behavior.
- `kakon.3` must verify that the server validates signer proof and local
  admission without accessing client private keys.
- Any key rotation, delegation, revocation, cross-machine trust, or recognized
  role continuity requires a separate TE and DI.
