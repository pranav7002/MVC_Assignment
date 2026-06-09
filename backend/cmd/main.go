package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
)

func main() {

    r := chi.NewRouter()

    log.Println("running on :8080")

    http.ListenAndServe(":8080", r)
}