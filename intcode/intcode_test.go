package intcode

import (
	"reflect"
	"testing"
)

const badOp = 77777

const (
	pos = positionMode
	imm = immediateMode
	rel = relativeMode
)

func m(op int, modes ...int) int {
	v := 0
	for _, m := range modes {
		v *= 10
		v += m
	}
	return (v * 100) + op
}
func imm2(op int) int { return m(op, immediateMode, immediateMode) }

func Test_mode(t *testing.T) {
	tests := []struct {
		name      string
		val, want int
	}{
		{"rel-imm-add", m(addOp, rel, imm), 2100 + addOp},
		{"rel-imm-jmp", m(jmpTrueOp, rel, imm), 2100 + jmpTrueOp},
		{"pos-imm-add", m(addOp, pos, imm), 100 + addOp},
		{"imm-pos-add", m(addOp, imm, pos), 1000 + addOp},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val != tt.want {
				t.Errorf("mode() = %v, want %v", tt.val, tt.want)
			}
		})
	}
}

func TestMem_RunWithFixedInputs_relBase(t *testing.T) {
	tests := []struct {
		name     string
		mem      []int
		afterMem []int
	}{
		{
			"uses relative base",
			[]int{m(adjRelBaseOp, imm), 2, m(multOp, rel, rel), 5, 6, 9, haltOp, 11, 13, -1},
			[]int{m(adjRelBaseOp, imm), 2, m(multOp, rel, rel), 5, 6, 9, haltOp, 11, 13, 143},
		},
		{
			"with relative base and many outputs",
			[]int{109, 1, 204, -1, haltOp},
			[]int{109, 1, 204, -1, haltOp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewFromOps(tt.mem)
			code.RunWithFixedInput([]int{})
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("RunWithFixedInput() = %v, want %v", code.mem, tt.afterMem)
			}
		})
	}
}

func TestMem_RunWithFixedInputs_basicOps(t *testing.T) {
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
			[]int{imm2(addOp), 3, 5, 1, haltOp},
			[]int{imm2(addOp), 8, 5, 1, haltOp},
		},
		{
			"multiply position mode",
			[]int{multOp, 0, 4, 1, haltOp},
			[]int{multOp, 198, 4, 1, haltOp},
		},
		{
			"multiply immediate mode",
			[]int{imm2(multOp), 3, 5, 1, haltOp},
			[]int{imm2(multOp), 15, 5, 1, haltOp},
		},
		{
			"jump true when true immediate mode",
			[]int{imm2(jmpTrueOp), trueV, 4, badOp, haltOp},
			[]int{imm2(jmpTrueOp), trueV, 4, badOp, haltOp},
		},
		{
			"jump true when true position mode",
			[]int{jmpTrueOp, 3, 4, trueV, 5, haltOp},
			[]int{jmpTrueOp, 3, 4, trueV, 5, haltOp},
		},
		{
			"jump true when false immediate mode",
			[]int{imm2(jmpTrueOp), falseV, 4, haltOp, badOp},
			[]int{imm2(jmpTrueOp), falseV, 4, haltOp, badOp},
		},
		{
			"jump false when true immediate mode",
			[]int{imm2(jmpFalseOp), trueV, 4, haltOp, badOp},
			[]int{imm2(jmpFalseOp), trueV, 4, haltOp, badOp},
		},
		{
			"jump false when false immediate mode",
			[]int{imm2(jmpFalseOp), falseV, 4, badOp, haltOp},
			[]int{imm2(jmpFalseOp), falseV, 4, badOp, haltOp},
		},
		{
			"lt when false immediate mode",
			[]int{imm2(ltOp), 22, 21, 5, haltOp, -1},
			[]int{imm2(ltOp), 22, 21, 5, haltOp, falseV},
		},
		{
			"lt when true immediate mode",
			[]int{imm2(ltOp), 21, 22, 5, haltOp, -1},
			[]int{imm2(ltOp), 21, 22, 5, haltOp, trueV},
		},
		{
			"eq when false immediate mode",
			[]int{imm2(eqOp), 22, 21, 5, haltOp, -1},
			[]int{imm2(eqOp), 22, 21, 5, haltOp, falseV},
		},
		{
			"eq when true immediate mode",
			[]int{imm2(eqOp), 21, 21, 5, haltOp, -1},
			[]int{imm2(eqOp), 21, 21, 5, haltOp, trueV},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewFromOps(tt.mem)
			code.RunWithFixedInput([]int{})
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("RunWithFixedInput() = %v, want %v", tt.mem, tt.afterMem)
			}
		})
	}
}

func TestMem_RunWithFixedInput_resize(t *testing.T) {
	tests := []struct {
		name       string
		mem        []int
		afterMem   []int
		inputs     []int
		wantOutput []int
	}{
		{
			"no resize for reads",
			[]int{outputOp, 6, haltOp},
			[]int{outputOp, 6, haltOp},
			[]int{},
			[]int{0},
		},
		{
			"resize for writes",
			[]int{imm2(addOp), 3, 5, 6, haltOp},
			[]int{imm2(addOp), 3, 5, 6, haltOp, 0, 8},
			[]int{},
			[]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewFromOps(tt.mem)
			outputs := code.RunWithFixedInput(tt.inputs)
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("RunWithFixedInput() = %v, want %v", code.mem, tt.afterMem)
			}
			if !reflect.DeepEqual(outputs, tt.wantOutput) {
				t.Errorf("outputs = %v, want %v", outputs, tt.wantOutput)
			}
		})
	}
}

func TestMem_RunWithFixedInput_inputOutput(t *testing.T) {
	tests := []struct {
		name       string
		mem        []int
		afterMem   []int
		inputs     []int
		wantOutput []int
	}{
		{
			"input",
			[]int{inputOp, 1, haltOp},
			[]int{inputOp, 88, haltOp},
			[]int{88},
			[]int{},
		},
		{
			"multiple inputs",
			[]int{inputOp, 1, inputOp, 3, haltOp},
			[]int{inputOp, 88, inputOp, 99, haltOp},
			[]int{88, 99},
			[]int{},
		},
		{
			"output",
			[]int{outputOp, 3, haltOp, 667},
			[]int{outputOp, 3, haltOp, 667},
			[]int{},
			[]int{667},
		},
		{
			"multiple outputs",
			[]int{outputOp, 5, outputOp, 6, haltOp, 667, 668},
			[]int{outputOp, 5, outputOp, 6, haltOp, 667, 668},
			[]int{},
			[]int{667, 668},
		},
		{
			"output from negative relative base",
			[]int{109, 1, 204, -1, haltOp},
			[]int{109, 1, 204, -1, haltOp},
			[]int{},
			[]int{109},
		},
		// {
		// 	"output from negative relative base",
		// 	[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		// 	[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		// 	[]int{},
		// 	[]int{0},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewFromOps(tt.mem)
			outputs := code.RunWithFixedInput(tt.inputs)
			if !reflect.DeepEqual(code.mem, tt.afterMem) {
				t.Errorf("RunWithFixedInput() = %v, want %v", code.mem, tt.afterMem)
			}
			if !reflect.DeepEqual(outputs, tt.wantOutput) {
				t.Errorf("outputs = %v, want %v", outputs, tt.wantOutput)
			}
		})
	}
}
