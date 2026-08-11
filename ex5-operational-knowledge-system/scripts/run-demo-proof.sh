#!/usr/bin/env bash
set -eu

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
demo_socket="/tmp/ex5-demo-browser/runtime/embodiment.sock"

# Intent: Rehearse the one verified browser-first demo session, then print a
# read-only CLI proof against the exact same runtime without widening into a
# recording or second-demo toolkit. Source: DI-pudob.
"${repo_root}/scripts/setup-demo-browser.sh"
"${repo_root}/scripts/launch-demo-browser.sh"
"${repo_root}/scripts/verify-demo-browser.sh"

echo
echo "Browser demo is verified. Run this CLI proof against the same runtime:"
echo "go run ./cmd/oks-cli -socket ${demo_socket} pending-review"
