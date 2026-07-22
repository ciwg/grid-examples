# knowledge-link

`knowledge-link` is the durable family for typed links across operational
records. Source: `DI-votek`.

Frozen v1 scope in ex5:

- connect responsibilities to items
- connect runs to other operational records
- preserve typed relation names instead of untyped references only

Each durable knowledge-link message in this fourth frozen family owns:

- stable link identity
- from-type
- from-id
- to-type
- to-id
- relation
- notes
- actor
- durable timestamp

The fourth runtime slice signs and verifies this event kind under the family:

- `link_added`

The current code models this family through `Link` plus the signed
knowledge-link envelope log. Source: `DI-votek`.
