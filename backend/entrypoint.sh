#!/bin/sh
set -e

echo "Running migrations"
migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

echo "Starting server"
go run ./cmd
