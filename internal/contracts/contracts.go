package contracts

type (
	Gauge   float64
	Counter int64

	Metrics struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *Counter `json:"delta,omitempty"`
		Value *Gauge   `json:"value,omitempty"`
	}
)
