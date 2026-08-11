# Makerspace Record Envelope and Local Author-Key Policy

TE ID: TE-zadam

## Status

decided

## Decision under test

How Ex7's shared `makerspace-record-v1` envelope should bind a family payload
to semantic author evidence, and how the single-runtime demo should create,
store, recognize, and locally revoke author keys without claiming portable
identity or global authority.

This TE narrows TODO `bubuz`, slice 1, after the clean canonical boundary
(TE-malap / DI-nohos) and the four-family set (TE-fivod / DI-bosur) were
locked.

## Assumptions and scope

- Every durable record is canonical CBOR `grid(...)`; its selector pCID chooses
  one of the four frozen makerspace family specifications. The shared envelope
  owns common evidence slots, while each family owns its payload bytes and
  validation.
- Alice, Carol, and Dave are local demonstration identities. Their names and
  initials are presentation/bootstrap data, not a portable identity registry.
- The browser sends a local request naming the selected demo author. The server
  may access local private keys, but the browser, event store, and photos may
  not.
- Mallory can provide an unknown pCID, malformed bytes, a valid record signed
  by an unrecognized key, or a record signed by a locally revoked key. Ex7 must
  preserve parseable exact bytes without letting those bytes alter its current
  makerspace projection.
- Relay carriage, remote key discovery, key delegation, cross-makerspace
  rotation, and a public revocation protocol are outside this first slice.

## Alternatives

### A. Per-demo-author Ed25519 keyring and signed common envelope

On first initialization, generate one Ed25519 key pair for each bootstrap demo
author and persist the private material only at the approved dynamic path
`<runtime-root>/authors.json` with mode `0600`. The file contains the active or
revoked local status, current public key, and private key for each local demo
author. Its values are a single-machine simulation boundary, not an identity
service.

The canonical envelope has these protocol-defined slots after its pCID selector:

```text
[ record_id,
  signer_label,
  created_at_rfc3339,
  canonical_json_payload_bytes,
  author_key_id,
  author_public_key_bytes,
  author_signature_bytes ]
```

`author_key_id` is `ed25519:` plus the lower-case hexadecimal SHA-256
fingerprint of `author_public_key_bytes`; it is a local key reference, not a
pCID. The signature covers the exact canonical Grid encoding of the same slots
with `author_signature_bytes` represented as CBOR null. A record is locally
projectable only when the signature verifies, its key ID matches its public
key, and its author/key pair is currently active in `authors.json`.

The primary framed record store keeps every well-framed canonical Grid byte
sequence exactly, including unknown-pCID, unrecognized-key, and revoked-key
records. The current makerspace projection consumes only recognized, active,
validly signed records for known pCIDs. Malformed frames still fail startup
closed because they cannot be preserved as an exact valid record.

### B. One runtime signing key plus an actor label

Generate one local Ed25519 key at `<runtime-root>/runtime-key.json`; each record
contains an Alice/Carol/Dave label but is signed by the runtime key.

### C. External author keys only

Require the browser or a separate external signer to provide every author
signature. The server stores no private keys and rejects actions from local demo
members without externally supplied key material.

## Scenario analysis

### Normal operation

Under A, Alice's observation carries Alice's locally provisioned public key and
signature; Carol's hold-clearance record carries Carol's. The runtime can say
precisely that it verified a local demo key associated with the chosen author,
while the browser remains a convenience ingress. A later reader can verify the
record's mathematics without trusting the UI label alone.

B is operationally smaller, but the honest author is the runtime. Calling the
labelled person the semantic author would be false. C gives strongest
separation, but turns the browser-first local demo into an undeployed key
management product before it can exercise its actual makerspace flows.

### Failure, corruption, and incomplete writes

A writes and fsyncs the keyring on initial setup before the first signed record
is accepted. The existing append-before-projection rule then applies to exact
framed record bytes. A missing, unreadable, malformed, or permission-insecure
keyring causes startup/action failure rather than an unsigned fallback.

B has fewer keys but the same durable key-file failure. C avoids a private
keyring but makes normal authoring fail whenever the external signer is absent.
None of the choices permit partial record frames to drive projection state.

### Concurrent actors and mixed-version nodes

With A, Alice and Carol have distinct keys, so two signed claims remain
distinguishable even if their browser requests are serialized by the current
single process. A later implementation can preserve an unknown pCID or key
exactly and decline to interpret it. The fixed signing view ensures every
verifier sees the same bytes.

B merges all author evidence into one key, weakening later assessment. C can
interoperate with independently held keys but has no smooth local-demo path.

### Long-horizon evolution and key change

A can mark a local demo key revoked in `authors.json`; existing record bytes
remain intact, but new projection rebuilds no longer treat that key as active.
Replacing a local demo key is a local bootstrap operation with no claim that
other runtimes recognize the replacement. A later durable key-continuity or
revocation family can supersede this policy without rewriting existing record
bytes.

B can revoke only the entire runtime author. C needs a portable discovery,
rotation, and recovery design immediately. A therefore supplies a visible,
testable local revoke state without pretending it solves cross-node continuity.

### Trust boundary and authority

In A, signature validity proves only that the configured local demo key made
the record. It does not prove Carol has universal authority to clear a hold;
the projection separately applies the local bootstrap recognition rule.
Unknown, unrecognized, and revoked-key records are retained but do not alter
the display. This preserves evidence without treating byte possession or a
signature as a command.

B makes it harder to distinguish a runtime service action from a human claim.
C has good separation but no configured local trust path for Ex7's shipped
demo. None of the alternatives establish global membership, universal
revocation, or automatic enforcement.

### Scale and operational complexity

A adds a small keyring, deterministic signature test fixtures, and one shared
envelope implementation. It creates the smallest credible author-evidence
boundary for the demo and a direct path to future relay assessment. B is less
code but lower-fidelity evidence. C is significantly more operational work and
requires a browser/key-management embodiment beyond current scope.

## Conclusions

B is rejected because it makes the runtime the actual semantic author while
the product presents Alice or Carol as the author. C is rejected for this
slice because it defers ordinary local demo actions behind an undeclared
external signing system.

A survives and is recommended: use per-demo-author Ed25519 keys at
`<runtime-root>/authors.json`, the seven-slot common envelope above, a
null-signature canonical signing view, and a projection rule that accepts only
recognized active local authors. Preserve all well-framed canonical bytes;
preserve malformed evidence only as a startup-failing corruption condition.

## Output to decision framing

The surviving DF choice is:

- **A (recommended):** per-demo-author local Ed25519 keyring,
  `<runtime-root>/authors.json`, the listed seven envelope slots, and
  preservation-without-projection for unknown, unrecognized, or revoked-key
  records.

If selected, a DI will lock the envelope, signer, key-storage, and local
revocation details. Only then may Ex7 create the specs, registry, codec, and
runtime code.

## Decision status

locked: Alternative A, per-demo-author local Ed25519 keyring and the listed
seven-slot envelope, by DI-simus in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications and future work

- The `authors.json` keyring is a production-data dynamic path under the
  caller-approved runtime-root. It is created on first persistent runtime
  initialization, read at startup, written only during initialization or an
  explicit future local key-management action, and is never returned by HTTP.
- The framed record-store filename and CAS/photo storage paths are separate
  file/path decisions; this TE approves neither.
- A future portable author/role/revocation protocol must have its own frozen
  pCID and supersede the relevant local-policy claim rather than silently
  changing the meaning of existing keys.
