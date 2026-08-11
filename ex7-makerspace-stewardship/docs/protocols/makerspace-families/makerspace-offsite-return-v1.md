# Makerspace Off-Site Return v1

Status: frozen

Family: `makerspace.offsite.return.v1`

This family is a borrower's observation and claim that a specifically loaned
tool was returned. It uses `../makerspace-record-v1.md`; its pCID is the CIDv1
of this exact file's bytes.

## Payload

The canonical JSON object requires non-empty `tool_id`, `loan_record_id`, and
`condition`. It may contain `photos`, an array of photo references shaped
exactly as in Equipment Observation v1. No other fields are allowed.

## Meaning and local assessment

The signer promises that they observed or claim the identified tool was
returned in the stated condition, linked to the named prior loan record. This
is evidence, not a global command erasing loan history. Ex7 may clear a
displayed active loan only when the return links that loan and its signer matches
the recorded borrower; otherwise it retains valid evidence without applying
that local projection effect.

## Versioning

This document is immutable. A different return or condition model requires a
new versioned family specification and pCID.
