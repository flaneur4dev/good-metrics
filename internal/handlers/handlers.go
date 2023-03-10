package handlers

import (
	"net/http"
	"strings"
)

type Adder interface {
	Add(t, n, v string) error
}

func New(rep Adder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		// if ct := r.Header.Get("Content-Type"); ct != "text/plain" {
		// 	http.Error(w, "Bad request", http.StatusBadRequest)
		// 	return
		// }

		// "path = /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>"
		pathSlice := strings.Split(r.URL.Path, "/")
		if len(pathSlice) != 5 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err := rep.Add(pathSlice[2], pathSlice[3], pathSlice[4])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
