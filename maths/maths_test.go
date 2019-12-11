package maths

import "testing"

func TestAbsInt(t *testing.T) {
	tests := []struct {
		name string
		x    int
		want int
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"-1", -1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsInt(tt.x); got != tt.want {
				t.Errorf("AbsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {

	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"smaller-neg", -2, -3, -2},
		{"smaller", 0, 1, 1},
		{"bigger", 2, 1, 2},
		{"same", 2, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt(tt.a, tt.b); got != tt.want {
				t.Errorf("MaxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxUint(t *testing.T) {
	tests := []struct {
		name string
		a, b uint
		want uint
	}{
		{"smaller", 0, 1, 1},
		{"bigger", 2, 1, 2},
		{"same", 2, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxUint(tt.a, tt.b); got != tt.want {
				t.Errorf("MaxUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt(t *testing.T) {

	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"smaller-neg", -2, -3, -3},
		{"smaller", 0, 1, 0},
		{"bigger", 2, 1, 1},
		{"same", 2, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.a, tt.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinUint(t *testing.T) {

	tests := []struct {
		name string
		a, b uint
		want uint
	}{
		{"smaller", 0, 1, 0},
		{"bigger", 2, 1, 1},
		{"same", 2, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinUint(tt.a, tt.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
