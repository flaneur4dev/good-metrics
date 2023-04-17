package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
)

type MemStorage struct {
	metrics       map[string]cs.Metrics
	filePath      string
	key           string
	isRestored    bool
	storeInterval time.Duration
	mTimer        *time.Timer
}

func New(fp, k string, siv float64, re bool) *MemStorage {
	ms := &MemStorage{
		metrics:       map[string]cs.Metrics{},
		filePath:      fp,
		key:           k,
		isRestored:    re,
		storeInterval: time.Duration(siv),
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
		if v.MType == utils.GaugeName {
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

	if !(nm.MType == utils.GaugeName && nm.Value != nil) && !(nm.MType == utils.CounterName && nm.Delta != nil) {
		return cs.Metrics{}, e.ErrUnkownMetricType
	}

	if ms.key != "" {
		var msg string
		switch nm.MType {
		case utils.GaugeName:
			msg = fmt.Sprintf(utils.GaugeTmpl, nm.ID, *nm.Value)
		case utils.CounterName:
			msg = fmt.Sprintf(utils.CounterTmpl, nm.ID, *nm.Delta)
		}

		if !utils.IsEqualSign256(msg, nm.Hash, ms.key) {
			return cs.Metrics{}, e.ErrCompromisedData
		}
	}

	if nm.MType == utils.CounterName {
		if m, ok := ms.metrics[n]; ok {
			nv := *m.Delta + *nm.Delta
			nm.Delta = &nv

			msg := fmt.Sprintf(utils.CounterTmpl, nm.ID, *nm.Delta)
			nm.Hash = utils.Sign256(msg, ms.key)
		}
	}

	ms.metrics[n] = nm
	if ms.storeInterval == 0 {
		ms.toFile()
	}

	return nm, nil
}

func (ms *MemStorage) Check() error {
	return e.ErrNoUsedDB
}

func (ms *MemStorage) Close() error {
	if ms.mTimer != nil {
		ms.mTimer.Stop()
	}
	ms.toFile()
	return nil
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
}

func (ms *MemStorage) fromFile() {
	if ms.filePath == "" {
		return
	}

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
}
