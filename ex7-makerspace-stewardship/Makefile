# Intent: Distinguish Ex7's local record server from its ingress and browser evidence proofs. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help serve ingress-proof browser-proof check

help:
	@printf '%s\n' 'Ex7 Makerspace Stewardship:' '  serve          Start the local record server on port 7037.' '  ingress-proof  Run the signed-record ingress proof.' '  browser-proof  Run the two-agent Chrome proof.' '  check          Run Go verification and error checks.'

serve:
	go run ./cmd/makerspace-stewardship

ingress-proof:
	bash scripts/run-record-ingress-proof.sh

browser-proof:
	bash scripts/run-two-agent-browser-proof.sh

check:
	go test ./...
	go vet ./...
	errcheck ./...
