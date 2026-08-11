# Workflow and evidence guide

Equipment observation, safety disposition, off-site loan, and off-site return
are separate pCID-selected record families. An observation and a requested
safety hold are distinct records. A loan retains borrower, deadline, and
accepted policy snapshot; a return links the loan record it addresses.

The agent first verifies exact record bytes and signatures, then retains them.
It projects a known family only when the record's full public-key fingerprint
matches its local recognition policy. A valid unknown pCID or unrecognized key
is retained without known makerspace meaning. A recognized key without the
locally recognized steward role cannot clear a safety hold. Source: DI-tohak;
DI-piruf.

The displayed Alice/Carol/Dave names are local projection/bootstrap data. They
are not login credentials or author evidence. Signed participant history and
exact-byte carriage are implemented; browser signing, account-based authoring,
running relay/discovery, blob retrieval, and portable governance remain
outside this embodiment.

For a real browser evidence run, use
`scripts/run-two-agent-browser-proof.sh`. It shows Bob's unsigned terminal
request, Alice's local approval, and Bob's final projection of Alice's returned
signed record. The runner prints its disposable evidence directory, containing
`bob-final.png`, `alice-approval-response.json`, and the agent and Chrome
logs. Source: DI-fuzar; DI-hibok; DI-kasaz.
