.PHONY: run test lint docker-up docker-down build

run:
	REDIS_ADDR=localhost:6380 go run ./cmd/server

test:
	go test ./... -race -count=1

lint:
	golangci-lint run ./...

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/server

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down
