test:
	@go test ./... -v

benchmark:
	@go test -bench=. -v -count=10 -run=^# ./...

run/server:
	@go run ./cmd/server

build/server:
	@go build -o ./bin/server -v -race ./cmd/server

run/cli:
	@go run ./cmd/cli

build/cli:
	@go build -o ./bin/cli -v -race ./cmd/cli

build/docker:
	@docker build --tag cacher:dev .

run/docker:
	@echo Press Ctrl+c to stop and remove the container
	@docker run --rm --name cacher -p 8595:8595 cacher:dev

up/docker:
	@docker run -d --name cacher -p 8595:8595 cacher:dev
	@echo The server is running in background
