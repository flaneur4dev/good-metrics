package main

import (
	"log"
	"net/http"

	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

func main() {
	fmt.Println(os.Environ())

	addr, _ := utils.EnvVar("ADDRESS", "localhost:8080").(string)
	storeFile, _ := utils.EnvVar("STORE_FILE", "/tmp/devops-metrics-db.json").(string)
	rawStoreInterval, _ := utils.EnvVar("STORE_INTERVAL", "300sec").(string)
	restore, _ := utils.EnvVar("RESTORE", true).(bool)

	storeInterval, err := strconv.Atoi(strings.TrimRight(rawStoreInterval, "sec"))
	if err != nil {
		log.Fatal("Incorrect parameter!")
	}

	r := chi.NewRouter()
	ms := storage.New(storeFile, storeInterval, restore)
	defer ms.Close()

	r.Get("/", handlers.HandleMetrics(ms))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(ms))
	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(ms))

	r.Post("/value/", handlers.HandleMetricJSON(ms))
	r.Post("/update/", handlers.HandleUpdateJSON(ms))

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
