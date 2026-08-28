# Intent: Keep Ex2's single relay, second logical relay, browser build, and container simulation as separate documented entry points. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help relay second-relay web-build docker-demo check

help:
	@printf '%s\n' 'Ex2 Grid Editor:' '  relay         Start the local relay on port 7015.' '  second-relay  Start the peer relay on port 7016.' '  web-build     Rebuild the browser bundle.' '  docker-demo   Start the two-relay Docker simulation.' '  check         Run Go and browser verification.'

relay:
	go run ./cmd/grid-relay --listen 127.0.0.1:7015 --data-root .grid-editor/alice

second-relay:
	go run ./cmd/grid-relay --listen 127.0.0.1:7016 --data-root .grid-editor/bob --peer http://127.0.0.1:7015

web-build:
	cd web && npm run build

docker-demo:
	docker-compose up --build

check:
	go vet ./...
	go test ./...
	errcheck ./...
	cd web && npm test && npm run build
