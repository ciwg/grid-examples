# knowledge-responsibility

`knowledge-responsibility` is the family for first-class durable
responsibilities and role-bearing operational duties. Source: `DI-sarib`.

Frozen v1 scope in ex5:

- stable responsibility identities
- summary text
- role keys
- links to items and runs

Each durable knowledge-responsibility message in this fifth frozen family owns:

- stable responsibility identity
- title
- summary
- team
- role keys
- tags
- actor
- durable timestamp

The fifth runtime slice signs and verifies this event kind under the family:

- `responsibility_created`

The current code models this family through `Responsibility` plus the signed
knowledge-responsibility envelope log. Source: `DI-sarib`.
