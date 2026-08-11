# Canonical Makerspace Records

## Decision Intent Log

### DI-nohos

- ID: DI-nohos
- Date: 2026-08-11 10:24:42
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Re-create Ex7's unreleased durable evidence boundary as canonical,
  pCID-selected PromiseGrid records; do not retain JSONL events as a supported
  durable compatibility format.
- Intent: Give makerspace observations, safety-hold/inspection evidence, and
  voluntary loan/return commitments exact frozen protocol meanings, durable
  author evidence, and a future-safe carriage boundary without fabricating the
  provenance of pre-recreation development data.
- Constraints: Use only the top-level semantic action `promise`; let each
  frozen family pCID define payload meaning. Browser JSON remains local
  ingress/egress only. Do not auto-convert, replay, or locally re-author old
  `events.jsonl` evidence. Preserve append-before-projection, fsync, and
  fail-closed durability behavior. Do not claim global authority, consensus,
  universal identity, automatic enforcement, or a relay embodiment before it
  is implemented and evidenced.
- Affects: `docs/thought-experiments/TE-malap-canonical-makerspace-records.md`,
  `service/`, `docs/`, `README.md`, `TODO/`, frozen family specifications,
  tests, and the repository handle ledger.

### DI-bosur

- ID: DI-bosur
- Date: 2026-08-11 10:31:22
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Freeze one shared `makerspace-record-v1` envelope specification
  and four focused durable evidence families—equipment observation, safety
  disposition, off-site loan, and off-site return—under
  `docs/protocols/makerspace-families/`, with their fixed CIDv1 mappings in
  `docs/protocols/makerspace-family-pcid-registry.md`.
- Intent: Give every current Ex7 durable claim a precise, independently
  extensible pCID-defined meaning while keeping the browser, static demo
  catalogue, member list, area policy, qualifications, and steward recognition
  honestly local until their own interoperable protocols are designed.
- Constraints: The four specs use only the top-level semantic action
  `promise`. Each admitted durable record must carry semantic author evidence
  verified against the runtime's local author-key policy. A missing,
  unrecognized, revoked, or invalid author key does not support local state
  projection; the exact bytes remain preservable as unknown or untrusted
  evidence. Relay-carriage signatures are separate future work. Do not design
  portable governance, membership revocation, steward delegation, or
  cross-makerspace key continuity in this slice.
- Affects: `docs/thought-experiments/TE-fivod-initial-makerspace-family-set.md`,
  `docs/protocols/makerspace-record-v1.md`,
  `docs/protocols/makerspace-families/`,
  `docs/protocols/makerspace-family-pcid-registry.md`, `service/`, `docs/`,
  `README.md`, tests, and the repository handle ledger.

### DI-simus

- ID: DI-simus
- Date: 2026-08-11 10:34:48
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Use one locally generated Ed25519 key pair per bootstrap demo
  author in `<runtime-root>/authors.json` (mode `0600`), and use the signed
  seven-slot `makerspace-record-v1` envelope defined by TE-zadam.
- Intent: Preserve an honest distinction between Alice's, Carol's, and Dave's
  durable claims while keeping the browser untrusted for key access and making
  local author evidence independently verifiable from the stored record.
- Constraints: The slots after the pCID selector are `record_id`,
  `signer_label`, `created_at_rfc3339`, `canonical_json_payload_bytes`,
  `author_key_id`, `author_public_key_bytes`, and `author_signature_bytes`.
  The author key ID is `ed25519:` plus lower-case hexadecimal SHA-256 of the
  public key. Sign the canonical encoding with the signature slot as CBOR null.
  Only active, locally recognized key/author pairs drive Ex7's projection;
  unknown pCIDs and unrecognized or revoked keys retain exact well-framed bytes
  without semantic projection. Malformed frames fail startup closed. This is a
  local policy, not portable identity, authority, key continuity, revocation,
  or relay carriage.
- Affects: `docs/thought-experiments/TE-zadam-makerspace-envelope-and-local-author-keys.md`,
  `docs/protocols/makerspace-record-v1.md`, `<runtime-root>/authors.json`,
  `service/`, tests, browser documentation, and the repository handle ledger.

### DI-pihav

- ID: DI-pihav
- Date: 2026-08-11 10:35:42
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Freeze the four typed payload schemas selected by TE-nazik and use
  content-addressed photo references (`blob_cid`, `media_type`, `name`) rather
  than embedded data URLs in durable records.
- Intent: Keep each pCID's promise predicate narrow and readable, preserve the
  exact accepted loan terms and return linkage, and let large photo bytes have
  independent content identity without making browser encoding the evidence
  format.
- Constraints: Observation requires `tool_id` and `observation`; safety
  disposition requires `tool_id`, `disposition` (`hold` or `clear`), and
  `assessment`; loan requires `tool_id`, `borrower_id`, `due_at`,
  `policy_version`, and `policy`; return requires `tool_id`, `loan_record_id`,
  and `condition`. `borrower_id` equals the signer label at local admission.
  Browser-created holds emit a linked observation then disposition record.
  All CID fields are CIDv1 base32 text. Blob-store path and atomic write order
  remain a separate decision; no data URL is durable evidence.
- Affects: `docs/thought-experiments/TE-nazik-makerspace-payload-schema.md`,
  `docs/protocols/makerspace-families/`, `service/`, `web/`, tests, docs, and
  the repository handle ledger.

### DI-zumat

- ID: DI-zumat
- Date: 2026-08-11 10:41:36
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Persist canonical makerspace evidence under `<runtime-root>` using
  `records.frames` with fixed `MSR1\n` header and fsynced length-prefixed
  canonical-CBOR transaction frames; store photos as fsynced CID-addressed
  blobs and apply projections only after the complete frame is durable.
- Intent: Retain exact signed record bytes, make linked observation/hold
  actions atomic to the projection, ensure every referenced photo exists before
  evidence names it, and expose interrupted writes as fail-closed evidence
  corruption rather than silently selecting a history prefix.
- Constraints: Approved paths are `<runtime-root>/authors.json`,
  `<runtime-root>/records.frames`, `<runtime-root>/blobs/<cidv1-base32>`, and
  `<runtime-root>/tmp/<opaque-temporary-name>`. Write decoded blobs to temp,
  fsync/close, rename, and fsync blob directory before appending/fsyncing the
  complete frame. Frame payload is canonical CBOR array of exact Grid record
  byte strings, with unsigned 64-bit big-endian length. Unknown or untrusted
  well-framed records remain exact bytes but do not project. Missing referenced
  blobs or malformed/partial frames fail startup closed. Orphan blobs are not
  automatically deleted; cleanup is deferred. No legacy JSONL import.
- Affects: `docs/thought-experiments/TE-tafug-makerspace-record-and-blob-store.md`,
  `<runtime-root>/records.frames`, `<runtime-root>/blobs/`,
  `<runtime-root>/tmp/`, `service/`, tests, docs, and the repository handle
  ledger.

### DI-sazir

- ID: DI-sazir
- Date: 2026-08-11 10:42:38
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Limit each decoded photo blob to 2 MiB and each canonical-CBOR
  transaction-frame payload to 1 MiB.
- Intent: Bound local parser and storage allocation while preserving the
  browser demonstration's existing photo scale; photo bytes are blobs, so
  record frames carry only small CID references and signed evidence.
- Constraints: Reject an oversized upload or frame before durable projection.
  Preserve existing fail-closed behavior for malformed/partial evidence. A
  larger blob or frame requires a superseding decision and explicit tests.
- Affects: `<runtime-root>/blobs/`, `<runtime-root>/records.frames`,
  `service/`, tests, docs, and the repository handle ledger.

### DI-lazim

- ID: DI-lazim
- Date: 2026-08-11 11:19:29
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Ex7 is a decentralized product of independently signing
  participant agents. Each participant holds their own private key and their
  agent retains and locally projects exact signed records; relay/feed services
  provide non-authoritative byte carriage only.
- Intent: Ensure durable makerspace promises are authored by their real
  participant agents, not impersonated by a shared HTTP process, while keeping
  trust and projection explicitly local to each participating agent.
- Constraints: No central server, keyholder, durable-history authority, or
  global state authority. A relay cannot sign for, authorize, revoke, or
  reinterpret a participant record. DI-simus's single-runtime author-keyring
  approach is superseded; do not implement it. Browser-to-local-agent signing,
  key continuity, and revocation require separate decisions before runtime
  code.
- Affects: `docs/thought-experiments/TE-mumut-decentralized-participant-agents.md`,
  `docs/thought-experiments/TE-zadam-makerspace-envelope-and-local-author-keys.md`,
  `service/`, `web/`, relay/feed design, tests, docs, and the repository handle
  ledger.
- Supersedes: DI-simus

### DI-janup

- ID: DI-janup
- Date: 2026-08-11 11:41:11
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Make participant-owned agents with participant signing identities Ex7's core identity boundary. Treat existing makerspace accounts only as UI login, discovery/bootstrap, and local policy inputs.
- Intent: Keep familiar account UX without confusing a website session with evidence that a participant made an exact PromiseGrid promise.
- Constraints: A kiosk/account session cannot sign for a participant. It may request a signature from the participant's reachable agent or retain an unsigned draft. Relay/feed carriage is non-authoritative. Browser/device identity and device authorization are deferred.
- Affects: `docs/thought-experiments/TE-folok-participant-agent-and-account-bootstrap.md`, `service/`, `web/`, account/bootstrap docs, tests, and handle ledger.
- Supersedes: DI-simus

### DI-sinov

- ID: DI-sinov
- Date: 2026-08-11 11:42:28
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Implement Ex7 as independent local participant agents with signed exact-record ingress, per-agent durable storage and projection, account-only UI bootstrap, and relay/feed carriage that transports exact bytes without semantic authority.
- Intent: Remove the single shared runtime as an authority while letting each participant retain, verify, assess, and exchange PromiseGrid evidence independently.
- Constraints: Relay does not sign for participants, validate makerspace semantics, compute current state, or decide trust. Ingress preserves unknown/untrusted well-framed bytes and projects only locally recognized evidence. Existing accounts are not author credentials. No central server or required relay.
- Affects: `docs/thought-experiments/TE-zajop-participant-agent-runtime-and-relay.md`, `service/`, `cmd/`, `web/`, relay implementation, tests, docs, and handle ledger.

### DI-rifib

- ID: DI-rifib
- Date: 2026-08-11 12:03:08
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Retain DI-zumat's MSR1 framing, blob-first atomic durability, and
  fail-closed replay mechanics, but scope every runtime path to one
  participant agent and remove `authors.json` from the durable-store boundary.
- Intent: Preserve exact-byte evidence mechanics without reviving the
  superseded single-runtime author-keyring model. A participant agent owns its
  private key under its own identity boundary; its record/blob store is not a
  shared makerspace authority.
- Constraints: The approved per-agent patterns are
  `<agent-runtime-root>/records.frames`,
  `<agent-runtime-root>/blobs/<cidv1-base32>`, and
  `<agent-runtime-root>/tmp/<opaque-temporary-name>`. Key storage is defined
  by the participant-identity slice, not by the framed-store implementation.
  Preserve unknown or locally untrusted well-framed record bytes without
  projection. No legacy JSONL import.
- Affects: `docs/thought-experiments/TE-tafug-makerspace-record-and-blob-store.md`,
  `docs/ex7-decentralized-redesign-roadmap.md`, `service/`, tests, docs, and
  the repository handle ledger.
- Supersedes: DI-zumat

### DI-basun

- ID: DI-basun
- Date: 2026-08-11 13:57:38
- Author: jj@thesalleys.com (JJ)
- Status: superseded
- Decision: Establish one Ed25519 root signing key per independently owned
  participant agent. The root key signs ordinary Ex7 makerspace promises and
  participant key-continuity and key-revocation promises directly. Freeze
  `participant-key-continuity-v1`, `participant-key-revocation-v1`, and the
  signed public `participant-peer-card-v1` bootstrap artifact under
  `docs/protocols/participant-identity/`, with their immutable CIDv1 mappings
  in `docs/protocols/participant-identity-pcid-registry.md`.
- Intent: Give participants a portable, independently verifiable identity
  anchor without making an account, browser profile, relay, delegated device
  key, or shared service an authoring authority.
- Constraints: Continuity and revocation use the top-level semantic action
  `promise` and preserve local trust assessment. A peer card is public,
  signed discovery data, not a durable promise, account credential, membership
  assertion, current-key authority, or global identity registry. No delegated
  operational/device keys, recovery witnesses, or account-derived signatures
  are introduced in this slice. Exact local private-key path and bootstrap
  implementation remain a subsequent Slice 1 decision.
- Affects: `docs/thought-experiments/TE-mamop-participant-root-identity-and-bootstrap.md`,
  `docs/thought-experiments/TE-baliv-key-continuity-and-revocation.md`,
  `docs/protocols/participant-identity/`,
  `docs/protocols/participant-identity-pcid-registry.md`, `service/`, tests,
  docs, and the repository handle ledger.
- Supersedes: none

### DI-tohak

- ID: DI-tohak
- Date: 2026-08-11 21:37:00
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Re-create Ex7 by Alternative A from TE-biban: retain the four
  makerspace family specifications only as preserved protocol-design artifacts
  until their exact-byte pCIDs are verified and the running product conforms
  to them. Treat `makerspace-record-v1` as an explicitly Ex7-scoped record
  profile, not a PromiseGrid-wide envelope rule. The first implementation
  claim may be made only after live exact-byte record storage/replay,
  per-family validation and projection, semantic author-evidence checks, and
  adversarial conformance tests exist.
- Intent: Repair the gap between Ex7's frozen protocol prose and its actual
  JSONL local-demo runtime, so that every future PromiseGrid claim names an
  exact contract and evidence rather than a design aspiration.
- Constraints: The present `events.jsonl` app remains a local-demo baseline,
  not an implementation of any makerspace pCID. Browser/account interaction
  is an embodiment convenience, never author evidence. Relay carriage, key
  continuity/recovery, delegated devices, blob retrieval, and portable
  governance remain separate work until each has its own source-grounded
  contract and evidence. Do not call the current Grid tag value an official
  universal allocation. Preserve unknown well-framed pCID bytes without
  assigning known semantics once the live record path exists.
- Affects: `docs/thought-experiments/TE-biban-source-grounded-ex7-recreation.md`,
  `TODO/TODO-bubuz-canonical-makerspace-records.md`,
  `TODO/TODO-giman-decentralized-redesign-roadmap.md`,
  `docs/ex7-decentralized-redesign-roadmap.md`, `docs/protocols/`, `service/`,
  `web/`, tests, implementation claims, and the repository handle ledger.
- Supersedes: DI-basun

### DI-piruf

- ID: DI-piruf
- Date: 2026-08-11 22:05:00
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Accept externally signed exact Ex7 records at the first live
  ingress boundary. Retain valid records for unknown families or unrecognized
  keys, and project only known-family records whose full public-key fingerprint
  is recognized by an injected local policy for the stated participant/role
  context.
- Intent: Make signature validity and local makerspace assessment separate, so
  the record path can become real without attributing process-made signatures
  to people or treating a chosen display label as evidence of a steward role.
- Constraints: The runtime creates and retains no participant private keys in
  this slice. Browser and account routes may submit drafts or exact signed
  bytes but never sign for a participant. No account credential, relay, or
  carriage behavior is introduced. Recognition policy is local bootstrap data,
  not portable membership, identity recovery, or a universal role registry.
- Affects: `docs/thought-experiments/TE-gozov-signed-ingress-and-local-recognition.md`,
  `service/`, `web/`, tests, runtime documentation, and the repository handle
  ledger.
- Supersedes: none

### DI-likoh

- ID: DI-likoh
- Date: 2026-08-11 22:30:00
- Author: jj@thesalleys.com (JJ)
- Status: active
- Decision: Load the agent-local recognition policy from
  `<runtime-root>/recognition.json`: a mode-`0600`, versioned JSON file of
  label and Ed25519-public-key entries. The agent fails closed on malformed
  content; an empty policy requires an explicit command-line bootstrap mode.
- Intent: Give a real operator a durable, inspectable local policy input
  without turning a browser, account, or private-key store into an authoring or
  recognition authority.
- Constraints: The file contains public keys only. It is never modified by the
  HTTP/UI process. Duplicate labels/keys, empty labels, invalid base64, wrong
  key size, unknown schema version, insecure permissions, and malformed JSON
  fail startup. Key continuity, revocation, and portable membership remain
  separate work.
- Affects: `docs/thought-experiments/TE-bilad-local-recognition-policy-configuration.md`,
  `<runtime-root>/recognition.json`, `service/`, `cmd/`, tests, and operator
  documentation.
- Supersedes: none

## Scope

Implement the source-grounded recreation locked by TE-biban / DI-tohak. The
four family files and registry are preserved protocol-design artifacts, not
implemented claims, until their exact pCID calculations and the live
conformance boundary are proven. The implementation order is in
[`TODO-giman-decentralized-redesign-roadmap.md`](TODO-giman-decentralized-redesign-roadmap.md).

## Open slices

- [ ] Verify every frozen makerspace specification's exact-byte pCID and its
  cited record-profile document identity; publish no implementation claim yet.
- [ ] Re-create durable storage and replay over exact canonical Ex7 records,
  with per-family validation and derived local projection.
- [ ] Add semantic author-evidence admission and unknown-family byte
  preservation through a separately locked participant embodiment.
- [ ] Replace browser-facing local-demo language with implementation claims,
  testing evidence, and explicit deferrals only after the live path passes.
