# Makerspace Safety Disposition v1

Status: frozen

Family: `makerspace.safety.disposition.v1`

This family is a promiser's safety disposition for one tool. It uses the common
envelope in `../makerspace-record-v1.md`; its pCID is the CIDv1 of this exact
file's bytes.

## Payload

The canonical JSON object requires non-empty `tool_id`, non-empty `assessment`,
and `disposition`, which is exactly `hold` or `clear`. It may contain non-empty
`basis_record_id` and `photos`, an array of photo references shaped exactly as
in Equipment Observation v1. No other fields are allowed.

## Meaning and local assessment

The signer promises they make the stated safety disposition after the stated
assessment. A hold may link the observation that prompted it; a clear records
inspection assessment and may link its basis observation. It creates no global
authority. Ex7's local projection may recognize an active configured steward's
valid clear as clearing a displayed hold and may treat a valid hold as safety
evidence. Recognition is local bootstrap policy.

## Versioning

This document is immutable. A change to disposition vocabulary, fields, or
assessment meaning requires a new versioned family specification and pCID.
