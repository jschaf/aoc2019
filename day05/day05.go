package main

import (
	"aoc2019/files"
	"aoc2019/lines"
	"fmt"
	"log"
	"math"
)

func main() {
	opsRaw := files.ReadFirstLine("day05/input.txt")

	ops1 := lines.ParseCommaSeparatedInts(opsRaw)
	code1 := NewIntCode(ops1)
	out1 := code1.RunWithInput(1)
	fmt.Printf("part 1 output: %d\n", out1)

	code2 := NewIntCode(ops1)
	out2 := code2.RunWithInput(5)
	fmt.Printf("part 2 output: %d\n", out2)
}

type intCode struct {
	mem []int
}

func NewIntCode(mem []int) *intCode {
	m := make([]int, len(mem))
	copy(m, mem)
	return &intCode{
		mem: m,
	}
}

const (
	addOp      = 1
	multOp     = 2
	inputOp    = 3
	outputOp   = 4
	jmpTrueOp  = 5
	jmpFalseOp = 6
	ltOp       = 7
	eqOp       = 8
	haltOp     = 99
)

const (
	falseV = 0
	trueV  = 1
)

const (
	positionMode  = 0
	immediateMode = 1
)

func (ic *intCode) RunWithInput(input int) (output int) {
	ip := 0
	mem := ic.mem
	output = math.MaxInt64
	for mem[ip] != haltOp {
		op := mem[ip]
		switch op % 100 {
		case haltOp:
			return

		case addOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.load(ip+3, immediateMode)
			ic.store(out, p1+p2)
			ip += 4

		case multOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.load(ip+3, immediateMode)
			ic.store(out, p1*p2)
			ip += 4

		case inputOp:
			out := ic.load(ip+1, immediateMode)
			ic.store(out, input)
			ip += 2

		case outputOp:
			out := ic.load(ip+1, immediateMode)
			output = mem[out]
			ip += 2

		case jmpTrueOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			if p1 > 0 {
				ip = p2
			} else {
				ip += 3
			}

		case jmpFalseOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			if p1 == 0 {
				ip = p2
			} else {
				ip += 3
			}

		case ltOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.load(ip+3, immediateMode)
			if p1 < p2 {
				ic.store(out, trueV)
			} else {
				ic.store(out, falseV)
			}
			ip += 4

		case eqOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.load(ip+3, immediateMode)
			if p1 == p2 {
				ic.store(out, trueV)
			} else {
				ic.store(out, falseV)
			}
			ip += 4

		default:
			log.Fatalf("unknown op code %d", op)
		}
	}
	return output
}

func mode(op, pos int) int {
	mask := op / 100
	for i := 0; i < pos; i++ {
		mask /= 10
	}
	return mask % 10
}

func (ic *intCode) load(pos int, mode int) int {
	switch mode {
	case positionMode:
		return ic.mem[ic.mem[pos]]
	case immediateMode:
		return ic.mem[pos]
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}

func (ic *intCode) store(pos int, val int) {
	ic.mem[pos] = val
}
