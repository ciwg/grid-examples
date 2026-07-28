#!/bin/sh
set -eu

case "${1-}" in
  describe)
    cat <<'EOF'
{"id":"writer-egg","version":"0.1.0","description":"Example installed package that writes through the basket","commands":[{"path":["writer","create"],"summary":"Create a writer note"}],"families":[{"name":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1"}],"claims":[{"protocol_pcid":"pcid:writer.note.v1","role":"family-validator","summary":"Validates writer note envelopes."}]}
EOF
    ;;
  validate)
    body="$(cat)"
    case "$body" in
      *'"family":"writer.note.v1"'*) exit 0 ;;
      *) echo "wrong family" >&2; exit 1 ;;
    esac
    ;;
  run)
    if [ "${2-}" != "writer create" ]; then
      echo "unknown writer command" >&2
      exit 1
    fi
    if [ "${3-}" = "" ]; then
      echo "missing record id" >&2
      exit 1
    fi
    cat <<EOF
{"output":"created $3","cas":[{"alias":"body1","body":"payload for $3"}],"records":[{"family":"writer.note.v1","protocol_pcid":"pcid:writer.note.v1","record_id":"$3","signer":"writer-egg","timestamp":"2026-07-28T00:00:00Z","payload":{"title":"Writer","body_ref":"\$cas:body1"}}]}
EOF
    ;;
  *)
    echo "unknown verb" >&2
    exit 1
    ;;
esac
