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

var address, rawPollInterval, rawReportInterval, key string

func main() {
	flag.Parse()

	address = utils.StringEnv("ADDRESS", address)
	rawPollInterval = utils.StringEnv("POLL_INTERVAL", rawPollInterval)
	rawReportInterval = utils.StringEnv("REPORT_INTERVAL", rawReportInterval)
	key = utils.StringEnv("KEY", key)

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

				if key != "" {
					msg := fmt.Sprintf(utils.GaugeTmpl, k, v())
					metrics.CMetrics.Hash = utils.Sign256(msg, key)
				}

				b, err := json.Marshal(metrics.CMetrics)
				if err != nil {
					fmt.Println(err)
					continue
				}

				mc.Fetch(http.MethodPost, "/update/", bytes.NewReader(b))
			}

			metrics.CMetrics.AddDelta("PollCount", metrics.PollCount)

			if key != "" {
				msg := fmt.Sprintf(utils.CounterTmpl, "PollCount", metrics.PollCount)
				metrics.CMetrics.Hash = utils.Sign256(msg, key)
			}

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
	flag.StringVar(&key, "k", "", "secret key")
}
