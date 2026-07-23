#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
source_root="${repo_root}/sample-data/newcomer-runtime"
target_root="${1:-}"

if [ -z "${target_root}" ]; then
  echo "usage: scripts/load-sample-data.sh TARGET_ROOT" >&2
  exit 1
fi

if [ ! -d "${source_root}" ]; then
  echo "sample corpus missing: ${source_root}" >&2
  exit 1
fi

if [ -e "${target_root}" ]; then
  if [ -d "${target_root}" ]; then
    if find "${target_root}" -mindepth 1 -print -quit | grep -q .; then
      echo "target root must not already contain data: ${target_root}" >&2
      exit 1
    fi
  else
    echo "target root exists and is not a directory: ${target_root}" >&2
    exit 1
  fi
else
  mkdir -p "${target_root}"
fi

# Intent: Keep the newcomer sample corpus checked in and deterministic while
# making the load step explicit and fail-closed instead of mutating an existing
# runtime in place. Source: DI-rubav
cp -a "${source_root}/." "${target_root}/"
echo "loaded ex5 newcomer sample data into ${target_root}"
