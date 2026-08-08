# How OKAR Works

OKAR is a node-local runtime for operational-knowledge packages. The runtime
owns shared infrastructure; packages own domain behavior.

1. A built-in or installed package declares its commands, protocol families,
   and implementation claims.
2. At startup the runtime validates the package and derives a route table from
   those claims. A route can be direct, parser, or transform.
3. The operator can inspect the candidate routes and the runtime's preferred
   plan. Route policies make the selection deterministic and locally
   configurable.
4. Durable records are validated by their owning package when the family is
   known. Unknown families are retained as exact bytes for relay rather than
   discarded.
5. The runtime keeps append-only history and content-addressed blobs in CAS.
   Relay exchange carries records, route metadata, signatures, and proofs;
   discovery alone never grants trust.
6. Workflow lifecycle events are stored as pCID-selected grid CBOR bytes in
   CAS. Import records an artifact; activation, deactivation, and revocation
   are local lifecycle decisions. The JSON cache is disposable and rebuilt from
   CAS on startup.

## What it is today

OKAR has a CLI, package activation, route planning, durable records, CAS,
relay, local trust policy, and nine first-party packages. It is still an
intermediate runtime: packages are not yet independently running agents, and
route plans are not yet executed through a full worker-dispatch system.

Source: DI-moksu, DI-puvok, DI-bavuk.

## Meeting walkthrough: workflow loading

Run:

    go run ./cmd/moks workflow demo procedure-execution

The disposable demo proves four things:

1. The procedure-execution directory was captured as one immutable CAS artifact.
2. Its manifest and local package dependencies were verified.
3. The workflow was explicitly activated for local availability.
4. The exact retained artifact was extracted for inspection.

The final JSON status should show an active state and ready set to true. The
demo uses a temporary runtime and removes it afterward, so it does not change
the local runtime under .moks/.
