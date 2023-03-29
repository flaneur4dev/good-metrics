package storage

import (
	"fmt"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
)

type MemStorage map[string]cs.Metrics

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
