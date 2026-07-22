# operational-run

`operational-run` is the durable family for performed operational execution
records in `ex5`. Source: `DI-vamok`.

Frozen v1 scope in ex5:

- stable run identity
- knowledge-item target
- revision number
- actor and timestamp
- outcome and notes
- current `place_id`
- current `resource_ids`
- current `responsibility_ids`
- machine
- location

Each durable `operational-run` message in this sixth frozen family owns:

- stable run identity
- knowledge-item target
- revision number
- actor
- sequence and durable timestamp
- operational outcome
- notes
- place reference
- resource references
- responsibility references
- machine
- location

The sixth runtime slice signs and verifies this event kind under the family:

- `run_recorded`

Evidence, approvals, and links remain separate durable families that anchor to
the run. The current code models the family through `RunRecord` projections and
the signed `operational-run` envelope log. Source: `DI-vamok`.
