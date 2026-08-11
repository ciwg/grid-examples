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

## Scope

Implement the decision locked by TE-malap. First define Ex7's active
makerspace family specifications and fixed pCID registry, then convert the
record store and state projection, author-evidence admission, tests, and
documentation as coherent slices.

## Open slices

- [x] Define and freeze the initial makerspace family-spec set and central pCID
  registry. Family boundaries and paths are locked by DI-bosur; shared envelope
  and local key-storage details are locked by DI-simus; payload fields and blob
  references are locked by DI-pihav. The four fixed mappings are published in
  `../docs/protocols/makerspace-family-pcid-registry.md`.
- [ ] Re-create durable storage and replay over exact canonical Grid bytes.
- [ ] Add semantic author-evidence admission and unknown-family byte
  preservation.
- [ ] Replace browser-facing legacy claims with aligned documentation and
  verification evidence.
