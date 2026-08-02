# Local registry allow-list policy shape

TE ID: TE-malob

## Status

decided

## Decision under test

How an EX6 runtime stores and evaluates its local OCI registry allow-list
before a digest-pinned workflow-adapter image may be acquired.

This refines TE-sobir / DI-harib. It does not decide Docker pull mechanics,
credentials, registry authentication, or image garbage collection.

## Assumptions and trust model

- Alice may publish a package manifest naming a digest-pinned OCI image.
- Bob owns the local runtime, its network policy, and its registry trust
  decisions. Alice's package cannot enlarge Bob's allowed network destinations.
- Image references are registry-qualified `host/repository@sha256:<digest>`
  references. Bare Docker image IDs and mutable tags are not portable registry
  dependencies.
- The policy must persist across runtime restart and remain inspectable by a
  local operator.
- No registry credential or token is stored in the runtime in this first slice.

## Alternatives

1. **Exact-host runtime allow-list.** Persist a sorted set of canonical
   registry hosts, including an explicit port when one is used. A reference is
   allowed only when its parsed registry host exactly equals one entry.
2. **Suffix/wildcard allow-list.** Permit entries such as `*.example.com` or
   `.example.com` to authorize multiple registry hosts.
3. **Per-package registry grants.** Persist policy keyed by package ID plus
   registry host, rather than one runtime-wide host set.

## Scenario analysis

### Normal operation

Bob runs `moks registry allow registry.example.com`. A package declaring
`registry.example.com/moks/procedure@sha256:...` becomes eligible for the
later image-availability check. `moks registry list` shows the stored host,
and `moks registry remove registry.example.com` removes permission without
deleting installed packages or retained workflow artifacts.

Option 2 reduces setup for a registry fleet but makes a short policy entry
authorize hosts Bob may not have reviewed. Option 3 lets Bob grant one package
a narrower exception, but package IDs are locally installed metadata that can
be replaced or versioned while the network trust question is simply whether a
host is allowed.

### Failure and malformed references

Exact-host matching rejects an image reference with no explicit registry host,
a tag, an invalid digest, an empty host, or a host not in the set. A failed
policy write leaves the in-memory allow-list unchanged. A deleted cache does
not matter because the policy file is durable local configuration.

Suffix matching needs normalization rules for dots, ports, case, public-suffix
boundaries, and lookalike hosts such as `registry.example.com.attacker.test`.
Per-package grants need package-version and uninstall/reinstall semantics.

### Concurrent actors and mixed versions

If Bob removes a host while an old artifact remains active, subsequent starts
must become not-ready rather than silently using the previous image. The
artifact and package remain retained for audit. A newer runtime that does not
understand the allow-list must not assume it grants execution; it should report
the portable image dependency unavailable.

Exact-host entries compose predictably across nodes. Wildcards make different
operators' interpretation of the same policy harder to audit. Per-package
grants make peer diagnosis require comparing package installation state as
well as network policy.

### Long-horizon evolution and scale

An exact-host set is small, deterministic, easy to serialize alongside current
local policy, and can later grow additive metadata such as a transport rule or
per-host credential reference. A future per-package exception can be added as
an explicit second policy layer without changing the baseline meaning of host
approval.

Wildcard matching is difficult to tighten later because existing entries may
already grant broad authority. Per-package grants are premature until EX6 has
package provenance/version policy beyond the current self-check.

## Conclusions

Reject suffix/wildcard matching for the first slice: its convenience does not
outweigh ambiguous and overly broad network authority.

Defer per-package grants: they add lifecycle and replacement semantics before
the base host-trust policy is proven.

Exact-host runtime allow-list survives. Store canonical host[:port] strings in
the existing runtime-owned policy state; expose `moks registry allow <host>`,
`moks registry list`, and `moks registry remove <host>`; reject empty hosts,
paths, schemes, wildcard characters, and non-canonical duplicate entries.

## Decision status

Locked by DI-hapak. The CLI family is `moks registry allow|list|remove`; its
runtime-wide scope is intentional because it governs local network trust rather
than one package's lifecycle.

## Implications for open work

- Persist the set beside existing local policy state, using atomic write and
  restart tests.
- Parse the registry host from an OCI digest reference before policy matching;
  never compare raw string prefixes.
- Include registry-policy readiness in workflow verification before Docker pull
  behavior is added.
- Do not add credentials, wildcard matching, or image transfer to this slice.
