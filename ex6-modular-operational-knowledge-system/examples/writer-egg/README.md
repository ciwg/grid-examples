# Writer Egg Example

This is the smallest installed executable package in the repo that performs a
real durable write through the basket.

It exists to prove the installed-package contract, not to model a full domain
package. Source: `DI-rovum`; `DI-nupad`.

## Try It

From the ex6 repo root:

```bash
go run ./cmd/moks package install ./examples/writer-egg
go run ./cmd/moks writer create writer-1
go run ./cmd/moks relay export ./writer-relay.json
```
