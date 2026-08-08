.PHONY: run build test test-integration tidy

run:
	go run ./cmd/detector

build:
	go build -o bin/detector ./cmd/detector

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

tidy:
	go mod tidy
