# Intent: Expose Ex3's loopback, remote-capability, container, and embodiment workflows without wrapping its destructive reset script. Source: DI-jigap
.DEFAULT_GOAL := help

.PHONY: help relay remote-relay docker-up docker-down web-build sidecar-build check

help:
	@printf '%s\n' 'Ex3 WebSocket Grid Editor:' '  relay          Start the loopback relay on port 7025.' '  remote-relay   Start the documented remote-access demo relay.' '  docker-up      Start the two-relay Docker setup in the background.' '  docker-down    Stop the Docker setup.' '  web-build      Rebuild the browser bundle.' '  sidecar-build  Build the Neovim sidecar.' '  check          Run Go, browser, and sidecar checks.'

relay:
	go run ./cmd/grid-relay --listen 127.0.0.1:7025 --data-root .grid-editor/alice

remote-relay:
	go run ./cmd/grid-relay --listen 0.0.0.0:7025 --remote-access-token ex3-demo-access

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

web-build:
	cd web && npm run build

sidecar-build:
	cd cmd/grid-nvim-sidecar && npm run build

check:
	go vet ./...
	go test ./...
	errcheck ./...
	cd web && npm test && npm run build
	cd cmd/grid-nvim-sidecar && npm run build
