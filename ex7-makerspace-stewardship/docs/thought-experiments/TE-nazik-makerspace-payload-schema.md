# Makerspace Payload Schema and Blob References

TE ID: TE-nazik

## Status

decided

## Decision under test

Which exact payload fields the four initial makerspace families should freeze,
including whether attached photos travel inside durable records or by
content-addressed reference.

This TE narrows TODO `bubuz`, slice 1, after DI-nohos, DI-bosur, and DI-simus.

## Assumptions and scope

- The common envelope already supplies record ID, signer label, creation time,
  author key material, and signature. Family payloads must not duplicate those
  fields.
- Alice may record a condition observation with photos; Carol may place or
  clear a safety hold; Alice may promise to return a drill under recorded
  terms; Alice may observe/claim the drill's return condition.
- Payload bytes are canonical JSON inside the canonical Grid envelope. Any CID
  field is CIDv1 base32 text with the `b` prefix.
- Tool, member, and area IDs are local bootstrap references in this slice.
  They identify which local projection entry a record concerns; they are not
  global identity claims.

## Alternatives

### A. Typed payloads with content-addressed photo references

Freeze these fields:

| Family | Required payload fields | Optional payload fields |
| --- | --- | --- |
| equipment observation | `tool_id`, `observation` | `photos` |
| safety disposition | `tool_id`, `disposition`, `assessment` | `basis_record_id`, `photos` |
| off-site loan | `tool_id`, `borrower_id`, `due_at`, `policy_version`, `policy` | none |
| off-site return | `tool_id`, `loan_record_id`, `condition` | `photos` |

`disposition` is exactly `hold` or `clear`. A safety hold made from the browser
first writes an equipment-observation record and then a safety-disposition
record whose `basis_record_id` names that observation. A clear record supplies
inspection text in `assessment`; it may name its basis record when available.

A photo entry is an object with exactly `blob_cid`, `media_type`, and `name`.
The durable payload contains only the CIDv1 reference, accepted image media
type, and display filename. The photo bytes belong in a local
content-addressed blob store defined in the later record-store slice.

`borrower_id` must equal the envelope signer label at local admission.
`loan_record_id` names the accepted loan record that the borrower reports as
returned. The local projection checks that signer/borrower relation and tool
match; those checks are local assessment, not global enforcement.

### B. Embed data URLs in canonical payloads

Retain the current `data:image/...` payload fields inside observation,
disposition, and return records.

### C. One generic payload with optional fields

Use one object containing `tool_id`, `text`, `kind`, optional loan data, and
optional photos for all four pCIDs.

## Scenario analysis

### Normal operation

With A, Alice's condition observation names the table saw and references its
photo object by content identity. If she places a hold, the disposition points
to that observation, so later readers can distinguish the observation from the
separate local safety assessment. Carol's clear record carries her inspection
assessment. Alice's loan captures the precise policy presented at acceptance;
her return points back to that loan.

B makes the record self-contained but duplicates large image data in every
durable copy. C makes each pCID's payload rules vague, undermining the family
boundaries just selected.

### Failure, corruption, and incomplete writes

A lets the blob store verify bytes by CID before the referencing record is
admitted. A missing blob becomes an explicit incomplete-evidence condition;
the record and its claimed photo reference remain intact. The record store can
retain small, bounded canonical evidence frames and fail closed on malformed
ones.

B makes ordinary photos consume the entire record-size budget and couples
record replay to browser data-URL encoding. C permits incomplete combinations
of fields that are difficult to distinguish from corruption.

### Concurrent actors and mixed-version nodes

A gives Alice's observation and Carol's disposition independent record IDs and
author evidence. A node that cannot fetch a blob can still retain and assess
the exact record while reporting that photo bytes are unavailable. A future
family can add a new photo/media rule without changing the four v1 specs.

B requires every carrier to handle the image bytes. C forces newer payload
meanings into the generic schema and makes unknown fields unclear.

### Long-horizon evolution

A keeps a loan's accepted policy snapshot frozen in its own family and links a
return to a specific prior loan. A calibration or maintenance workflow can
compose these records without a new pCID. New interoperable meanings gain a
new family; they do not append more optional branches to v1.

B makes historical records larger and harder to move. C increases pressure to
change an already frozen generic payload whenever the makerspace adds behavior.

### Trust boundary and authority

A treats every payload as the signer's promise/claim: no record grants global
permission. The local projection separately decides whether Carol's active
recognized local key makes a `clear` disposition effective, and whether Alice
is the loan's borrower. The blob CID proves byte identity, not the truth of the
photo or a right to use the tool.

B does not improve authority assessment. C obscures the exact predicate a
signer is promising to meet for each pCID.

### Scale and operational complexity

A adds a blob store and link validation, but keeps immutable evidence records
small and exchangeable. B is initially less code but compounds storage and
bandwidth. C has fewer type definitions but pushes complexity into validators,
UI conditionals, and future protocol versions.

## Conclusions

B is rejected because browser data URLs are not a durable, scalable evidence
format. C is rejected because optional-field unions erase pCID-specific payload
meaning. A survives as the most PromiseGrid-complete design for the current
makerspace scope.

## Decision status

locked: Alternative A, by DI-pihav in
`../../TODO/TODO-bubuz-canonical-makerspace-records.md`.

## Implications and future work

- The later store slice must define a content-addressed blob-store path and
  atomic ordering between blob writes and referencing record admission.
- A missing referenced blob remains explicit incomplete evidence; no payload
  may substitute a data URL or inferred photo.
- Browser JSON may carry uploads at ingress, but the durable record must carry
  only the resulting blob reference.
