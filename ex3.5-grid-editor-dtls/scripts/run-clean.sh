#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
dir="$(cd "${script_dir}/.." && pwd)"

cd "${dir}"

compose_cmd() {
	if docker compose version >/dev/null 2>&1; then
		echo "docker compose"
		return 0
	fi
	if command -v docker-compose >/dev/null 2>&1; then
		echo "docker-compose"
		return 0
	fi
	return 1
}

if compose="$(compose_cmd)"; then
	:
else
	echo "no docker compose command found (need docker compose or docker-compose)" >&2
	exit 127
fi

echo "== reset docker runtime state =="
# Intent: Keep the demo reset script usable on hosts that still ship the
# standalone `docker-compose` binary instead of the `docker compose` plugin.
# Source: DI-samuv
if ${compose} down -v --remove-orphans; then
	echo "reset complete"
else
	status=$?
	echo "${compose} down failed with status ${status}" >&2
	exit "${status}"
fi

echo "== run compose regression =="
if ${compose} up --build; then
	echo "compose run complete"
else
	status=$?
	echo "${compose} up failed with status ${status}" >&2
	exit "${status}"
fi
