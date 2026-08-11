#!/bin/sh
set -eu

if [ -z "${MOKS_RECORD_FIXTURE:-}" ]; then
  echo "MOKS_RECORD_FIXTURE must name the canonical record generator" >&2
  exit 1
fi

case "${1-}" in
  describe)
    cat <<'EOF'
{"id":"writer-agent","version":"0.1.0","description":"Example installed package that writes through the runtime","commands":[{"path":["writer","create"],"summary":"Create a writer note"}],"families":[{"name":"writer.note.v1","protocol_pcid":"bafkreigwh6qript7zma7gu6fgxixmno2eglo3v2bhwpqr3dg5utiyagmca"}],"claims":[{"protocol_pcid":"bafkreigwh6qript7zma7gu6fgxixmno2eglo3v2bhwpqr3dg5utiyagmca","role":"family-validator","summary":"Validates writer note envelopes."}]}
EOF
    ;;
  validate)
    family="$("$MOKS_RECORD_FIXTURE" inspect)"
    if [ "$family" != "writer.note.v1" ]; then
      echo "wrong family" >&2
      exit 1
    fi
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
    # Intent: External package records name their frozen pCID explicitly. Source: DI-solan.
    record="$("$MOKS_RECORD_FIXTURE" --pcid bafkreigwh6qript7zma7gu6fgxixmno2eglo3v2bhwpqr3dg5utiyagmca writer.note.v1 "$3" writer-agent 2026-07-28T00:00:00Z '{"title":"Writer","body_ref":"$cas:body1"}')"
    printf '{"output":"created %s","cas":[{"alias":"body1","body":"payload for %s"}],"records":["%s"]}\n' "$3" "$3" "$record"
    ;;
  *)
    echo "unknown verb" >&2
    exit 1
    ;;
esac
