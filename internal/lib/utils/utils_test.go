package utils

import "testing"

func TestStringEnv(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "#1",
			value: "test 1",
			want:  "test 1",
		},
		{
			name:  "#2",
			value: "test 2",
			want:  "test 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := StringEnv("SOME_ENV", tt.value); res != tt.want {
				t.Errorf("EnvVar() = %v, want: %v", res, tt.want)
			}
		})
	}
}

func TestBoolEnv(t *testing.T) {
	tests := []struct {
		name  string
		value bool
		want  bool
	}{
		{
			name:  "#1",
			value: true,
			want:  true,
		},
		{
			name:  "#2",
			value: false,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := BoolEnv("SOME_ENV", tt.value); res != tt.want {
				t.Errorf("EnvVar() = %v, want: %v", res, tt.want)
			}
		})
	}
}
