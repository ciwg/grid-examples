#!/bin/sh
set -eu

case "${1-}" in
  describe)
    printf '%s\n' '{"id":"procedure-execution-adapter","version":"0.1.0","description":"Docker-confined adapter for the procedure-execution workflow","workflow_adapters":[{"name":"procedure-execution","image":"sha256:e02dfedf0daa5770bb785d11b4c4c8f51e377ad8144d73bf0846d5e64fb9410d","input_pcid":"bafkreiawxq2i7q57tks6f5viofxkko2jf2txmlurbp3i33svynytyjswfq","output_pcid":"bafkreiamprv3apzowjzqbkp3hnrhrla5aq7lp5kyzbca5j3iv4v5jmhwa4","cpus":"0.5","memory":"128m","pids_limit":64,"timeout":"30s"}]}'
    ;;
  *)
    printf 'usage: %s describe\n' "$0" >&2
    exit 64
    ;;
esac
