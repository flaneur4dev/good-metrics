package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

func main() {
	ms := storage.New()
	r := chi.NewRouter()

	r.Get("/", handlers.HandleMetrics(ms))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(ms))
	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(ms))

	log.Fatal(http.ListenAndServe(":8080", r))
}
