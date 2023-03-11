package storage

import (
	"errors"
	"fmt"
	"strconv"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
)

type MemStorage struct {
	gauge   map[string]cs.Gauge
	counter map[string]cs.Counter
}

func New() *MemStorage {
	return &MemStorage{
		gauge:   map[string]cs.Gauge{},
		counter: map[string]cs.Counter{},
	}
}

func (ms *MemStorage) AllMetrics() ([]string, []string) {
	gm := make([]string, len(ms.gauge))
	cm := make([]string, len(ms.counter))

	i := 0
	for k, v := range ms.gauge {
		gm[i] = fmt.Sprintf("%s: %.3f", k, v)
		i++
	}

	i = 0
	for k, v := range ms.counter {
		cm[i] = fmt.Sprintf("%s: %d", k, v)
		i++
	}

	return gm, cm
}

func (ms *MemStorage) OneMetric(t, n string) (string, error) {
	switch t {
	case "gauge":
		v, ok := ms.gauge[n]
		if !ok {
			return "", errors.New("no such metric")
		}
		return fmt.Sprintf("%.3f", v), nil
	case "counter":
		v, ok := ms.counter[n]
		if !ok {
			return "", errors.New("no such metric")
		}
		return fmt.Sprintf("%d", v), nil
	default:
		return "", errors.New("unknown metric type")
	}
}

func (ms *MemStorage) Update(t, n, v string) error {
	switch t {
	case "gauge":
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return errors.New("invalid data")
		}
		ms.gauge[n] = cs.Gauge(f)
	case "counter":
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return errors.New("invalid data")
		}
		ms.counter[n] += cs.Counter(i)
	default:
		return errors.New("unknown metric type")
	}

	return nil
}
