# Package Family pCID Registry

Status: frozen baseline; append-only

This registry maps each first-party package-record family to the CIDv1 of its exact immutable specification bytes. It is the authoritative source for Ex6's fixed package pCIDs. A pCID identifies a shared wire-level semantic contract; it is not an identifier for a package, workflow, executable, or individual record. Source: DI-jusij.

## Workflow authors

Create a workflow by composing existing family pCIDs whenever its durable records have meanings already covered below. For example, a procedure workflow can use procedure, run, evidence, and approval records without adding a pCID.

Add a new pCID only when the workflow introduces a genuinely new interoperable durable-record meaning. Before doing so, create a new versioned immutable family specification in `package-families/`, calculate its CIDv1 from the exact file bytes, append its mapping here, and update the implementation claims and tests. Do not rewrite or delete existing specifications or mappings: retire a workflow or family by documenting deprecation, while retaining old entries for historical records.

Installed third-party packages follow the same rule but own their family specifications outside this first-party registry. A receiving node preserves unknown pCID records as exact bytes until it has a local interpreter and policy for that family.

## Runtime scope and recompilation

This registry covers only Ex6's built-in first-party family set. Creating, editing, enabling, disabling, or retiring a workflow that composes existing families does not change this registry and does not require recompiling Ex6. A workflow is a local composition of record contracts, not a new pCID by default.

An installed third-party package may introduce its own family specification and pCID without changing the Ex6 binary. Ex6 preserves unknown records as exact bytes; the installed package supplies any local validator or interpreter. Rebuilding Ex6 is required only when the built-in first-party family set or its built-in implementation changes.

Retiring a workflow stops new use but does not delete its historical records or any specification they reference.

## Frozen first-party mappings

| Family | Specification | Fixed CIDv1 pCID |
| --- | --- | --- |
| `moks.context.place.v1` | [moks-context-place-v1.md](package-families/moks-context-place-v1.md) | `bafkreibczv3ytaah2bgoelclio275nblrps4lrtwjvsa5rwo6r3aybndim` |
| `moks.context.resource.v1` | [moks-context-resource-v1.md](package-families/moks-context-resource-v1.md) | `bafkreidrnhbzqmhxstqq2ev325hwws6qeucaunumk7p7tbiwesnbs7rwma` |
| `moks.context.responsibility.v1` | [moks-context-responsibility-v1.md](package-families/moks-context-responsibility-v1.md) | `bafkreibwa7cc65rmie3z5tltubvz3uz34dinpsvm3nroczcb6exvaitidq` |
| `moks.correctiveaction.event.v1` | [moks-correctiveaction-event-v1.md](package-families/moks-correctiveaction-event-v1.md) | `bafkreicnlp64sd7obleleg5exm7c2yysg2wjv6wnmzx5nvx4uhpcftwvya` |
| `moks.inventory.count.v1` | [moks-inventory-count-v1.md](package-families/moks-inventory-count-v1.md) | `bafkreiftbdmcewooqz73ise7klyy7zpzrlqa2q5p5hpf4mcy5sesex33a4` |
| `moks.inventory.item.v1` | [moks-inventory-item-v1.md](package-families/moks-inventory-item-v1.md) | `bafkreiguud5rv5ybe7dpwmqq4xf2lkjmh5wajvpl2rfew2q5iu4ewathwm` |
| `moks.inventory.reconcile.v1` | [moks-inventory-reconcile-v1.md](package-families/moks-inventory-reconcile-v1.md) | `bafkreidxt63iq3zazcj6jdtuzkuwyyaht56fuustik6um7vyielpor7kiy` |
| `moks.knowledge.item.v1` | [moks-knowledge-item-v1.md](package-families/moks-knowledge-item-v1.md) | `bafkreigcuuwxlm3c2bvj73mvftzo4qbkyq5pqecx7zgtf4bnmdafpwjk64` |
| `moks.knowledge.lifecycle.v1` | [moks-knowledge-lifecycle-v1.md](package-families/moks-knowledge-lifecycle-v1.md) | `bafkreihffm26lnfpfzdn466pomobsmh2lyhl4kwfuuik3pzf4axm2ke2eu` |
| `moks.knowledge.revision.v1` | [moks-knowledge-revision-v1.md](package-families/moks-knowledge-revision-v1.md) | `bafkreidnyfrg5eifflrbjpb7kckaoxlwca77swygm5ae2cdo6wweccgs6e` |
| `moks.links.typed.v1` | [moks-links-typed-v1.md](package-families/moks-links-typed-v1.md) | `bafkreihhg4mxoh37frbzw22ujyuigu546qkn5pr5iyy4naxd6dwr3udsyu` |
| `moks.maintenance.finding.v1` | [moks-maintenance-finding-v1.md](package-families/moks-maintenance-finding-v1.md) | `bafkreif2pt5vcyijajwrurebyk2eoa4h77gxtj7iq7suwuw2egaogksw7u` |
| `moks.maintenance.item.v1` | [moks-maintenance-item-v1.md](package-families/moks-maintenance-item-v1.md) | `bafkreies7a4w3bpnl7beri3ofky3lewiztwkmalzkb3bu3mgkfh5vfklyi` |
| `moks.maintenance.service.v1` | [moks-maintenance-service-v1.md](package-families/moks-maintenance-service-v1.md) | `bafkreibsdabewazp4ibxcjp2uwyddczv2qim6dqyg7rckrzlzgpadvppi4` |
| `moks.ops.note.v1` | [moks-ops-note-v1.md](package-families/moks-ops-note-v1.md) | `bafkreia3ctr6uzj3ygnpdastrn4fydsvknnonkdpgxsydyj5qvz2qcs5gy` |
| `moks.procedures.item.v1` | [moks-procedures-item-v1.md](package-families/moks-procedures-item-v1.md) | `bafkreid3xltxowg4uftabtchtmqcdnr6vcmhxd7ziqmrxxe3v5ew6vu3yi` |
| `moks.procedures.use.v1` | [moks-procedures-use-v1.md](package-families/moks-procedures-use-v1.md) | `bafkreic3x6pdeeazo66unr637j5nib3lkp7y26kwrocltiwkjejxds6uve` |
| `moks.quarantine.event.v1` | [moks-quarantine-event-v1.md](package-families/moks-quarantine-event-v1.md) | `bafkreicqudjaew2khqckm7yasfsthrdwldps26wcsm7rjg6aswagjp7md4` |
| `moks.receiving.disposition.v1` | [moks-receiving-disposition-v1.md](package-families/moks-receiving-disposition-v1.md) | `bafkreifvlesjb2sqwvdholdrhpvttsmrrjka7as2ig3ijqm3ypnywnhfp4` |
| `moks.receiving.item.v1` | [moks-receiving-item-v1.md](package-families/moks-receiving-item-v1.md) | `bafkreihtq33ntlmhgu2jknbpehuq3tjcqh6y5ruaekhao3yl6ew32nsqbq` |
| `moks.receiving.receipt.v1` | [moks-receiving-receipt-v1.md](package-families/moks-receiving-receipt-v1.md) | `bafkreiar4u5irmomc2lkptxhrsdsk5podezhrji55247hhkh2ngzkcba7i` |
| `moks.runs.approval.v1` | [moks-runs-approval-v1.md](package-families/moks-runs-approval-v1.md) | `bafkreidbpikyfswoqv6uja2liswlao5t2jjpilklmkqqbgrl3fhtkpvpba` |
| `moks.runs.evidence.v1` | [moks-runs-evidence-v1.md](package-families/moks-runs-evidence-v1.md) | `bafkreiflfm2hhldipxdzmjzmnlufisf4tihzn7wjrdbdduasawqesvdnly` |
| `moks.runs.run.v1` | [moks-runs-run-v1.md](package-families/moks-runs-run-v1.md) | `bafkreia75oenisiarp2psot5tumzsw3snu5ofqzk5lzijirgs6r5jidox4` |
| `moks.training.completion.v1` | [moks-training-completion-v1.md](package-families/moks-training-completion-v1.md) | `bafkreigm66qcalyb3bedzjmnybbqwvn6cmpcjw32lsjertk6i6kdxolpcm` |
| `moks.training.item.v1` | [moks-training-item-v1.md](package-families/moks-training-item-v1.md) | `bafkreihsft37k2fefo7h3viuhzgmtdxy3k64lgypgnsvd63fbxtfyccq7y` |
| `moks.training.session.v1` | [moks-training-session-v1.md](package-families/moks-training-session-v1.md) | `bafkreic5pd5ijoxdwdqnbj6yiy2ufkm4sl6p4jqnfftyqezvcfo5updnee` |

## Verification

Recalculate every CIDv1 from the linked file's exact bytes before accepting a registry change. The future registry test must fail if a fixed value does not match its file or if a first-party family lacks an entry.
