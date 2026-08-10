# Ex1 implementation scope

## 2026-08-10 — Local draft profile scope declaration

This is a scope declaration for the Ex1 example, not a formal PromiseGrid
implementation-promise claim against a frozen upstream specification doc-CID.

### Local profiles implemented

| Local draft profile | Current pCID | Implemented components |
| --- | --- | --- |
| `order` | `bafkreiadgrunjzobfflgs4m4f7xzisuehhgd2s3w2ptl5xwvntis6ir4oy` | Intake and seller handlers; signed submit/final flow. |
| `pick_pack` | `bafkreiafy57gl5b2pn4pj5qkczund3jrp42tqq46rzb6ipw5wyxi6hlwoy` | Seller and warehouse handlers; signed request/result flow. |
| `accounting` | `bafkreib3bo5sdzsk276qnxtlppyyqfwjuafx33jqotc35mm53ucayrmon4` | Seller and accounting handlers; signed request/result flow. |
| `shipment` | `bafkreieivmdty4siq2smgjg27x62mgv5g5oqdr7h2l4hvbw46fk6z6vaki` | Seller and carrier handlers; signed request/result flow. |
| `kernel_register` | `bafkreidscjgld22lrvsr5wk3lwlzspnifv4e2ewsgsqsg3old3q3pvqpre` | Each app role and the kernel registration/dispatch boundary. |

The scope covers Ex1's pCID-selected `grid([42(pCID), payload, proof])`
envelopes, signed proof verification, the documented reusable-until-expiry
bearer-token checks, local agent handlers, exact-byte kernel dispatch, and raw
message artifact retention. It also covers each observer's local append-only
evidence for malformed input, invalid proof or capability, signed refusal,
timeout, and no-current-recipient dispatch. Each pCID is derived from the
named local draft profile bytes; changing those bytes creates a different
profile identity.

### Explicit non-claims

Ex1 does not claim:

- conformance to a frozen upstream PromiseGrid spec doc-CID;
- interoperability with an independently implemented peer solely because it
  uses a similar profile name;
- production role or signing-key continuity, key rotation, or external trust
  policy beyond deterministic demo fixture keys;
- one-time token redemption, replay prevention, exactly-once execution, or a
  durable token-redemption ledger;
- that a local observer-evidence record establishes another agent's intent,
  produces shared evidence, or constitutes a portable protocol artifact; or
- that the Docker topology, collector, analyzer, local filesystem layout, or
  runtime scheduling model is part of a portable protocol contract.

When Ex1 adopts a frozen upstream spec, a later entry will use the guide's
formal implementation-promise fields (`claim`, `spec`, `scope`,
`breaking-change`, and `notes`) and name that exact frozen spec doc-CID.
