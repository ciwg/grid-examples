# Makerspace Equipment Observation v1

Status: frozen

Family: `makerspace.equipment.observation.v1`

This family is a promiser's observation of one tool's condition. It uses the
common envelope in `../makerspace-record-v1.md`; its pCID is the CIDv1 of this
exact file's bytes.

## Payload

The canonical JSON object requires `tool_id` and `observation`, both non-empty
text. It may contain `photos`, an array of photo references. No other fields
are allowed. Each photo reference has exactly non-empty `blob_cid` (CIDv1
base32 text), `media_type` (accepted image type), and `name` (display filename,
not a filesystem path). Payloads never embed browser data URLs or image bytes.

## Meaning and local assessment

The signer promises that the payload records the signer's observation of the
identified local tool. It does not itself place or clear a safety hold; a linked
safety-disposition record expresses that separate claim. Photo bytes are
evidence references, not proof the observation is true. Ex7 may show valid
observations and may report referenced photo bytes as unavailable. Tool-ID
interpretation is local bootstrap policy.

## Versioning

This document is immutable. A new field or meaning requires a new versioned
family specification and pCID.
