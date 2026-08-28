# Intent: Keep Ex6's disposable walkthroughs and opt-in integration proofs explicit rather than treating them as one server startup. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help quickstart writer-example check docker-integration relay-integration

help:
	@printf '%s\n' 'Ex6 Operational Knowledge Agent Runtime:' '  quickstart          Run the built-in walkthrough in a temporary workspace.' '  writer-example      Run the installed writer-package walkthrough.' '  check               Run Go verification.' '  docker-integration  Run the opt-in Docker worker proof.' '  relay-integration   Run the opt-in two-runtime relay proof.'

quickstart:
	./examples/builtin-quickstart.sh

writer-example:
	./examples/install-writer-agent.sh

check:
	go test ./...
	go vet ./...
	errcheck ./...

docker-integration:
	MOKS_DOCKER_INTEGRATION=1 go test ./cmd/moks -run '^TestProcedureExecutionAdapterDockerEndToEnd$$' -count=1

relay-integration:
	MOKS_DEPLOYED_RELAY_INTEGRATION=1 go test ./cmd/moks -run '^TestDeployedTwoNodeWorkflowRelay$$' -count=1
