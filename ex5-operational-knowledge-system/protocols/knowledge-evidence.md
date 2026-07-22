# knowledge-evidence

`knowledge-evidence` is the durable family for structured evidence attached to
performed work. Source: `DI-kavup`; `DI-ribof`.

Frozen v1 scope in ex5:

- summary text
- structured facts
- optional copied attachment reference
- actor and timestamp
- stable evidence identity
- run target

Each durable knowledge-evidence message in this third frozen family owns:

- stable evidence identity
- run target
- actor
- sequence and durable timestamp
- evidence summary
- structured facts
- attachment name
- attachment path/reference
- attachment size

The third runtime slice signs and verifies this event kind under the family:

- `evidence_added`

Attachment bytes themselves remain outside this family and stay on the current
copied-file storage path. The current code models the family through `Evidence`
records attached to `RunRecord` and the signed knowledge-evidence envelope log.
Source: `DI-kavup`; `DI-ribof`.
