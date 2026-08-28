# Intent: Preserve Ex5's separate normal server, newcomer corpus, browser proof, and verification paths. Source: DI-jigap
.DEFAULT_GOAL := help

RUNTIME_ROOT ?= /tmp/ex5-newcomer-runtime

.PHONY: help serve sample newcomer browser-demo check

help:
	@printf '%s\n' 'Ex5 Operational Knowledge System:' '  serve         Start the server with its default runtime.' '  sample        Load the newcomer sample into RUNTIME_ROOT.' '  newcomer      Load the sample and start the server.' '  browser-demo  Run the verified Chrome browser demonstration.' '  check         Run the documented Go checks.' '  RUNTIME_ROOT=/path  Override the newcomer sample runtime.'

serve:
	go run ./cmd/operational-knowledge

sample:
	./scripts/load-sample-data.sh $(RUNTIME_ROOT)

newcomer:
	./scripts/load-sample-data.sh $(RUNTIME_ROOT)
	go run ./cmd/operational-knowledge -data-root $(RUNTIME_ROOT)

browser-demo:
	./scripts/run-demo-proof.sh

check:
	go test ./...
	go vet ./...
	errcheck ./...
