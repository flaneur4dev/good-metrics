package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

func main() {
	addr, _ := utils.EnvVar("ADDRESS", "localhost:8080").(string)
	addrSlice := strings.Split(addr, ":")
	if len(addrSlice) != 2 {
		fmt.Println("Incorrect parameter!")
		os.Exit(1)
	}

	port := addrSlice[1]
	ms := storage.New()
	r := chi.NewRouter()

	r.Get("/", handlers.HandleMetrics(ms))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(ms))
	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(ms))

	r.Post("/value/", handlers.HandleMetricJSON(ms))
	r.Post("/update/", handlers.HandleUpdateJSON(ms))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
