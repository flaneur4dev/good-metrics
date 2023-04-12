package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/middlewares/compression"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

var (
	address, storeFile, rawStoreInterval, key string
	restore                                   bool
)

func main() {
	flag.Parse()

	address = utils.StringEnv("ADDRESS", address)
	storeFile = utils.StringEnv("STORE_FILE", storeFile)
	rawStoreInterval = utils.StringEnv("STORE_INTERVAL", rawStoreInterval)
	restore = utils.BoolEnv("RESTORE", restore)
	key = utils.StringEnv("KEY", key)

	storeInterval, err := time.ParseDuration(rawStoreInterval)
	if err != nil {
		log.Fatal("Incorrect parameter: ", rawStoreInterval)
	}

	ms := storage.New(storeFile, key, storeInterval.Seconds(), restore)
	defer ms.Close()

	r := chi.NewRouter()
	r.Use(compression.HandleGzip)

	r.Get("/", handlers.HandleMetrics(ms))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(ms))
	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(ms))

	r.Post("/value/", handlers.HandleMetricJSON(ms))
	r.Post("/update/", handlers.HandleUpdateJSON(ms))

	err = http.ListenAndServe(address, r)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&address, "a", "localhost:8080", "server address")
	flag.StringVar(&storeFile, "f", "/tmp/devops-metrics-db.json", "store file")
	flag.StringVar(&rawStoreInterval, "i", "300s", "store interval")
	flag.StringVar(&key, "k", "", "secret key")
	flag.BoolVar(&restore, "r", true, "restore on start")
}
