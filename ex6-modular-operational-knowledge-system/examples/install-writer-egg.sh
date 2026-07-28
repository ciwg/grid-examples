#!/bin/sh
set -eu

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
workdir="$(mktemp -d)"

printf 'workspace: %s\n' "$workdir"
cd "$workdir"

go run "$repo_root/cmd/moks" package install "$repo_root/examples/writer-egg"
go run "$repo_root/cmd/moks" writer create writer-1
go run "$repo_root/cmd/moks" relay export "$workdir/writer-relay.json"

printf 'relay export: %s\n' "$workdir/writer-relay.json"
