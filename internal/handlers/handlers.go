package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
	e "github.com/flaneur4dev/good-metrics/internal/lib/mistakes"
)

type (
	Metrics interface {
		AllMetrics() ([]string, []string)
	}
	Metric interface {
		OneMetric(t, n string) (string, error)
	}
	Updater interface {
		Update(t, n, v string) (string, error)
	}
)

func setMetric(res *cs.Metrics, v string) {
	if res.MType == "gauge" {
		f, _ := strconv.ParseFloat(v, 64)
		value := cs.Gauge(f)
		res.Value = &value
		res.Delta = nil
	} else {
		d, _ := strconv.ParseInt(v, 10, 64)
		value := cs.Counter(d)
		res.Value = nil
		res.Delta = &value
	}
}

func HandleUpdate(rep Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")

		_, err := rep.Update(mType, mName, mValue)
		if err != nil {
			sc := http.StatusBadRequest
			if errors.Is(err, e.ErrUnkownMetricType) {
				sc = http.StatusNotImplemented
			}
			http.Error(w, err.Error(), sc)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HandleUpdateJSON(rep Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var res cs.Metrics
		err = json.Unmarshal(reqBody, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var value string
		if res.Value != nil && res.Delta == nil {
			value = fmt.Sprintf("%.3f", *res.Value)
		}
		if res.Delta != nil && res.Value == nil {
			value = fmt.Sprintf("%d", *res.Delta)
		}

		v, err := rep.Update(res.MType, res.ID, value)
		if err != nil {
			sc := http.StatusBadRequest
			if errors.Is(err, e.ErrUnkownMetricType) {
				sc = http.StatusNotImplemented
			}
			http.Error(w, err.Error(), sc)
			return
		}

		setMetric(&res, v)

		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleMetric(rep Metric) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")

		v, err := rep.OneMetric(mType, mName)
		if err != nil {
			sc := http.StatusNotFound
			if errors.Is(err, e.ErrUnkownMetricType) {
				sc = http.StatusNotImplemented
			}
			http.Error(w, err.Error(), sc)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(v))
	}
}

func HandleMetricJSON(rep Metric) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var res cs.Metrics
		err = json.Unmarshal(reqBody, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		v, err := rep.OneMetric(res.MType, res.ID)
		if err != nil {
			sc := http.StatusNotFound
			if errors.Is(err, e.ErrUnkownMetricType) {
				sc = http.StatusNotImplemented
			}
			http.Error(w, err.Error(), sc)
			return
		}

		setMetric(&res, v)

		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func HandleMetrics(rep Metrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gm, cm := rep.AllMetrics()

		t, err := template.New("metrics webpage").Parse(tmpl)
		if err != nil {
			http.Error(w, "Somthing went wrong", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			GMetics []string
			CMetics []string
		}{
			Title:   "Good metrics",
			GMetics: gm,
			CMetics: cm,
		}

		err = t.Execute(w, data)
		if err != nil {
			http.Error(w, "Somthing went wrong", http.StatusInternalServerError)
			return
		}
	}
}

const tmpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body style="margin-left: 50px;">
		<h1>Метрики</h1>
		<div style="display: flex; gap: 100px;">
			<div>
				<h2 style="color: #9a9a9a;">Тип gauge:</h2>
				<ul>
					{{range .GMetics}}<li>{{ . }}</li>{{else}}<p>Нет данных</p>{{end}}
				</ul>
			</div>
			<div>
				<h2 style="color: #9a9a9a;">Тип counter:</h2>
				<ul>
					{{range .CMetics}}<li>{{ . }}</li>{{else}}<p>Нет данных</p>{{end}}
				</ul>			
			</div>
		</div>
	</body>
</html>
`
