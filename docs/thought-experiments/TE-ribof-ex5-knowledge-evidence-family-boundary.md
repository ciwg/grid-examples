# Ex5 Knowledge-Evidence Family Boundary

TE ID: `TE-ribof`
## Status
decided

## Decision under test

What the third frozen ex5 PromiseGrid family should own when `knowledge-evidence`
is frozen: evidence metadata only, evidence metadata plus attachment references,
or evidence metadata plus copied attachment bytes inside the signed family.

Related TODO:

- `096` - `ex5-operational-knowledge-system/TODO/TODO-kavup-ex5-third-frozen-protocol-family.md`

## Assumptions

- `ex5` already has two frozen PromiseGrid families: `knowledge-item` and
  `knowledge-approval`.
- The current runtime models evidence as `evidence_added` events attached to
  runs, with optional copied attachments stored under `attachments/<run-id>/`.
- The current browser, CLI, and Neovim surfaces should remain on the same
  local HTTP adapter during this migration slice.
- The PromiseGrid dev guide favors staged, spec-first, claim-driven migration
  by narrow durable family.
- Attachment bytes may eventually move to a more content-addressed or relay-
  visible store, but that is not already implemented in ex5.

## Alternatives

### Alternative A

Freeze `knowledge-evidence` around durable evidence metadata only:
- evidence identity
- run target
- actor
- timestamp
- summary
- structured facts

Attachment name/path/size remain outside the frozen family.

### Alternative B

Freeze `knowledge-evidence` around durable evidence metadata plus attachment
reference metadata:
- all metadata in Alternative A
- attachment name
- attachment path/reference
- attachment size

The family claims the durable reference, but not the attachment bytes
themselves.

### Alternative C

Freeze `knowledge-evidence` around metadata plus the copied attachment bytes as
part of the signed family itself.

## Scope and systems affected

- `protocols/knowledge-evidence.md`
- evidence creation paths in browser and CLI
- runtime storage and replay for evidence
- attachment storage contract clarity
- PromiseGrid claims and boundary docs

## Scenario analysis

### Scenario 1: normal operator evidence capture

Alice records a receiving run and attaches structured facts plus an uploaded
photo or document.

Alternative A:

- freezes the smallest durable concept
- leaves attachment references outside the frozen contract even though users
  experience them as part of one evidence artifact
- weakens the trust story when evidence depends on attachments

Alternative B:

- matches the current product concept well
- lets the evidence artifact claim both facts and the durable attachment
  reference without dragging raw bytes into the signed family
- preserves current copied-file storage under the existing adapter/runtime

Alternative C:

- makes the signed family much heavier immediately
- increases storage and serialization obligations
- entangles the first evidence-family slice with attachment-byte transport and
  migration concerns

Result:

- B best matches the shipped behavior without overreaching.

### Scenario 2: failure, corruption, or incomplete writes

Bob restarts after partial evidence writes or a missing attachment file.

Alternative A:

- leaves the signed family unable to say whether an attachment reference was
  supposed to exist

Alternative B:

- lets replay verify the durable evidence record and its attachment reference
- still leaves attachment-byte existence as a separate storage integrity
  question, which matches current ex5 reality

Alternative C:

- demands immediate signed-byte treatment for large files and partial writes
- greatly enlarges the failure surface

Result:

- B is the cleanest match to current storage integrity boundaries.

### Scenario 3: staged migration and mixed coverage

Carol uses browser evidence upload while Dave later audits the run history from
CLI or Neovim.

Alternative A:

- creates a gap where evidence attachments are visible in embodiments but not
  part of the frozen evidence claim

Alternative B:

- keeps user-visible evidence artifacts coherent across embodiments
- allows later attachment-byte storage changes without redefining the evidence
  family itself

Alternative C:

- ties the early evidence slice to a much larger storage redesign

Result:

- B preserves the staged migration shape.

### Scenario 4: long-horizon evolution

Ellen wants later transport or CAS work without rewriting what evidence means.

Alternative A:

- under-specifies evidence too much

Alternative B:

- freezes the evidence concept at the right level: facts plus attachment
  reference
- leaves byte transport/storage evolution for later implementation slices

Alternative C:

- risks freezing current copied-file storage assumptions too deeply

Result:

- B gives the best long-horizon seam.

### Scenario 5: trust-boundary clarity

Mallory audits what exactly a signed evidence artifact proves.

Alternative A:

- proves only summary/facts, not whether an attachment was part of the record

Alternative B:

- proves the evidence summary/facts plus that a named attachment reference and
  size were part of the durable artifact
- does not overclaim that the raw attachment bytes are yet carried by the
  PromiseGrid family

Alternative C:

- proves the most, but at the cost of a much larger first slice

Result:

- B is the clearest honest trust statement for current ex5.

## Conclusions

Rejected alternatives:

- Alternative A: under-freezes the current evidence artifact by dropping
  attachment references out of the family.
- Alternative C: overreaches into attachment-byte handling too early.

Surviving alternative:

- Alternative B: one `knowledge-evidence` family covering evidence metadata
  plus attachment references, but not attachment bytes themselves.

Locked decision:

- Lock Alternative B.
- Add a stable durable evidence ID before signing the new family artifacts.

## Implications for open TODOs and pending DIs

- `096` should freeze `knowledge-evidence` around facts plus attachment
  references if Alternative B is chosen.
- A later follow-on can decide whether attachment-byte storage should become a
  PromiseGrid-native layer or remain an adapter/runtime-local storage concern.
