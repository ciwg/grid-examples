# operational-place

`operational-place` is the durable family for first-class operational place
records in `ex5`. Source: `DI-pivul`.

Frozen v1 scope in ex5:

- stable place identity
- kind
- name
- summary
- parent place reference
- tags

Each durable `operational-place` message in this family owns:

- stable place identity
- actor
- sequence and durable timestamp
- place kind
- human-readable place name
- summary text
- parent place reference
- tags

The runtime signs and verifies this event kind under the family:

- `place_created`

Child-place lists, resource membership lists, related runs, and links stay as
derived projection state rebuilt from exchanged place/resource/run/link
artifacts. The current code models the family through `Place` projections and
the signed `operational-place` envelope log. Source: `DI-pivul`.
