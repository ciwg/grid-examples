# accounting

Signed accounting request/result protocol for the `ex1-order-flow` example.

- Outer shape: `grid([42(pCID), payload, proof])`
- Payload carries `kind = "request"` or `kind = "result"`
- Payload carries capability token bytes in `capability_token`
- Proof signs the pCID-tagged payload signable view
- Payload never repeats the protocol name
