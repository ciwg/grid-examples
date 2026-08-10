# accounting

## Status and scope

This is an Ex1 local draft/example profile, not a frozen upstream PromiseGrid
specification. It carries seller-to-accounting requests and accounting-to-seller
results.

## Envelope and proof

The outer shape is `grid([42(pCID), payload, proof])`. Slot 1 is a CBOR
`AccountingMessage`; slot 2 is an EdDSA COSE proof over the pCID-tagged slot
vector. Seller verifies accounting results; accounting verifies seller requests.

## Payload and valid kinds

`kind` is required and is either `request` or `result`.

- `request` requires `customer_order_ref`, `parent_order_cid`, and
  `capability_token`; payment reference, currency, and amount are optional.
- `result` requires `customer_order_ref`, `parent_order_cid`,
  `capability_token`, and `status` (`recorded` or `refused`). It may carry
  ledger entry, invoice reference, and notes.

## Capability and validation

Each sender issues a signed bearer token to its receiver. It binds issuer,
audience, this pCID, kind, action `send`, issued-at, expiry, and token ID. It
is reusable until its five-minute expiry; no redemption or replay ledger exists.
Unexpected pCID, malformed bytes, invalid proof, or invalid/expired token are
locally rejected without a required signed reply. A valid request receives one
signed result on the profile-defined path; `refused` is a business result.

## Evidence and evolution

Raw signed envelopes are retained locally and sent envelopes are published with
parent links to the collector. `parent_order_cid` lets later readers relate the
accounting result to the order. Editing this local draft changes its pCID.
