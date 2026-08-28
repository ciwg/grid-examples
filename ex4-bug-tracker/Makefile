# Intent: Present Ex4's server, seeded demonstration, and full verification as distinct, documented entry points. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help serve demo verify

help:
	@printf '%s\n' 'Ex4 Bug Tracker:' '  serve   Start the application server on port 7035.' '  demo    Start the seeded demonstration server.' '  verify  Run Go, error, and browser-signing verification.'

serve:
	go run ./cmd/bug-tracker

demo:
	bash scripts/run-demo.sh

verify:
	bash scripts/verify.sh
