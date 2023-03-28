package storage

import (
	"fmt"
	// "strconv"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
)

// type MemStorage struct {
// 	gauge   map[string]cs.Gauge
// 	counter map[string]cs.Counter
// }

type MemStorage map[string]cs.Metrics

// func New() *MemStorage {
// 	return &MemStorage{
// 		gauge:   map[string]cs.Gauge{},
// 		counter: map[string]cs.Counter{},
// 	}
// }

func New() MemStorage {
	return MemStorage{}
}

func (ms MemStorage) AllMetrics() (gm, cm []string) {
	for k, v := range ms {
		if v.MType == "gauge" {
			gm = append(gm, fmt.Sprintf("%s: %f", k, *v.Value))
		} else {
			cm = append(cm, fmt.Sprintf("%s: %d", k, *v.Delta))
		}
	}
	return
}

func (ms MemStorage) OneMetric(t, n string) (cs.Metrics, error) {
	m, ok := ms[n]
	if !ok || m.MType != t {
		return cs.Metrics{}, e.ErrNoMetric
	}
	return m, nil
}

func (ms MemStorage) Update(n string, nm cs.Metrics) (cs.Metrics, error) {
	if n == "" {
		return cs.Metrics{}, e.ErrInvalidData
	}

	if !(nm.MType == "gauge" && nm.Value != nil) && !(nm.MType == "counter" && nm.Delta != nil) {
		return cs.Metrics{}, e.ErrUnkownMetricType
	}

	if nm.MType == "counter" {
		if m, ok := ms[n]; ok {
			nv := *m.Delta + *nm.Delta
			m.Delta = &nv
			ms[n] = m
			return m, nil
		}
	}

	ms[n] = nm
	return nm, nil
}

// func (ms *MemStorage) AllMetrics() ([]string, []string) {
// 	gm := make([]string, len(ms.gauge))
// 	cm := make([]string, len(ms.counter))

// 	i := 0
// 	for k, v := range ms.gauge {
// 		gm[i] = fmt.Sprintf("%s: %f", k, v)
// 		i++
// 	}

// 	i = 0
// 	for k, v := range ms.counter {
// 		cm[i] = fmt.Sprintf("%s: %d", k, v)
// 		i++
// 	}

// 	return gm, cm
// }

// func (ms *MemStorage) OneMetric(t, n string) (string, error) {
// 	switch t {
// 	case "gauge":
// 		v, ok := ms.gauge[n]
// 		if !ok {
// 			return "", e.ErrNoMetric
// 		}
// 		return fmt.Sprintf("%f", v), nil
// 	case "counter":
// 		v, ok := ms.counter[n]
// 		if !ok {
// 			return "", e.ErrNoMetric
// 		}
// 		return fmt.Sprintf("%d", v), nil
// 	default:
// 		return "", e.ErrUnkownMetricType
// 	}
// }

// func (ms *MemStorage) Update(t, n, v string) (string, error) {
// 	switch t {
// 	case "gauge":
// 		f, err := strconv.ParseFloat(v, 64)
// 		if err != nil {
// 			return "", e.ErrInvalidData
// 		}

// 		ms.gauge[n] = cs.Gauge(f)
// 		return v, nil
// 	case "counter":
// 		i, err := strconv.ParseInt(v, 10, 64)
// 		if err != nil {
// 			return "", e.ErrInvalidData
// 		}

// 		ms.counter[n] += cs.Counter(i)
// 		nv := ms.counter[n]
// 		return fmt.Sprintf("%d", nv), nil
// 	default:
// 		return "", e.ErrUnkownMetricType
// 	}
// }
