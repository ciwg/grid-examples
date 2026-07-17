#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
dir="$(cd "${script_dir}/.." && pwd)"

cd "${dir}"

echo "== reset docker runtime state =="
if docker compose down -v --remove-orphans; then
	echo "reset complete"
else
	status=$?
	echo "docker compose down failed with status ${status}" >&2
	exit "${status}"
fi

echo "== run compose regression =="
if docker compose up --build; then
	echo "compose run complete"
else
	status=$?
	echo "docker compose up failed with status ${status}" >&2
	exit "${status}"
fi

