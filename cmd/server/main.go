package main

import (
	"log"
	"net/http"

	"github.com/flaneur4dev/good-metrics/internal/handlers"
	"github.com/flaneur4dev/good-metrics/internal/storage"
)

func main() {
	http.HandleFunc("/update/", handlers.New(storage.New()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
