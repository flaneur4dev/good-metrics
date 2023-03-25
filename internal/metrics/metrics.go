package metrics

import (
	"math/rand"
	"runtime"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
)

var (
	ms          = &runtime.MemStats{}
	PollCount   cs.Counter
	RandomValue = rand.Float64()
	List        = map[string]func() cs.Gauge{
		"gauge/Alloc":         func() cs.Gauge { return cs.Gauge(ms.Alloc) },
		"gauge/BuckHashSys":   func() cs.Gauge { return cs.Gauge(ms.BuckHashSys) },
		"gauge/Frees":         func() cs.Gauge { return cs.Gauge(ms.Frees) },
		"gauge/GCCPUFraction": func() cs.Gauge { return cs.Gauge(ms.GCCPUFraction) },
		"gauge/GCSys":         func() cs.Gauge { return cs.Gauge(ms.GCSys) },
		"gauge/HeapAlloc":     func() cs.Gauge { return cs.Gauge(ms.HeapAlloc) },
		"gauge/HeapIdle":      func() cs.Gauge { return cs.Gauge(ms.HeapIdle) },
		"gauge/HeapInuse":     func() cs.Gauge { return cs.Gauge(ms.HeapInuse) },
		"gauge/HeapObjects":   func() cs.Gauge { return cs.Gauge(ms.HeapObjects) },
		"gauge/HeapReleased":  func() cs.Gauge { return cs.Gauge(ms.HeapReleased) },
		"gauge/HeapSys":       func() cs.Gauge { return cs.Gauge(ms.HeapSys) },
		"gauge/LastGC":        func() cs.Gauge { return cs.Gauge(ms.LastGC) },
		"gauge/Lookups":       func() cs.Gauge { return cs.Gauge(ms.Lookups) },
		"gauge/MCacheInuse":   func() cs.Gauge { return cs.Gauge(ms.MCacheInuse) },
		"gauge/MCacheSys":     func() cs.Gauge { return cs.Gauge(ms.MCacheSys) },
		"gauge/MSpanInuse":    func() cs.Gauge { return cs.Gauge(ms.MSpanInuse) },
		"gauge/MSpanSys":      func() cs.Gauge { return cs.Gauge(ms.MSpanSys) },
		"gauge/Mallocs":       func() cs.Gauge { return cs.Gauge(ms.Mallocs) },
		"gauge/NextGC":        func() cs.Gauge { return cs.Gauge(ms.NextGC) },
		"gauge/NumForcedGC":   func() cs.Gauge { return cs.Gauge(ms.NumForcedGC) },
		"gauge/NumGC":         func() cs.Gauge { return cs.Gauge(ms.NumGC) },
		"gauge/OtherSys":      func() cs.Gauge { return cs.Gauge(ms.OtherSys) },
		"gauge/PauseTotalNs":  func() cs.Gauge { return cs.Gauge(ms.PauseTotalNs) },
		"gauge/StackInuse":    func() cs.Gauge { return cs.Gauge(ms.StackInuse) },
		"gauge/StackSys":      func() cs.Gauge { return cs.Gauge(ms.StackSys) },
		"gauge/Sys":           func() cs.Gauge { return cs.Gauge(ms.Sys) },
		"gauge/TotalAlloc":    func() cs.Gauge { return cs.Gauge(ms.TotalAlloc) },
	}
)

func Update() {
	PollCount++
	RandomValue = rand.Float64()
	runtime.ReadMemStats(ms)
}
