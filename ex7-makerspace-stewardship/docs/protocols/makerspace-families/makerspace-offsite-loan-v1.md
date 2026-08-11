# Makerspace Off-Site Loan v1

Status: frozen

Family: `makerspace.offsite.loan.v1`

This family is a borrower's voluntary return commitment for one loanable tool.
It uses `../makerspace-record-v1.md`; its pCID is the CIDv1 of this exact
file's bytes.

## Payload

The canonical JSON object has exactly non-empty `tool_id`, `borrower_id`,
`due_at`, `policy_version`, and `policy`. `due_at` is a UTC RFC3339 return
commitment time. `borrower_id` must equal the common envelope signer label at
Ex7 local admission. `policy_version` and `policy` are the exact version and
text presented at acceptance, never reconstructed from later area policy.

## Meaning and local assessment

The signer promises to return the identified tool by the stated time under the
recorded accepted terms. This is evidence of a voluntary commitment, not a
global checkout command or universal authorization grant. Ex7 may display an
active loan only when its local loanability and qualification policy passes; a
later policy change must not rewrite accepted terms.

## Versioning

This document is immutable. A different commitment or terms model requires a
new versioned family specification and pCID.
