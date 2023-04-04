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
	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
)

type (
	Metrics interface {
		AllMetrics() ([]string, []string)
	}
	Metric interface {
		OneMetric(t, n string) (cs.Metrics, error)
	}
	Updater interface {
		Update(n string, nm cs.Metrics) (cs.Metrics, error)
	}
)

func HandleUpdate(rep Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mType := chi.URLParam(r, "mType")
		mName := chi.URLParam(r, "mName")
		mValue := chi.URLParam(r, "mValue")

		res := cs.Metrics{ID: mName, MType: mType}

		switch mType {
		case utils.GaugeName:
			f, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			value := cs.Gauge(f)
			res.Value = &value
		case utils.CounterName:
			d, err := strconv.ParseInt(mValue, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			value := cs.Counter(d)
			res.Delta = &value
		}

		_, err := rep.Update(res.ID, res)
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

		v, err := rep.Update(res.ID, res)
		if err != nil {
			sc := http.StatusBadRequest
			if errors.Is(err, e.ErrUnkownMetricType) {
				sc = http.StatusNotImplemented
			}
			http.Error(w, err.Error(), sc)
			return
		}

		b, err := json.Marshal(v)
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

		var rv string
		switch mType {
		case utils.GaugeName:
			rv = fmt.Sprintf("%.3f", *v.Value)
		case utils.CounterName:
			rv = fmt.Sprintf("%d", *v.Delta)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rv))
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

		b, err := json.Marshal(v)
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

		w.Header().Set("Content-Type", "text/html")

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
