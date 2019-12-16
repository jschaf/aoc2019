package main

import (
	"aoc2019/files"
	ic "aoc2019/intcode"
	"testing"
)

func newIntCode(codes ...[]int) *ic.Mem {
	mem := make([]int, 0)
	for _, code := range codes {
		mem = append(mem, code...)
	}
	return ic.NewFromOps(mem)
}

func Test_part1(t *testing.T) {
	line := files.ReadFirstLine("input.txt")
	paintWhiteTurnLeft := []int{
		ic.InputOp, 1,
		ic.Mode(ic.OutputOp, ic.ImmediateMode), white,
		ic.Mode(ic.OutputOp, ic.ImmediateMode), turnLeft,
	}
	tests := []struct {
		name string
		code *ic.Mem
		want int
	}{
		{"paint 1 white", newIntCode(paintWhiteTurnLeft, []int{ic.HaltOp}), 1},
		{"paint 5 white - 1 dupe", newIntCode(
			paintWhiteTurnLeft,
			paintWhiteTurnLeft,
			paintWhiteTurnLeft,
			paintWhiteTurnLeft,
			paintWhiteTurnLeft,
			[]int{ic.HaltOp},
		), 4},
		{"input", ic.ParseLine(line), 2511},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part1(tt.code); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}
