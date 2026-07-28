#!/bin/sh
set -eu

case "${1-}" in
  describe)
    cat <<'EOF'
{
  "id": "replace-me",
  "version": "0.1.0",
  "description": "Example ex6 package",
  "commands": [
    {
      "path": ["replace", "me"],
      "summary": "Example command"
    }
  ],
  "families": [
    {
      "name": "replace.me.v1",
      "protocol_pcid": "pcid:replace.me.v1"
    }
  ],
  "claims": [
    {
      "protocol_pcid": "pcid:replace.me.v1",
      "role": "family-validator",
      "summary": "Validates replace.me.v1 envelopes."
    }
  ]
}
EOF
    ;;
  validate)
    body="$(cat)"
    case "$body" in
      *'"family":"replace.me.v1"'*) exit 0 ;;
      *) echo "unsupported family" >&2; exit 1 ;;
    esac
    ;;
  run)
    if [ "${2-}" != "replace me" ]; then
      echo "unknown command" >&2
      exit 1
    fi
    printf '%s\n' "replace-me package executed"
    ;;
  *)
    echo "usage: $0 {describe|validate|run}" >&2
    exit 1
    ;;
esac
