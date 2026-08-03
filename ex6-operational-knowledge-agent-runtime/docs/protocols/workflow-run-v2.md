# OKAR Workflow Run Protocol v2

Status: implementation-local profile; current selector is
`bafkreiejstdppcv25xazlzgszwrnjrf6rmnn3sosomkwhushkml3weit64`.

## Scope

This document defines the pCID-selected CBOR event shape for local workflow
runs. Version 2 adds one durable UTC event timestamp so an operator projection
can order current runs by their actual accepted event time. Source: DI-gihor.

The runtime continues to accept retained version-1 events with selector
`bafkreifmttp5fwt3yvxvkb7ni6kwg3j3arl7mbjsyzszf7s7crxrncch24`; those events
have no timestamp and sort after timestamped v2 events. New events always use
the v2 selector. This is an implementation-local profile, not a claim that
this Markdown file has completed a self-hashing publication process.

## Envelope

Each v2 event is the exact canonical CBOR bytes of:

```cddl
grid(
  [ selector: 42(pCID),
    state: tstr,
    run_cid: bstr,
    workflow: tstr,
    input_cid: bstr,
    output_cid: bstr,
    reason: tstr,
    parents: [* bstr],
    recorded_at_unix_nano: uint,
  ]
)
```

`recorded_at_unix_nano` is the accepting runtime's UTC Unix-nanosecond time.
It orders local operator projections; it does not establish distributed clock
truth, causal order, authorization, or a peer timestamp claim.

## Acceptance and compatibility

1. New events require the v2 selector, eight slots, canonical CBOR, and a
   non-zero timestamp representable as a Go UTC time.
2. Retained v1 events require their original selector and seven slots. Their
   projected `updated_at` is absent and they remain inspectable and executable.
3. Replay validates parent, run, workflow, handoff, and transition rules before
   a timestamped event can become a current run head.
4. Overview sorts current run heads newest-first by `updated_at`, then run ID
   for a deterministic tie break; v1 zero-time heads follow v2 heads.
