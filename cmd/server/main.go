package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

func main() {
	addr, _ := utils.EnvVar("ADDRESS", ":8080").(string)
	ms := storage.New()
	r := chi.NewRouter()

	r.Get("/", handlers.HandleMetrics(ms))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(ms))
	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(ms))

	r.Post("/value/", handlers.HandleMetricJSON(ms))
	r.Post("/update/", handlers.HandleUpdateJSON(ms))

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
