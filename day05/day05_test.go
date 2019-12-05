package main

import (
	"reflect"
	"testing"
)

const badOp = 77777

func Test_eval_basicOps(t *testing.T) {
	tests := []struct {
		name     string
		mem      []int
		afterMem []int
	}{
		{
			"add position mode",
			[]int{addOp, 0, 4, 1, haltOp},
			[]int{addOp, 100, 4, 1, haltOp},
		},
		{
			"add immediate mode",
			[]int{1100 + addOp, 3, 5, 1, haltOp},
			[]int{1100 + addOp, 8, 5, 1, haltOp},
		},
		{
			"multiply position mode",
			[]int{multOp, 0, 4, 1, haltOp},
			[]int{multOp, 198, 4, 1, haltOp},
		},
		{
			"multiply immediate mode",
			[]int{1100 + multOp, 3, 5, 1, haltOp},
			[]int{1100 + multOp, 15, 5, 1, haltOp},
		},
		{
			"jump true when true immediate mode",
			[]int{1100 + jmpTrueOp, trueV, 4, badOp, haltOp},
			[]int{1100 + jmpTrueOp, trueV, 4, badOp, haltOp},
		},
		{
			"jump true when true position mode",
			[]int{jmpTrueOp, 3, 4, trueV, 5, haltOp},
			[]int{jmpTrueOp, 3, 4, trueV, 5, haltOp},
		},
		{
			"jump true when false immediate mode",
			[]int{1100 + jmpTrueOp, falseV, 4, haltOp, badOp},
			[]int{1100 + jmpTrueOp, falseV, 4, haltOp, badOp},
		},
		{
			"jump false when true immediate mode",
			[]int{1100 + jmpFalseOp, trueV, 4, haltOp, badOp},
			[]int{1100 + jmpFalseOp, trueV, 4, haltOp, badOp},
		},
		{
			"jump false when false immediate mode",
			[]int{1100 + jmpFalseOp, falseV, 4, badOp, haltOp},
			[]int{1100 + jmpFalseOp, falseV, 4, badOp, haltOp},
		},
		{
			"lt when false immediate mode",
			[]int{1100 + ltOp, 22, 21, 5, haltOp, -1},
			[]int{1100 + ltOp, 22, 21, 5, haltOp, falseV},
		},
		{
			"lt when true immediate mode",
			[]int{1100 + ltOp, 21, 22, 5, haltOp, -1},
			[]int{1100 + ltOp, 21, 22, 5, haltOp, trueV},
		},
		{
			"eq when false immediate mode",
			[]int{1100 + eqOp, 22, 21, 5, haltOp, -1},
			[]int{1100 + eqOp, 22, 21, 5, haltOp, falseV},
		},
		{
			"eq when true immediate mode",
			[]int{1100 + eqOp, 21, 21, 5, haltOp, -1},
			[]int{1100 + eqOp, 21, 21, 5, haltOp, trueV},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewIntCode(tt.mem)
			code.RunWithInput(0)
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("eval() = %v, want %v", tt.mem, tt.afterMem)
			}
		})
	}
}

func Test_eval_inputOutput(t *testing.T) {
	tests := []struct {
		name       string
		mem        []int
		afterMem   []int
		input      int
		wantOutput int
	}{
		{
			"input",
			[]int{3, 1, 99},
			[]int{3, 88, 99},
			88,
			0,
		},
		{
			"output",
			[]int{4, 0, 99},
			[]int{4, 0, 99},
			0,
			4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewIntCode(tt.mem)
			code.RunWithInput(tt.input)
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("eval() = %v, want %v", tt.mem, tt.afterMem)
			}
		})
	}
}
