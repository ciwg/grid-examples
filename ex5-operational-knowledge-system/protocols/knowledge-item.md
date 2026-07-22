# knowledge-item

`knowledge-item` is the durable family for versioned operational knowledge
content. Source: `DI-mibor`.

Frozen v1 scope in ex5:

- procedures
- training content
- maintenance content
- receiving checks
- inventory audits

Each durable knowledge-item message in this first frozen family owns:

- stable item identity
- kind
- lifecycle event type
- actor
- sequence and durable timestamp
- current summary metadata for the affected revision or lifecycle transition
- responsibility links when the item is first created
- revision number when a revisioned change occurs
- revision body content when the event owns durable text
- supersedence or lifecycle notes when the event owns them

The first runtime slice signs and verifies these event kinds under this
family:

- `knowledge_item_created`
- `revision_added`
- `knowledge_item_status_changed`
- `knowledge_item_superseded`

Live drafting, participant presence, local form state, and other embodiment
state are intentionally out of scope for this frozen family. The current code
models the family through `KnowledgeItem`, `KnowledgeRevision`, and the signed
knowledge-item envelope log. Source: `DI-mibor`.
