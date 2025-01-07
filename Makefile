test:
	@go test ./... -v

run/server:
	@go run ./cmd/server

build/server:
	@go build -o bin/server -v ./cmd/server
