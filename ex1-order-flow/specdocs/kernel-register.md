# kernel_register

## Status and scope

This is an Ex1 local draft/example profile, not a frozen upstream PromiseGrid
specification. It lets an app role declare which local profile pCIDs it receives
from the kernel.

## Envelope, payload, and proof

The outer shape is `grid([42(pCID), payload, proof])`. Slot 1 is a CBOR
`KernelRegisterMessage` containing a non-empty `role` and `receive_pcids` as
binary CID bytes. Slot 2 is an EdDSA COSE proof over the pCID-tagged slot
vector. The kernel verifies the proof against the declared local role before
recording subscriptions.

## Validation and response

Registration has no capability-token field. Invalid envelope bytes, an
unexpected pCID, malformed payload, invalid proof, or invalid CID bytes are
locally rejected by the kernel without a required signed reply. A valid
registration permits only pCID-selected exact-byte dispatch; it grants neither
business authority nor a claim about whether any downstream promise is kept.

## Evidence and evolution

The profile describes a transient local kernel subscription. Business messages
remain durable as the signed raw artifacts retained by their sending and
receiving agents and by the collector. Editing this local draft changes its
pCID; peers using the older pCID require the older profile bytes.
