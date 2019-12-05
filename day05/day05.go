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
	out1 := eval(ops1, 1)
	fmt.Printf("part 1 output: %d\n", out1)

	ops2 := lines.ParseCommaSeparatedInts(opsRaw)
	out2 := eval(ops2, 5)
	fmt.Printf("part 2 output: %d\n", out2)
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

func eval(mem []int, input int) (output int) {
	ip := 0
	output = math.MaxInt64
	for mem[ip] != haltOp {
		op := mem[ip]
		switch op % 100 {
		case haltOp:
			return

		case addOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			out := load(mem, ip+3, immediateMode)
			mem[out] = p1 + p2
			ip += 4

		case multOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			out := load(mem, ip+3, immediateMode)
			mem[out] = p1 * p2
			ip += 4

		case inputOp:
			out := load(mem, ip+1, immediateMode)
			mem[out] = input
			ip += 2

		case outputOp:
			out := load(mem, ip+1, immediateMode)
			output = mem[out]
			ip += 2

		case jmpTrueOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			if p1 > 0 {
				ip = p2
			} else {
				ip += 3
			}

		case jmpFalseOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			if p1 == 0 {
				ip = p2
			} else {
				ip += 3
			}

		case ltOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			store := mem[ip+3]
			if p1 < p2 {
				mem[store] = trueV
			} else {
				mem[store] = falseV
			}
			ip += 4

		case eqOp:
			p1 := load(mem, ip+1, mode(op, 0))
			p2 := load(mem, ip+2, mode(op, 1))
			store := mem[ip+3]
			if p1 == p2 {
				mem[store] = trueV
			} else {
				mem[store] = falseV
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

func load(mem []int, ip int, mode int) int {
	switch mode {
	case positionMode:
		return mem[mem[ip]]
	case immediateMode:
		return mem[ip]
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}
