# OKAR Package Template

Copy this directory to start a new installed executable package.

After copying:

1. rename the package ID and command path
2. decide whether the package reuses an existing family from the [built-in registry](../../docs/protocols/package-family-pcid-registry.md) or introduces a new shared durable-record meaning
3. for a new family, add an immutable versioned specification in your package's `protocols/` directory and calculate its CIDv1 from the exact file bytes
4. put that exact CIDv1 in the manifest, `describe` output, claim block, and every emitted canonical record
5. implement real `validate` and `run` behavior

Do not copy a symbolic `pcid:...` value or derive a pCID from the family name at
runtime. Workflows normally compose existing family pCIDs and do not get a pCID
of their own. See the [package author guide](../../docs/package-author-guide.md)
for the complete contract. Source: `DI-jusij`; `DI-solan`.

Then install with:

```bash
moks package install ./your-package
```
