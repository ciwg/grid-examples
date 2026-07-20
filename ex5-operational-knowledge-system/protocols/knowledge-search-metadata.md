# knowledge-search-metadata

`knowledge-search-metadata` is the family for searchable latest-state metadata
used to retrieve durable records.

Current ex5 implementation intent:

- titles
- summaries
- tags
- category/kind labeling
- context for search and drilldown

The current code currently folds these fields into the main projected records,
but the family is called out explicitly so search metadata can later be
separated without changing the public model.
