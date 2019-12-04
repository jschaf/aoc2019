package main

import (
	"strconv"
	"testing"
)

func Test_isMatch(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{111111, true},
		{111110, false},
		{123456, false},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.n), func(t *testing.T) {
			if got := isMatch(tt.n); got != tt.want {
				t.Errorf("isMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countMatches2(t *testing.T) {
	tests := []struct {
		n    int
		want bool
	}{
		{679999, false},
		{111111, false},
		{111199, true},
		{111189, false},
		{111122, true},
		{112234, true},
		{122345, true},
		{223456, true},
		{222456, false},
		{111223, true},
		{111222, false},
		{111110, false},
		{123456, false},
		{112233, true},
		{112244, true},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.n), func(t *testing.T) {
			if got := countMatches2(tt.n, tt.n) == 1; got != tt.want {
				t.Errorf("isMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

var countMatches2Result bool

func Benchmark_countMatches2(b *testing.B) {
	var r bool
	for n := 0; n < b.N; n++ {
		r = isMatch2(172930)
		// r = countMatches2(172930, 683082)
	}
	countMatches2Result = r
}
