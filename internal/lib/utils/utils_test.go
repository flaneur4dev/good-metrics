package utils

import (
	"testing"
)

func TestEnvVar(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{
			name:  "#1",
			value: "test 42",
			want:  "test 42",
		},
		{
			name:  "#2",
			value: 42,
			want:  42,
		},
		{
			name:  "#3",
			value: false,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := EnvVar("ENV", tt.value); res != tt.want {
				t.Errorf("EnvVar() = %v, want: %v", res, tt.want)
			}
		})
	}
}
