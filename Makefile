.PHONY: run build test tidy

run:
	go run ./cmd/detector

build:
	go build -o bin/detector ./cmd/detector

test:
	go test ./...

tidy:
	go mod tidy
