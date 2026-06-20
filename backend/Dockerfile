FROM golang:1.26-alpine
WORKDIR /app

RUN apk add --no-cache curl postgresql-client \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-arm64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
