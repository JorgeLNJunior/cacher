test:
	@go test ./... -v

run/server:
	@go run ./cmd/server

build/server:
	@go build -o ./bin/server -v ./cmd/server

run/cli:
	@go run ./cmd/cli

build/cli:
	@go build -o ./bin/cli -v ./cmd/cli
