package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
)

type MemStorage struct {
	metrics       map[string]cs.Metrics
	storeInterval time.Duration
	filePath      string
	isRestored    bool
	mTimer        *time.Timer
}

func New(fp string, iv int, re bool) *MemStorage {
	ms := &MemStorage{
		metrics:       map[string]cs.Metrics{},
		filePath:      fp,
		storeInterval: time.Duration(iv),
		isRestored:    re,
	}

	if ms.isRestored {
		ms.fromFile()
	}

	if ms.storeInterval > 0 {
		ms.intervalSave()
	}

	return ms
}

func (ms *MemStorage) AllMetrics() (gm, cm []string) {
	for k, v := range ms.metrics {
		if v.MType == "gauge" {
			gm = append(gm, fmt.Sprintf("%s: %f", k, *v.Value))
		} else {
			cm = append(cm, fmt.Sprintf("%s: %d", k, *v.Delta))
		}
	}
	return
}

func (ms *MemStorage) OneMetric(t, n string) (cs.Metrics, error) {
	m, ok := ms.metrics[n]
	if !ok || m.MType != t {
		return cs.Metrics{}, e.ErrNoMetric
	}
	return m, nil
}

func (ms *MemStorage) Update(n string, nm cs.Metrics) (cs.Metrics, error) {
	if n == "" {
		return cs.Metrics{}, e.ErrInvalidData
	}

	if !(nm.MType == "gauge" && nm.Value != nil) && !(nm.MType == "counter" && nm.Delta != nil) {
		return cs.Metrics{}, e.ErrUnkownMetricType
	}

	if nm.MType == "counter" {
		if m, ok := ms.metrics[n]; ok {
			nv := *m.Delta + *nm.Delta
			nm.Delta = &nv
			// ms.metrics[n] = m
			// return m, nil
		}
	}

	ms.metrics[n] = nm
	if ms.storeInterval == 0 {
		ms.toFile()
	}

	return nm, nil
}

func (ms *MemStorage) Close() {
	if ms.mTimer != nil {
		ms.mTimer.Stop()
	}
	ms.toFile()
}

func (ms *MemStorage) intervalSave() {
	ms.toFile()
	ms.mTimer = time.AfterFunc(ms.storeInterval*time.Second, func() {
		ms.intervalSave()
	})
}

func (ms *MemStorage) toFile() {
	if ms.filePath == "" {
		return
	}

	// fd, err := os.OpenFile(ms.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	data, err := json.Marshal(ms.metrics)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(ms.filePath, data, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	// _, err = fd.Write(b)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// if err := fd.Close(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

func (ms *MemStorage) fromFile() {
	if ms.filePath == "" {
		return
	}

	// fd, err := os.Open(ms.filePath)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data, err := os.ReadFile(ms.filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var res map[string]cs.Metrics
	err = json.Unmarshal(data, &res)
	if err != nil {
		fmt.Println(err)
		return
	}

	ms.metrics = res

	// if err := fd.Close(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
