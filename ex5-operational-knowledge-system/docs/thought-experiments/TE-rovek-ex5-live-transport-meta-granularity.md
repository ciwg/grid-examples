# ex5 live transport meta granularity

TE ID: TE-rovek
## Status
decided

## Decision under test

How `/api/meta` should describe live-draft transport preference now that ex5
ships different preferred live paths for different embodiments:

- browser prefers websocket carriage under the local HTTP adapter
- CLI uses no live-draft carriage of its own
- Neovim prefers the local Unix socket and keeps HTTP fallback

The current single `live_draft_preferred_transport` field compresses those
differences into one global answer that is no longer fully honest.

## Assumptions

- `/api/meta` remains the local capability surface for embodiment discovery.
- The fix should improve contract clarity, not add a new transport.
- Mixed-version clients matter because some consumers may already parse the
  current singular field.
- Mallory is not central here; the real risk is misleading clients and docs.

## Alternatives

### Alternative A: embodiment-aware live transport fields replace the global preference

Remove the singular global `live_draft_preferred_transport` field and replace
it with explicit embodiment-aware fields such as browser live transport and
terminal/Neovim live transport preferences.

### Alternative B: keep the global field and add embodiment-specific fields beside it

Preserve `live_draft_preferred_transport` for compatibility, but add new
embodiment-specific live transport fields and document the singular field as a
legacy coarse summary.

### Alternative C: keep the singular field and explain the nuance only in docs

Leave the runtime contract unchanged and rely on prose to explain that the
single field is only an approximation.

## Scenario analysis

### Scenario 1: browser capability discovery

Alice builds a browser-side tool that reads `/api/meta`.

Alternative A lets the browser see its own live carriage preference directly.

Alternative B also works, though the presence of both old and new fields adds
some transitional redundancy.

Alternative C leaves the browser with a misleading single field that may name
the terminal transport instead.

### Scenario 2: Neovim or terminal discovery

Bob wants a terminal client to know whether it should prefer local Unix socket
live carriage or HTTP fallback.

Alternative A tells him directly.

Alternative B also tells him directly while leaving the old summary field in
place for older consumers.

Alternative C remains ambiguous.

### Scenario 3: compatibility with existing tests or tooling

Carol already has scripts that read `live_draft_preferred_transport`.

Alternative A is the cleanest contract long-term, but it creates an immediate
breaking change for those consumers.

Alternative B allows a staged transition: new clients use the embodiment-aware
fields, while old clients still see the older coarse summary until it is no
longer needed.

Alternative C avoids breakage, but preserves the contract flaw.

### Scenario 4: PromiseGrid contract honesty

Dave wants `/api/meta` to describe the shipped embodiment split precisely.

Alternative A is the cleanest final shape because it removes the false idea of
one universal live-draft preference.

Alternative B is slightly less pure, but still honest if the old field is
clearly downgraded to compatibility/summary status.

Alternative C is no longer honest enough.

### Scenario 5: long-horizon evolution

Ellen later adds another embodiment or changes browser transport again.

Alternative A scales cleanly because each embodiment can carry its own
capability field.

Alternative B also scales, though eventually the coarse summary field becomes
more legacy baggage.

Alternative C grows less and less credible as embodiment diversity increases.

## Conclusions

Rejected:

- Alternative C. It leaves the contract under-specified.

Surviving:

- Alternative A: replace the global preference with embodiment-aware fields
- Alternative B: keep the global field as compatibility and add embodiment-
  aware fields beside it

Recommendation before DF:

- Alternative B

Why it survived pre-lock:

- It fixed the contract honesty problem while preserving compatibility.

## Implications for TODOs and pending DIs

- TODO `123` is locked to Alternative `A` instead: remove the singular field
  and publish embodiment-specific live transport fields directly.
- The shipped `Meta` contract should expose the browser and Neovim live-draft
  transport lanes explicitly and stop implying that one global live transport
  preference exists.
