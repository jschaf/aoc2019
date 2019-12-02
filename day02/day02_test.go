package main

import (
	"reflect"
	"testing"
)

func newCode(ops ...[]int) []int {
	var codes []int
	for _, op := range ops {
		for _, x := range op {
			codes = append(codes, x)
		}
	}
	return codes
}

func add(x, y, pos int) []int {
	return []int{addOp, x, y, pos}
}

func mult(x, y, pos int) []int {
	return []int{multOp, x, y, pos}
}

func halt() []int {
	return []int{haltOp}
}

func Test_evalIntCode(t *testing.T) {
	tests := []struct {
		name string
		code []int
		want []int
	}{
		{"no op", newCode(halt()), newCode(halt())},
		{"1 op add", newCode(add(0, 0, 0), halt()), []int{2, 0, 0, 0, 99}},
		{"1 op mult", newCode(mult(1, 2, 3), halt()), []int{2, 1, 2, 2, 99}},
		{"1 op mult stored last", []int{2, 4, 4, 5, 99, 0}, []int{2, 4, 4, 5, 99, 9801}},
		{"modify halt", []int{1, 1, 1, 4, 99, 5, 6, 0, 99}, []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
		{
			"multi step",
			[]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50},
			[]int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evalIntCode(tt.code)
			if got := tt.code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evalIntCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
