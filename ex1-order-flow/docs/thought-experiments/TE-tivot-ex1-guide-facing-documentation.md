# Ex1 guide-facing documentation completion

TE ID: TE-tivot

## Status

decided

## Decision status

Locked by DI-motiv: extend the README with the concise walkthrough, cover the
happy path plus refusal and timeout fixtures, name the retained evidence and
its local boundaries, and state the incomplete-run rule.

## Decision under test

How should Ex1 present its final guide-facing documentation so a reader can
identify the example's claims and provisional limits, reproduce a demo, and
inspect retained evidence without mistaking host-local runtime artifacts for a
portable PromiseGrid protocol contract?

## Assumptions and trust model

- Ex1 implements five named local-draft pCID profiles, not a frozen upstream
  PromiseGrid specification or interoperability claim. Source: DI-garis and
  DI-josir.
- The Docker demo begins with an empty host-local runtime root and preserves
  final artifacts after the run. Source: DI-rokol.
- Raw envelopes, message records, and local observations are evidence from the
  component that retained them. A timeout or local observation does not prove
  another agent's intent. Source: DI-vihoz, DI-riguz, DI-purum, and DI-zosiz.
- The existing README is the first-entry document; `CHANGELOG.md` states local
  implementation scope; `docs/testing.md` explains regression verification;
  and `docs/design.md` contains the profile inventory and detailed design.
- Alice is a guide reader running the local demo. Mallory may provide stale
  artifacts, alter a fixture, or mistake a collector record for a universally
  trusted protocol fact.

## Alternatives

### A. Extend the README with a concise guided demo and evidence walkthrough

Keep the README as the entry point. Add a small guide-facing section that links
to the scope declaration and testing guide, gives copyable demo commands, names
the preserved runtime root and relevant files, and states the local-evidence
boundary. Keep detailed contract material in the existing design and profile
documents.

### B. Create a separate `docs/demo.md` operator guide

Keep the README brief and place the complete run-and-inspection sequence in a
new guide, linked from the README.

### C. Rely on the existing README, design, scope declaration, and testing guide

Add only cross-links and expect readers to assemble the run and inspection
procedure from those documents.

## Scenario analysis

### Normal first run

Alice opens the README, runs the happy-path fixture, and sees the preserved
runtime root. Under A, she immediately finds the command, the scope link, and
the exact per-role/collector files worth inspecting. Under B, she must follow
one additional link but receives a fuller operator document. Under C, she can
run the demo but must infer which `message-cas`, JSONL, and analyzer artifacts
answer which question.

### Refusal and timeout interpretation

Alice runs `warehouse-refusal.json` or `carrier-timeout.json`. A can point her
to the seller's local `observations.jsonl` and explain that a signed refusal is
the responder's artifact while `timeout_observed` is only the seller's local
record. B can say the same with more room. C leaves the distinction scattered
between design and testing prose, increasing the risk that she treats a
timeout as proof that the carrier refused.

### Corruption, incomplete output, and failed runs

Mallory or a failed container leaves an incomplete runtime root. A can state
that each new run clears the selected root, that a nonzero script exit means no
completed scenario claim, and that the operator should rerun rather than merge
partial artifacts. B offers more space for troubleshooting. C provides no
single place that differentiates an incomplete run from preserved evidence of a
completed one.

### Mixed versions and long-horizon evolution

Alice changes a local profile document, changing its pCID. A sends her from the
README to the design inventory and local scope declaration, which explicitly
state the draft profile boundary. B does the same through a separate guide. C
risks making the README's demo instructions look like an interoperability
promise because the provisional status is not visible from the first-entry
document.

### Scale and maintenance

A adds a compact, maintained entry-point section and reuses the existing
specialized documents. B makes the operator flow easier to expand but creates
another document that can drift from the demo wrapper. C has no duplicate text
but imposes the highest navigation and interpretation cost on readers.

## Conclusions

C is rejected because it fails the required one-stop guide-facing outcome.
B remains viable if Ex1 later gains several operator modes or troubleshooting
procedures, but it adds a separate document before the current walkthrough is
large enough to need one.

A survives and is recommended: extend the README with a concise guided demo
and evidence-inspection section, link to `CHANGELOG.md`, `docs/design.md`, and
`docs/testing.md`, and keep detailed protocol and test explanations in those
existing documents. The walkthrough must call the artifacts host-local and
observer-local rather than protocol-wide claims.

## Decisions still requiring DF

1. **Guide layout:** extend the README with the concise walkthrough
   (recommended), create a separate `docs/demo.md`, or rely on cross-links
   only?
2. **Demo coverage:** show the happy path plus one refusal and one timeout
   fixture (recommended), or happy path only?
3. **Evidence walkthrough:** name the collector DAG, per-role raw-envelope
   directories, and local `observations.jsonl` with their evidence boundaries
   (recommended), or name only the runtime root?
4. **Incomplete-run guidance:** say that a failed script run is not a completed
   demo claim and must be rerun from a fresh selected root (recommended), or
   omit failure guidance?

## Implications for open work

- `lubav.7` can proceed after these DFs are locked and recorded in the `lubav`
  decision log.
- The resulting edit is documentation-only and must preserve the distinction
  between local draft profiles, host-local runtime layout, and pCID-defined
  contract meaning.
- The README must continue linking the exercise-local testing guide required by
  the repository testing policy.
