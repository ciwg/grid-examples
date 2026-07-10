#!/usr/bin/env bash
set -euo pipefail

# Intent: Each demo run starts from an empty temp root, but the resulting
# artifacts remain on disk afterward for manual inspection. Source: DI-rokol

script_dir="$(cd "$(dirname "$0")" && pwd)"
data_root="${EX1_DATA_ROOT:-/tmp/grid-examples-ex1-data}"
fixture="${1:-happy-path.json}"

cd "$script_dir"

compose_cmd=()
if command -v docker-compose >/dev/null 2>&1; then
  compose_cmd=(docker-compose)
else
  compose_cmd=(docker compose)
fi

"${compose_cmd[@]}" down

cleanup_data_root() {
  if [ ! -d "$data_root" ]; then
    return 0
  fi
  if rm -rf "$data_root" 2>/dev/null; then
    return 0
  fi
  # Intent: If an earlier container run left root-owned files behind, fall back
  # to a one-shot container cleanup so each new run still starts from an empty
  # temp root without deleting the preserved results afterward. Source: DI-rokol; DI-sabol
  docker run --rm -v "$data_root:/data" alpine:3.22 sh -c 'rm -rf /data/* /data/.[!.]* /data/..?*'
  rm -rf "$data_root"
}

cleanup_data_root

mkdir -p \
  "$data_root/collector" \
  "$data_root/intake" \
  "$data_root/seller" \
  "$data_root/warehouse" \
  "$data_root/accounting" \
  "$data_root/carrier"

export EX1_DATA_ROOT="$data_root"
export HOST_UID="$(id -u)"
export HOST_GID="$(id -g)"

"${compose_cmd[@]}" build
"${compose_cmd[@]}" up -d collector kernel seller warehouse accounting carrier
sleep 2
"${compose_cmd[@]}" run --rm intake "/app/fixtures/$fixture"
"${compose_cmd[@]}" run --rm analyze /data
"${compose_cmd[@]}" down

echo "preserved run data: $data_root"
