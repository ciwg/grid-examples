# knowledge-search-metadata

`knowledge-search-metadata` is not a standalone durable family in the shipped
`ex5` runtime. It is the name for derived searchable latest-state metadata used
to retrieve durable records. Source: `DI-fusok`.

Current ex5 implementation intent:

- titles
- summaries
- tags
- category/kind labeling
- context for search and drilldown

The current code folds these fields into the main projected records and derives
search behavior from the already-frozen operational families. `ex5` does not
ship a separate append-only signed search-metadata log. A later TE could still
revisit peer-visible search/index exchange, but the current shipped boundary is
derived projection state rather than a sixth durable family. Source: `DI-fusok`.
