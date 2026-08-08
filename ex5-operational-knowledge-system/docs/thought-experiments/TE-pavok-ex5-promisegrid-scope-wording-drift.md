# ex5 PromiseGrid scope wording drift

TE ID: TE-pavok
## Status
decided

## Decision under test

How the top-level ex5 PromiseGrid docs should describe the shipped scope now
that the runtime already includes:

- eight frozen signed families
- origin-aware ongoing peer exchange
- incremental relay-feed exchange
- a dedicated remote relay binary
- direct local Unix-socket terminal embodiments

The current problem is not missing runtime behavior. The problem is that some
summary docs still speak as if ex5 stops before signed envelopes or relay.

## Assumptions

- The implementation already ships signed-envelope runtime slices and relay
  behavior.
- Browser remains on the local HTTP adapter for the current shipped scope.
- Direct non-HTTP terminal embodiments and the dedicated remote relay are now
  real product/runtime behavior, not merely future intent.
- Mallory is irrelevant here except that inaccurate docs can create operator
  or reviewer misunderstanding about what the runtime actually promises.

## Alternatives

### Alternative A: describe the full shipped signed-envelope and relay scope directly

Rewrite the top-level PromiseGrid summary surfaces so they say clearly that
ex5 already ships its current signed-envelope, peer-exchange, relay-feed, and
remote relay layers, while still marking browser-side non-HTTP embodiment work
and broader ERP follow-ons as future scope.

### Alternative B: keep the “local runtime layer” framing and qualify later

Preserve the older statement that ex5 is mainly the local runtime / local
embodiment layer, then rely on later sections to explain that signed envelopes
and relay have in fact shipped within that local-first framing.

## Scenario analysis

### Scenario 1: new reader at the README

Alice reads only the README to decide whether ex5 already implements signed
families and relay behavior.

Alternative A tells her directly what ships today.

Alternative B makes the first impression misleading and depends on careful
later reading to correct it.

### Scenario 2: PromiseGrid review or external summary

Bob performs an alignment review and scans the top-level claims docs.

Alternative A gives him one consistent story: signed runtime families and
relay behavior are shipped within the current ex5 scope.

Alternative B invites false negatives because the summary wording understates
what the code already does.

### Scenario 3: future-scope boundaries

Carol wants to know what is still not shipped.

Alternative A can state that precisely: browser-side non-HTTP embodiment
contracts and other broader future-scope expansions remain open.

Alternative B muddies that boundary by grouping shipped signed-envelope and
relay features together with genuinely unshipped work.

### Scenario 4: long-horizon maintenance

Dave updates docs after future follow-on work.

Alternative A creates a stable baseline where the summary surfaces match the
runtime and later diffs only need to describe real new scope.

Alternative B preserves a stale rhetorical frame that future editors must keep
explaining around.

## Conclusions

Rejected:

- Alternative B. It is now too misleading for the shipped runtime.

Surviving:

- Alternative A: describe the full shipped signed-envelope and relay scope
  directly

Recommendation:

- Alternative A

Why:

- It is the most honest PromiseGrid wording for the code already shipping.
- It removes a recurring source of review churn.
- It keeps future-scope boundaries explicit without understating the present.

## Implications for TODOs and pending DIs

- TODO `122` is locked to Alternative `A`.
- The implementation should sweep the main summary surfaces and re-run a
  wording alignment pass to confirm no stale “pre-relay” framing remains in
  those top-level docs.
