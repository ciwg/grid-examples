# shipment

## Status and scope

This is an Ex1 local draft/example profile, not a frozen upstream PromiseGrid
specification. It carries seller-to-carrier booking requests and carrier-to-seller
results.

## Envelope and proof

The outer shape is `grid([42(pCID), payload, proof])`. Slot 1 is a CBOR
`ShipmentMessage`; slot 2 is an EdDSA COSE proof over the pCID-tagged slot
vector. Seller verifies carrier results; carrier verifies seller requests.

## Payload and valid kinds

`kind` is required and is either `request` or `result`.

- `request` requires `customer_order_ref`, `parent_order_cid`,
  `parent_pick_pack_cid`, `parent_accounting_cid`, `capability_token`, and
  package details; destination and service level are optional.
- `result` requires `customer_order_ref`, all three parent CIDs,
  `capability_token`, and `status` (`booked` or `refused`). It may carry carrier,
  tracking, label-artifact, and note fields.

## Capability and validation

Each sender issues a signed bearer token to its receiver. It binds issuer,
audience, this pCID, kind, action `send`, issued-at, expiry, and token ID. It
is reusable until its five-minute expiry; no redemption or replay ledger exists.
Unexpected pCID, malformed bytes, invalid proof, or invalid/expired token are
locally rejected without a required signed reply. A valid request receives one
signed result on the profile-defined path. A carrier timeout is observed by the
seller; its durable observer-evidence form is outside this profile's current
scope.

## Evidence and evolution

Raw signed envelopes are retained locally and sent envelopes are published with
all parent links to the collector. The parent CIDs relate a result to its order,
pick-pack, and accounting work. Editing this local draft changes its pCID.
