package storage

import (
	"testing"

	cs "github.com/flaneur4dev/good-metrics/internal/contracts"
)

func TestAdd(t *testing.T) {
	ms := &MemStorage{
		gauge:   map[string]cs.Gauge{},
		counter: map[string]cs.Counter{},
	}

	tests := []struct {
		name string
		t    string
		n    string
		v    string
		want error
	}{
		{
			name: "#1",
			t:    "counter",
			n:    "metric1",
			v:    "42",
			want: nil,
		},
		{
			name: "#2",
			t:    "gauge",
			n:    "metric2",
			v:    "42.42",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ms.Update(tt.t, tt.n, tt.v); err != tt.want {
				t.Errorf("Add() = %v, want: %v", err, tt.want)
			}
		})
	}
}
