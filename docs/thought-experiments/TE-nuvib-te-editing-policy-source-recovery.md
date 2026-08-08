# TE Editing Policy Source Recovery

TE ID: TE-nuvib

## Status

decided

## Decision under test

Determine how this repository should apply its TE editing policy after the
policy's named source artifacts, `TE-dabol` and `TE-vudaf`, proved absent from
all local branches, reachable history, sibling CIWG workspaces, and unreachable
Git objects.

## Assumptions

- The category rules written in `AGENTS.md` remain present and are the only
  recoverable policy text in this repository.
- Recreating absent historical TEs would fabricate provenance and could
  silently change their original meaning.
- Existing TEs still need a safe, auditable editing procedure.

## Alternatives

### A. Keep the impossible prerequisite

Retain the requirement to read the two absent files before every TE edit.

### B. Recreate the missing TE files from the surrounding summary

Write new files at the referenced paths based on the policy summary in
`AGENTS.md`.

### C. Make the recoverable policy text self-contained

Remove only the references that require absent artifacts, retain the seven
category rules and holistic-reading requirement, and require review of the
affected TE plus every extant cited or citing TE that can be located.

## Scenario analysis

### Normal TE refinement

Under A, a routine refinement cannot begin because its mandatory reading list
does not exist. Under B, the agent can proceed but readers cannot distinguish
recovered text from original policy. Under C, the existing category rules guide
the refinement and the available corpus supplies the factual context.

### Historical-policy dispute

Under A, the dispute remains blocked without evidence. Under B, a newly written
file could be mistaken for historical evidence. Under C, the repository states
the loss plainly and preserves the currently recoverable policy as the operative
rule, leaving any later authoritative recovery additive and auditable.

### Long-horizon recovery

Under A, later recovery cannot repair work blocked in the meantime. Under B,
later recovery risks a conflict with fabricated replacements. Under C, an
authoritative recovered source can be added later and supersede this recovery
decision without rewriting the historical record.

## Conclusion

Alternative A is rejected because it makes valid TE maintenance impossible.
Alternative B is rejected because it fabricates missing historical provenance.
Alternative C is selected: `AGENTS.md` becomes self-contained for the
recoverable TE-editing procedure. The repair preserves the existing categories,
their append-only rules, and the holistic-reading default; it removes only the
unfulfillable mandatory reads of absent files. Source: DI-surob.

## Implications

- TE-ravuk can now receive a policy-compliant append-only refinement that
  records its completed decision state.
- If authoritative copies of `TE-dabol` or `TE-vudaf` are later recovered, a
  new TE must assess whether they supersede this recovery decision.

## Decision status

locked by DI-surob.

## Refinements

None.
