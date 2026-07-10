# kernel_register

Signed app/kernel registration protocol for the `ex1-order-flow` example.

- Outer shape: `grid([42(pCID), payload, proof])`
- Payload carries the registering role ID and the pCIDs it promises to receive
- Proof signs the pCID-tagged payload signable view
- Payload never repeats the protocol name
