#!/bin/sh
set -eu

repo_root="$(cd "$(dirname "$0")/.." && pwd)"
workdir="$(mktemp -d)"

printf 'workspace: %s\n' "$workdir"
cd "$workdir"

go run "$repo_root/cmd/moks" context place create place-1 Receiving Inbound-area
go run "$repo_root/cmd/moks" context resource create res-1 Scale Bench-scale place-1
go run "$repo_root/cmd/moks" procedures create proc-1 DockCheck Check-the-dock
go run "$repo_root/cmd/moks" procedures record-use proc-1 run-1 alice ok followed-v1
go run "$repo_root/cmd/moks" runs evidence add run-1 ev-1 photo kind=image,shift=night blob payload
go run "$repo_root/cmd/moks" runs approve run-1 ap-1 accepted looks-good
go run "$repo_root/cmd/moks" procedures inspect proc-1
go run "$repo_root/cmd/moks" runs inspect run-1
go run "$repo_root/cmd/moks" relay export "$workdir/relay.json"

printf 'relay export: %s\n' "$workdir/relay.json"
