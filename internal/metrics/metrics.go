package metrics

import (
	"math/rand"
	"runtime"
)

type (
	gauge   float64
	counter int64
)

var (
	ms          = &runtime.MemStats{}
	PollCount   counter
	RandomValue = rand.Float64()
	List        = map[string]func() gauge{
		"gauge/Alloc":         func() gauge { return gauge(ms.Alloc) },
		"gauge/BuckHashSys":   func() gauge { return gauge(ms.BuckHashSys) },
		"gauge/Frees":         func() gauge { return gauge(ms.Frees) },
		"gauge/GCCPUFraction": func() gauge { return gauge(ms.GCCPUFraction) },
		"gauge/GCSys":         func() gauge { return gauge(ms.GCSys) },
		"gauge/HeapAlloc":     func() gauge { return gauge(ms.HeapAlloc) },
		"gauge/HeapIdle":      func() gauge { return gauge(ms.HeapIdle) },
		"gauge/HeapInuse":     func() gauge { return gauge(ms.HeapInuse) },
		"gauge/HeapObjects":   func() gauge { return gauge(ms.HeapObjects) },
		"gauge/HeapReleased":  func() gauge { return gauge(ms.HeapReleased) },
		"gauge/HeapSys":       func() gauge { return gauge(ms.HeapSys) },
		"gauge/LastGC":        func() gauge { return gauge(ms.LastGC) },
		"gauge/Lookups":       func() gauge { return gauge(ms.Lookups) },
		"gauge/MCacheInuse":   func() gauge { return gauge(ms.MCacheInuse) },
		"gauge/MCacheSys":     func() gauge { return gauge(ms.MCacheSys) },
		"gauge/MSpanInuse":    func() gauge { return gauge(ms.MSpanInuse) },
		"gauge/MSpanSys":      func() gauge { return gauge(ms.MSpanSys) },
		"gauge/Mallocs":       func() gauge { return gauge(ms.Mallocs) },
		"gauge/NextGC":        func() gauge { return gauge(ms.NextGC) },
		"gauge/NumForcedGC":   func() gauge { return gauge(ms.NumForcedGC) },
		"gauge/NumGC":         func() gauge { return gauge(ms.NumGC) },
		"gauge/OtherSys":      func() gauge { return gauge(ms.OtherSys) },
		"gauge/PauseTotalNs":  func() gauge { return gauge(ms.PauseTotalNs) },
		"gauge/StackInuse":    func() gauge { return gauge(ms.StackInuse) },
		"gauge/StackSys":      func() gauge { return gauge(ms.StackSys) },
		"gauge/Sys":           func() gauge { return gauge(ms.Sys) },
		"gauge/TotalAlloc":    func() gauge { return gauge(ms.TotalAlloc) },
	}
)

func Update() {
	PollCount++
	RandomValue = rand.Float64()
	runtime.ReadMemStats(ms)
}
