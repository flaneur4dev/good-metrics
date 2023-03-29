package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flaneur4dev/good-metrics/internal/api"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/metrics"
)

var (
	rawPollInterval, _   = utils.EnvVar("POLL_INTERVAL", "2sec").(string)
	rawReportInterval, _ = utils.EnvVar("REPORT_INTERVAL", "10sec").(string)
)

func main() {
	pollInterval, err := strconv.Atoi(strings.TrimRight(rawPollInterval, "seconds"))
	if err != nil {
		fmt.Println("Incorrect parameter!")
		os.Exit(1)
	}

	reportInterval, err := strconv.Atoi(strings.TrimRight(rawReportInterval, "seconds"))
	if err != nil {
		fmt.Println("Incorrect parameter!")
		os.Exit(1)
	}

	start := time.Now()
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		metrics.Update()

		if int(t.Sub(start).Seconds())%reportInterval == 0 {
			for k, v := range metrics.List {
				metrics.CMetrics.AddValue(k, v())

				b, err := json.Marshal(metrics.CMetrics)
				if err != nil {
					fmt.Println(err)
					continue
				}

				api.Fetch(http.MethodPost, "update/", bytes.NewReader(b))
			}

			metrics.CMetrics.AddDelta("PollCount", metrics.PollCount)

			b, err := json.Marshal(metrics.CMetrics)
			if err != nil {
				fmt.Println(err)
				continue
			}

			api.Fetch(http.MethodPost, "update/", bytes.NewReader(b))
		}
	}
}

// func main() {
// 	start := time.Now()
// 	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
// 	defer ticker.Stop()

// 	for t := range ticker.C {
// 		metrics.Update()

// 		if int(t.Sub(start).Seconds())%reportInterval == 0 {
// 			for k, v := range metrics.List {
// 				api.Fetch(http.MethodPost, fmt.Sprintf("update/gauge/%s/%.3f", k, v()), nil)
// 			}

// 			api.Fetch(http.MethodPost, fmt.Sprintf("update/counter/PollCount/%d", metrics.PollCount), nil)
// 		}
// 	}
// }
