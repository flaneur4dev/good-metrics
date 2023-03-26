package metrics

import (
	"math/rand"
	"runtime"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
)

type ClientMetrics cs.Metrics

func (cm *ClientMetrics) AddValue(n string, v cs.Gauge) {
	cm.ID = n
	cm.MType = "gauge"
	cm.Delta = nil
	cm.Value = &v
}

func (cm *ClientMetrics) AddDelta(n string, v cs.Counter) {
	cm.ID = n
	cm.MType = "counter"
	cm.Delta = &v
	cm.Value = nil
}

var (
	ms          = &runtime.MemStats{}
	CMetrics    = &ClientMetrics{}
	PollCount   cs.Counter
	randomValue = rand.Float64()
	List        = map[string]func() cs.Gauge{
		"Alloc":         func() cs.Gauge { return cs.Gauge(ms.Alloc) },
		"BuckHashSys":   func() cs.Gauge { return cs.Gauge(ms.BuckHashSys) },
		"Frees":         func() cs.Gauge { return cs.Gauge(ms.Frees) },
		"GCCPUFraction": func() cs.Gauge { return cs.Gauge(ms.GCCPUFraction) },
		"GCSys":         func() cs.Gauge { return cs.Gauge(ms.GCSys) },
		"HeapAlloc":     func() cs.Gauge { return cs.Gauge(ms.HeapAlloc) },
		"HeapIdle":      func() cs.Gauge { return cs.Gauge(ms.HeapIdle) },
		"HeapInuse":     func() cs.Gauge { return cs.Gauge(ms.HeapInuse) },
		"HeapObjects":   func() cs.Gauge { return cs.Gauge(ms.HeapObjects) },
		"HeapReleased":  func() cs.Gauge { return cs.Gauge(ms.HeapReleased) },
		"HeapSys":       func() cs.Gauge { return cs.Gauge(ms.HeapSys) },
		"LastGC":        func() cs.Gauge { return cs.Gauge(ms.LastGC) },
		"Lookups":       func() cs.Gauge { return cs.Gauge(ms.Lookups) },
		"MCacheInuse":   func() cs.Gauge { return cs.Gauge(ms.MCacheInuse) },
		"MCacheSys":     func() cs.Gauge { return cs.Gauge(ms.MCacheSys) },
		"MSpanInuse":    func() cs.Gauge { return cs.Gauge(ms.MSpanInuse) },
		"MSpanSys":      func() cs.Gauge { return cs.Gauge(ms.MSpanSys) },
		"Mallocs":       func() cs.Gauge { return cs.Gauge(ms.Mallocs) },
		"NextGC":        func() cs.Gauge { return cs.Gauge(ms.NextGC) },
		"NumForcedGC":   func() cs.Gauge { return cs.Gauge(ms.NumForcedGC) },
		"NumGC":         func() cs.Gauge { return cs.Gauge(ms.NumGC) },
		"OtherSys":      func() cs.Gauge { return cs.Gauge(ms.OtherSys) },
		"PauseTotalNs":  func() cs.Gauge { return cs.Gauge(ms.PauseTotalNs) },
		"StackInuse":    func() cs.Gauge { return cs.Gauge(ms.StackInuse) },
		"StackSys":      func() cs.Gauge { return cs.Gauge(ms.StackSys) },
		"Sys":           func() cs.Gauge { return cs.Gauge(ms.Sys) },
		"TotalAlloc":    func() cs.Gauge { return cs.Gauge(ms.TotalAlloc) },
		"RandomValue":   func() cs.Gauge { return cs.Gauge(randomValue) },
	}
)

func Update() {
	PollCount++
	randomValue = rand.Float64()
	runtime.ReadMemStats(ms)
}
