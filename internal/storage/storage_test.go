package storage

import (
	"testing"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
)

func TestAdd(t *testing.T) {
	ms := New("", 0, false)
	delta := cs.Counter(42)
	value := cs.Gauge(42.420)

	tests := []struct {
		name string
		m    cs.Metrics
		want error
	}{
		{
			name: "#1",
			m:    cs.Metrics{ID: "metric1", MType: "counter", Delta: &delta},
			want: nil,
		},
		{
			name: "#2",
			m:    cs.Metrics{ID: "metric2", MType: "gauge", Value: &value},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ms.Update(tt.m.ID, tt.m); err != tt.want {
				t.Errorf("Add() = %v, want: %v", err, tt.want)
			}
		})
	}
}
