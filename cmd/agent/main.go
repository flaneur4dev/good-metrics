package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/flaneur4dev/good-metrics/internal/api"
	"github.com/flaneur4dev/good-metrics/internal/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 //time.Second
)

func main() {
	start := time.Now()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		metrics.Update()

		if int(t.Sub(start).Seconds())%reportInterval == 0 {
			for k, v := range metrics.List {
				api.Fetch(http.MethodPost, fmt.Sprintf("update/%s/%.3f", k, v()), nil)
			}

			api.Fetch(http.MethodPost, fmt.Sprintf("update/counter/PollCount/%d", metrics.PollCount), nil)
			api.Fetch(http.MethodPost, fmt.Sprintf("update/gauge/RandomValue/%.3f", metrics.RandomValue), nil)
		}
	}
}
