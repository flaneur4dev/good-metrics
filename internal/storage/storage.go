package storage

import (
	"errors"
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

func (ms *MemStorage) Add(t, n, v string) error {
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
