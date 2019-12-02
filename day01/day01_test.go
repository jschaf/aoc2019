package main

import (
	"strconv"
	"testing"
)

func Test_calcRequiredFuel(t *testing.T) {
	type args struct {
		mass int
	}
	tests := []struct {
		mass int
		want int
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.mass), func(t *testing.T) {
			if got := calcRequiredFuel(tt.mass); got != tt.want {
				t.Errorf("calcRequiredFuel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcRequiredFuelRecursive(t *testing.T) {
	tests := []struct {
		mass int
		want int
	}{
		{12, 2},
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.mass), func(t *testing.T) {
			if got := calcRequiredFuelRecursive(tt.mass); got != tt.want {
				t.Errorf("calcRequiredFuelRecursive() = %v, want %v", got, tt.want)
			}
		})
	}
}
