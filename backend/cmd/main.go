package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/pranav7002/MVC_Assignment/internal/database"
)

func main() {
	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Error("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", 
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	cfg := config{
		addr: fmt.Sprintf(":%s", port),
		db: dbConfig{
			dsn: dbDSN,
		},
	}

	api := application{
		config: cfg,
		dbpool: nil,
	}

	pool, err := api.connectDB()
	if err != nil {
		slog.Error("Unable to create connection pool", "error", err)
		os.Exit(1)
	}
	api.dbpool = pool
	defer pool.Close() 

	database.Seed(pool)

	if err := api.run(api.mount()); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}