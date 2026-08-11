# Writer Package Example

This is the smallest installed executable package in the repo that performs a
real durable write through the runtime.

It exists to prove the installed-package contract, not to model a full domain
package. Its `writer.note.v1` family is defined by the immutable
[writer specification](protocols/writer-note-v1.md) and uses fixed pCID
`bafkreigwh6qript7zma7gu6fgxixmno2eglo3v2bhwpqr3dg5utiyagmca`.
Source: `DI-rovum`; `DI-nupad`; `DI-jusij`; `DI-solan`.

The executable requires `MOKS_RECORD_FIXTURE` to name the canonical record
generator when it is exercised in the repository test harness. It emits a
base64 `[][]byte` process wrapper containing an exact canonical Grid record;
the pCID is explicit and is not derived from `writer.note.v1` at runtime.

## Try It

From the ex6 repo root:

```bash
go run ./cmd/moks package install ./examples/writer-agent
go run ./cmd/moks writer create writer-1
go run ./cmd/moks relay export ./writer-relay.json
```
