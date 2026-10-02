#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exercise_dir="$(cd "${script_dir}/.." && pwd)"

gum_bin=""
if command -v gum >/dev/null 2>&1; then
	gum_bin="$(command -v gum)"
elif command -v go >/dev/null 2>&1; then
	candidate="$(go env GOPATH)/bin/gum"
	if [ -x "${candidate}" ]; then
		gum_bin="${candidate}"
	fi
fi

if [ -z "${gum_bin}" ]; then
	echo "Grid Charm topology chooser requires Charm Gum. Install it, then run this script again:" >&2
	echo "  go install github.com/charmbracelet/gum@latest" >&2
	exit 127
fi

# Intent: Let a presenter select an existing Ex3 topology without changing the
# relay's signed-envelope, pCID, or admission behavior. Source: DI-holoz
choice="$("${gum_bin}" choose 'One local relay (7025)' 'Two relays (7025 and 7026)')"
case "${choice}" in
	'One local relay (7025)')
		cd "${exercise_dir}"
		go run ./cmd/grid-relay --listen 127.0.0.1:7025 --data-root .grid-editor/charm
		;;
	'Two relays (7025 and 7026)')
		cd "${exercise_dir}"
		docker compose up -d --build
		;;
	*)
		echo "No topology selected." >&2
		exit 1
		;;
esac
