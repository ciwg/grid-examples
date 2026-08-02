# Repository tools

## `mint-handle`

`tools/mint-handle` prints and reserves one unused five-letter proquint handle
for a TODO, thought experiment, decision request, or decision intent. It
serializes allocation through `TODO/handle-namespace.tsv` and also checks all
existing repository identifiers, so concurrent agents cannot allocate the same
new handle.

Use the returned value as the handle portion of the record ID, for example
`DI-$(tools/mint-handle)`. The reservation ledger is append-only coordination
data and should be committed with the new record that consumes its handle.
