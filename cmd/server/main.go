package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/middlewares/compression"
	"github.com/flaneur4dev/good-metrics/internal/storage/memory"
	"github.com/flaneur4dev/good-metrics/internal/storage/pgdb"
)

type Storage interface {
	AllMetrics() ([]string, []string)
	OneMetric(t, n string) (cs.Metrics, error)
	Update(n string, nm cs.Metrics) (cs.Metrics, error)
	Check() error
	Close() error
}

var (
	address, storeFile, rawStoreInterval, key, dsn string
	restore                                        bool
)

func main() {
	flag.Parse()

	address = utils.StringEnv("ADDRESS", address)
	storeFile = utils.StringEnv("STORE_FILE", storeFile)
	rawStoreInterval = utils.StringEnv("STORE_INTERVAL", rawStoreInterval)
	key = utils.StringEnv("KEY", key)
	dsn = utils.StringEnv("DATABASE_DSN", dsn)
	restore = utils.BoolEnv("RESTORE", restore)

	storeInterval, err := time.ParseDuration(rawStoreInterval)
	if err != nil {
		log.Fatal("incorrect parameter: ", rawStoreInterval)
	}

	var s Storage
	if dsn != "" && storeFile == "/tmp/devops-metrics-db.json" {
		fmt.Println("db")
		s, err = pgdb.New(dsn, key)
		if err != nil {
			log.Fatal("can't connect to storage: ", err)
		}
	} else {
		fmt.Println("memory")
		s = memory.New(storeFile, key, storeInterval.Seconds(), restore)
	}
	defer s.Close()

	r := chi.NewRouter()
	r.Use(compression.HandleGzip)

	r.Get("/", handlers.HandleMetrics(s))
	r.Get("/ping", handlers.HandleStorageCheck(s))
	r.Get("/value/{mType}/{mName}", handlers.HandleMetric(s))

	r.Post("/update/{mType}/{mName}/{mValue}", handlers.HandleUpdate(s))
	r.Post("/value/", handlers.HandleMetricJSON(s))
	r.Post("/update/", handlers.HandleUpdateJSON(s))

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
	flag.StringVar(&dsn, "d", "", "db address")
	flag.BoolVar(&restore, "r", true, "restore on start")
}
