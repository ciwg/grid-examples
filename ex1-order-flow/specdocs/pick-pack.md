# pick_pack

## Status and scope

This is an Ex1 local draft/example profile, not a frozen upstream PromiseGrid
specification. It carries seller-to-warehouse pick-and-pack requests and
warehouse-to-seller results.

## Envelope and proof

The outer shape is `grid([42(pCID), payload, proof])`. Slot 1 is a CBOR
`PickPackMessage`; slot 2 is an EdDSA COSE proof over the pCID-tagged slot
vector. Seller verifies warehouse results; warehouse verifies seller requests.

## Payload and valid kinds

`kind` is required and is either `request` or `result`.

- `request` requires `customer_order_ref`, `parent_order_cid`, `items`, and
  `capability_token`; service level is optional.
- `result` requires `customer_order_ref`, `parent_order_cid`,
  `capability_token`, and `status` (`packed` or `refused`). It may carry
  package ID, weight, package count, and notes.

## Capability and validation

Each sender issues a signed bearer token to its receiver. It binds issuer,
audience, this pCID, kind, action `send`, issued-at, expiry, and token ID. It
is reusable until its five-minute expiry; no redemption or replay ledger exists.
Unexpected pCID, malformed bytes, invalid proof, or invalid/expired token are
locally rejected without a required signed reply. A valid request receives one
signed result on the profile-defined path; `refused` is a business result, not
proof that the warehouse broke its protocol promise.

## Evidence and evolution

Raw signed envelopes are retained locally and sent envelopes are published with
their parent links to the collector. `parent_order_cid` ties the result to the
order request. Editing this local draft changes its pCID and creates a distinct
profile for later messages.
