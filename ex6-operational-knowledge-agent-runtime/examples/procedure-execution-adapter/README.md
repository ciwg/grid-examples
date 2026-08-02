# Procedure Execution Adapter

This is the first real installed workflow-adapter package. Its Docker worker
accepts the procedure-execution input handoff, returns the matching completed
output handoff, and proposes the built-in run and procedure-use records. The
runtime validates and signs those records; the worker has no runtime state
access. Source: `DI-fofuh`.

Build the local image from the ex6 root:

```bash
docker build -f examples/procedure-execution-adapter/Dockerfile -t moks/procedure-execution-adapter:dev .
```

The tag is only a local build convenience. The package manifest uses Docker's
immutable local image ID, so an installed adapter cannot silently move when a
tag is rebuilt. Rebuild this image and update the manifest plus `package.sh`
together before installing it on another machine. Source: `DI-fofuh`.

Then install the package into a disposable runtime with:

```bash
go run ./cmd/moks package install ./examples/procedure-execution-adapter
```
