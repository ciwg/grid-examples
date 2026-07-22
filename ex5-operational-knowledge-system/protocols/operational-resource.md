# operational-resource

`operational-resource` is the durable family for first-class operational
resource records in `ex5`. Source: `DI-pivul`.

Frozen v1 scope in ex5:

- stable resource identity
- kind
- name
- summary
- current place reference
- tags

Each durable `operational-resource` message in this family owns:

- stable resource identity
- actor
- sequence and durable timestamp
- resource kind
- human-readable resource name
- summary text
- current place reference
- tags

The runtime signs and verifies this event kind under the family:

- `resource_created`

Related runs and links stay as derived projection state rebuilt from exchanged
resource/run/link artifacts. The current code models the family through
`Resource` projections and the signed `operational-resource` envelope log.
Source: `DI-pivul`.
