# knowledge-item

`knowledge-item` is the durable family for versioned operational knowledge
content.

Current ex5 implementation intent:

- procedures
- training content
- maintenance content

Each item is expected to carry:

- stable item identity
- kind
- current summary metadata
- responsibility links
- revision history
- revision body content

The current code models this family through `KnowledgeItem` plus
`KnowledgeRevision`.
