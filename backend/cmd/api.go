package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)              // rate limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // rate limiting and analytics
	r.Use(middleware.Logger) 
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr: app.config.addr,
		Handler: h,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) connectDB() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), app.config.db.dsn)
	if err != nil {
		return nil, err
	}

	log.Printf("db connected successfully")
	return dbpool, nil
}

type application struct {
	config config
	dbpool *pgxpool.Pool
	// logger
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}
