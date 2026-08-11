# OKAR Package Record v1

Status: frozen

This specification defines the canonical durable record carriage for every
first-party and installed OKAR package family.

Each record is canonical CBOR `grid()` bytes selected by its individual family
pCID. This document defines the carriage shared by those family specifications;
it is not itself the selector for every package record. Its slots are, in
order: family name, record ID, signer label, RFC3339 timestamp, canonical JSON
payload bytes, author key ID or null, author public key or null, and author
signature or null. The family name and canonical JSON payload are validated by
the active family contract.

An author signature signs the same canonical Grid envelope with the final three
author slots set to null. It is evidence of authorship, not authority. Unknown
families retain their exact canonical bytes without an inferred validator.
