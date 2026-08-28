# Intent: Surface Ex1's documented container scenarios and local verification without hiding their distinct topology. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help demo demo-refusal demo-timeout check

help:
	@printf '%s\n' 'Ex1 order flow:' '  demo          Run the Docker happy-path topology.' '  demo-refusal  Run the warehouse-refusal fixture.' '  demo-timeout  Run the carrier-timeout fixture.' '  check         Run Go tests and error checking.'

demo:
	bash docker/run-demo.sh

demo-refusal:
	bash docker/run-demo.sh warehouse-refusal.json

demo-timeout:
	bash docker/run-demo.sh carrier-timeout.json

check:
	go test ./...
	errcheck ./...
