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
	for i := len(modes) - 1; i >= 0; i-- {
		m := modes[i]
		v *= 10
		v += m
	}
	return (v * 100) + op
}

func rel1(op int) int { return m(op, relativeMode) }
func imm1(op int) int { return m(op, immediateMode) }
func imm2(op int) int { return m(op, immediateMode, immediateMode) }

func Test_mode(t *testing.T) {
	tests := []struct {
		name      string
		val, want int
	}{
		{"rel-imm-add", m(AddOp, rel, imm), 12e2 + AddOp},
		{"rel-imm-jmp", m(JmpTrueOp, imm, rel), 21e2 + JmpTrueOp},
		{"pos-imm-add", m(AddOp, pos, imm), 10e2 + AddOp},
		{"imm-pos-add", m(AddOp, imm, pos), 1e2 + AddOp},
		{"pos-imm-eq", m(EqOp, pos, imm), 10e2 + EqOp},
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
			[]int{m(AdjRelBaseOp, imm), 2, m(MultOp, rel, rel), 5, 6, 9, HaltOp, 11, 13, -1},
			[]int{m(AdjRelBaseOp, imm), 2, m(MultOp, rel, rel), 5, 6, 9, HaltOp, 11, 13, 143},
		},
		{
			"with relative base and many outputs",
			[]int{109, 1, 204, -1, HaltOp},
			[]int{109, 1, 204, -1, HaltOp},
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
			[]int{AddOp, 0, 4, 1, HaltOp},
			[]int{AddOp, 100, 4, 1, HaltOp},
		},
		{
			"add immediate mode",
			[]int{imm2(AddOp), 3, 5, 1, HaltOp},
			[]int{imm2(AddOp), 8, 5, 1, HaltOp},
		},
		{
			"multiply position mode",
			[]int{MultOp, 0, 4, 1, HaltOp},
			[]int{MultOp, 198, 4, 1, HaltOp},
		},
		{
			"multiply immediate mode",
			[]int{imm2(MultOp), 3, 5, 1, HaltOp},
			[]int{imm2(MultOp), 15, 5, 1, HaltOp},
		},
		{
			"jump true when true immediate mode",
			[]int{imm2(JmpTrueOp), TrueV, 4, badOp, HaltOp},
			[]int{imm2(JmpTrueOp), TrueV, 4, badOp, HaltOp},
		},
		{
			"jump true when true position mode",
			[]int{JmpTrueOp, 3, 4, TrueV, 5, HaltOp},
			[]int{JmpTrueOp, 3, 4, TrueV, 5, HaltOp},
		},
		{
			"jump true when false immediate mode",
			[]int{imm2(JmpTrueOp), FalseV, 4, HaltOp, badOp},
			[]int{imm2(JmpTrueOp), FalseV, 4, HaltOp, badOp},
		},
		{
			"jump false when true immediate mode",
			[]int{imm2(JmpFalseOp), TrueV, 4, HaltOp, badOp},
			[]int{imm2(JmpFalseOp), TrueV, 4, HaltOp, badOp},
		},
		{
			"jump false when false immediate mode",
			[]int{imm2(JmpFalseOp), FalseV, 4, badOp, HaltOp},
			[]int{imm2(JmpFalseOp), FalseV, 4, badOp, HaltOp},
		},
		{
			"lt when false immediate mode",
			[]int{imm2(LtOp), 22, 21, 5, HaltOp, -1},
			[]int{imm2(LtOp), 22, 21, 5, HaltOp, FalseV},
		},
		{
			"lt when true immediate mode",
			[]int{imm2(LtOp), 21, 22, 5, HaltOp, -1},
			[]int{imm2(LtOp), 21, 22, 5, HaltOp, TrueV},
		},
		{
			"eq when false immediate mode",
			[]int{imm2(EqOp), 22, 21, 5, HaltOp, -1},
			[]int{imm2(EqOp), 22, 21, 5, HaltOp, FalseV},
		},
		{
			"eq when true immediate mode",
			[]int{imm2(EqOp), 21, 21, 5, HaltOp, -1},
			[]int{imm2(EqOp), 21, 21, 5, HaltOp, TrueV},
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
			[]int{OutputOp, 6, HaltOp},
			[]int{OutputOp, 6, HaltOp},
			[]int{},
			[]int{0},
		},
		{
			"resize for writes",
			[]int{imm2(AddOp), 3, 5, 6, HaltOp},
			[]int{imm2(AddOp), 3, 5, 6, HaltOp, 0, 8},
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

func TestMem_RunWithFixedInput_quine(t *testing.T) {
	tests := []struct {
		name       string
		mem        []int
		afterMem   []int
		inputs     []int
		wantOutput []int
	}{
		{
			"output from quine",
			[]int{
				imm1(AdjRelBaseOp), 1,
				rel1(OutputOp), -1,
				m(AddOp, pos, imm), 100, 1, 100,
				m(EqOp, pos, imm), 100, 16, 101,
				m(JmpFalseOp, pos, imm), 101, 0,
				HaltOp},
			[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
			[]int{},
			[]int{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := NewFromOps(tt.mem)
			outputs := code.RunWithFixedInput(tt.inputs)
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
			[]int{InputOp, 1, HaltOp},
			[]int{InputOp, 88, HaltOp},
			[]int{88},
			[]int{},
		},
		{
			"multiple inputs",
			[]int{InputOp, 1, InputOp, 3, HaltOp},
			[]int{InputOp, 88, InputOp, 99, HaltOp},
			[]int{88, 99},
			[]int{},
		},
		{
			"output",
			[]int{OutputOp, 3, HaltOp, 667},
			[]int{OutputOp, 3, HaltOp, 667},
			[]int{},
			[]int{667},
		},
		{
			"multiple outputs",
			[]int{OutputOp, 5, OutputOp, 6, HaltOp, 667, 668},
			[]int{OutputOp, 5, OutputOp, 6, HaltOp, 667, 668},
			[]int{},
			[]int{667, 668},
		},
		{
			"output from negative relative base",
			[]int{109, 1, 204, -1, HaltOp},
			[]int{109, 1, 204, -1, HaltOp},
			[]int{},
			[]int{109},
		},
		{
			"output large number",
			[]int{104, 1125899906842624, 99},
			[]int{104, 1125899906842624, 99},
			[]int{},
			[]int{1125899906842624},
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
