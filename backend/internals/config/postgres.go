package config

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() (*pgxpool.Pool, error) {

    conn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        Get("DB_USER"),
        Get("DB_PASSWORD"),
        Get("DB_HOST"),
        Get("DB_PORT"),
        Get("DB_NAME"),
    )

    return pgxpool.New(context.Background(), conn)
}