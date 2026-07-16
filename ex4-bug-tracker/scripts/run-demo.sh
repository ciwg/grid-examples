#!/usr/bin/env bash
set -euo pipefail

# Intent: Keep the ex4 demo repeatable by starting from a fresh temp runtime
# root each time while still using the normal server entrypoint and seeded app
# methods. Source: DI-zogof

script_dir="$(cd "$(dirname "$0")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
data_root="${EX4_DATA_ROOT:-/tmp/grid-examples-ex4-demo}"
listen="${EX4_LISTEN:-127.0.0.1:7035}"

if [ -d "$data_root" ]; then
  rm -rf "$data_root"
fi

mkdir -p "$data_root"

cd "$repo_root"
exec go run ./cmd/bug-tracker --listen "$listen" --data-root "$data_root" --seed-demo
