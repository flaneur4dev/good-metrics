package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

var (
	re     bool
	ad, sf string
	iv     = "0sec"
)

func main() {
	flag.Parse()

	addr, _ := utils.EnvVar("ADDRESS", ad).(string)
	storeFile, _ := utils.EnvVar("STORE_FILE", sf).(string)
	rawStoreInterval, _ := utils.EnvVar("STORE_INTERVAL", iv).(string)
	restore, _ := utils.EnvVar("RESTORE", re).(bool)

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

func init() {
	flag.StringVar(&ad, "a", "localhost:8080", "server address")
	flag.StringVar(&sf, "f", "/tmp/devops-metrics-db.json", "store file")
	flag.BoolVar(&re, "r", true, "restore on start")

	flag.Func("i", "store interval", func(fl string) error {
		iv = fl + "sec"
		return nil
	})
}
