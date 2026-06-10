.PHONY: run build test migrate-up migrate-down seed

run:
	go run ./cmd

build:
	go build -o ./bin/server ./cmd

test:
	go test ./internals/...

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5433/coc?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5433/coc?sslmode=disable" down
