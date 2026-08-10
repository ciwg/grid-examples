# order

## Status and scope

This is an Ex1 local draft/example protocol profile. Its pCID is derived from
these exact bytes and is listed in `docs/design.md`; it is not a frozen upstream
PromiseGrid specification or an interoperability claim.

The profile carries an order submission from intake to seller and one final
seller result back to intake.

## Envelope and proof

The outer shape is `grid([42(pCID), payload, proof])`. Slot 1 is a CBOR
`OrderMessage`; slot 2 is an EdDSA COSE proof over the pCID-tagged slot vector.
The payload does not repeat the protocol name. Intake verifies seller finals;
seller verifies intake submissions.

## Payload and valid kinds

`kind` is required and is either `submit` or `final`.

- `submit` requires `customer_order_ref`, at least one `items` entry, and
  `capability_token`; it may carry requester, shipping, payment, service, and
  note fields.
- `final` requires `customer_order_ref`, `parent_order_cid`,
  `capability_token`, and `order_status` (`fulfilled`, `refused`, or `failed`).
  It may carry a failure stage, warehouse/accounting/shipment result CIDs,
  package/tracking/ledger identifiers, and a summary.

## Capability and validation

For each message, the sender issues a signed bearer token to the receiving
role. The token binds issuer, audience, this profile pCID, message kind, action
`send`, issued-at, expiry, and token ID. It is reusable until its five-minute
expiry; Ex1 has no redemption or replay ledger. Seller accepts `submit` only
when pCID, proof, token issuer/audience/kind/action, and expiry validate.
Intake applies the corresponding checks to `final`.

Malformed bytes, an unexpected pCID, invalid proof, or invalid/expired token
are rejected locally and do not require a signed reply. A valid `submit` may
receive a signed `final`: business refusal can be `order_status = refused`,
while validation or downstream execution failure can be `failed`.

## Evidence and evolution

Senders and recipients retain raw signed envelopes locally; the sender records
parent CIDs and publishes raw sent envelopes to the collector. `parent_order_cid`
and result CIDs make the final auditable against the earlier work. Changing
this file changes its pCID; old envelopes remain evidence for the old profile.
