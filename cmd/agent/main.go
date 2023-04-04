package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flaneur4dev/good-metrics/internal/client"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
	"github.com/flaneur4dev/good-metrics/internal/metrics"
)

var address, rawPollInterval, rawReportInterval string

func main() {
	flag.Parse()

	address, _ := utils.EnvVar("ADDRESS", address).(string)
	rawPollInterval, _ := utils.EnvVar("POLL_INTERVAL", rawPollInterval).(string)
	rawReportInterval, _ := utils.EnvVar("REPORT_INTERVAL", rawReportInterval).(string)

	pollInterval, err := time.ParseDuration(rawPollInterval)
	if err != nil {
		log.Fatal("Incorrect parameter: ", rawPollInterval)
	}

	reportInterval, err := time.ParseDuration(rawReportInterval)
	if err != nil {
		log.Fatal("Incorrect parameter: ", rawReportInterval)
	}

	mc := client.New(address, false)

	start := time.Now()
	ticker := time.NewTicker(time.Duration(pollInterval.Seconds()) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		metrics.Update()

		if int(t.Sub(start).Seconds())%int(reportInterval.Seconds()) == 0 {
			for k, v := range metrics.List {
				metrics.CMetrics.AddValue(k, v())

				b, err := json.Marshal(metrics.CMetrics)
				if err != nil {
					fmt.Println(err)
					continue
				}

				mc.Fetch(http.MethodPost, "/update/", bytes.NewReader(b))
			}

			metrics.CMetrics.AddDelta("PollCount", metrics.PollCount)

			b, err := json.Marshal(metrics.CMetrics)
			if err != nil {
				fmt.Println(err)
				continue
			}

			mc.Fetch(http.MethodPost, "/update/", bytes.NewReader(b))
		}
	}
}

func init() {
	flag.StringVar(&address, "a", "localhost:8080", "server address")
	flag.StringVar(&rawPollInterval, "p", "2s", "poll interval")
	flag.StringVar(&rawReportInterval, "r", "10s", "report interval")
}
