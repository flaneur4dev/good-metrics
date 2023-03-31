package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/flaneur4dev/good-metrics/internal/api"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/metrics"
)

var (
	ad  string
	piv = "2sec"
	riv = "10sec"
)

func main() {
	flag.Parse()

	rawReportInterval, _ := utils.EnvVar("REPORT_INTERVAL", riv).(string)
	rawPollInterval, _ := utils.EnvVar("POLL_INTERVAL", piv).(string)
	addr, _ := utils.EnvVar("ADDRESS", ad).(string)
	baseURL := "http://" + addr

	pollInterval, err := strconv.Atoi(strings.TrimRight(rawPollInterval, "sec"))
	if err != nil {
		log.Fatal("Incorrect parameter!")
	}

	reportInterval, err := strconv.Atoi(strings.TrimRight(rawReportInterval, "sec"))
	if err != nil {
		log.Fatal("Incorrect parameter!")
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

				api.Fetch(http.MethodPost, baseURL+"/update/", bytes.NewReader(b))
			}

			metrics.CMetrics.AddDelta("PollCount", metrics.PollCount)

			b, err := json.Marshal(metrics.CMetrics)
			if err != nil {
				fmt.Println(err)
				continue
			}

			api.Fetch(http.MethodPost, baseURL+"/update/", bytes.NewReader(b))
		}
	}
}

func init() {
	flag.StringVar(&ad, "a", "localhost:8080", "server address")

	flag.Func("p", "poll interval", func(fl string) error {
		piv = fl + "sec"
		return nil
	})

	flag.Func("r", "report interval", func(fl string) error {
		riv = fl + "sec"
		return nil
	})
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
