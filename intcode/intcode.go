package intcode

import (
	"fmt"
	"log"
	"math"
)

type Mem struct {
	mem []int
	in  chan int
	out chan int
}

func NewFromOps(mem []int) *Mem {
	m := make([]int, len(mem))
	copy(m, mem)
	return &Mem{
		mem: m,
		in:  make(chan int),
		out: make(chan int),
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

func (ic *Mem) RunWithFixedInput(inputs []int) (output int) {
	ip := 0
	curInputIdx := 0
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
			if curInputIdx >= len(inputs) {
				log.Fatalf("not enough inputs, had %d inputs but trying to use %d", len(inputs), curInputIdx)
			}
			ic.store(out, inputs[curInputIdx])
			curInputIdx += 1
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

func (ic *Mem) load(pos int, mode int) int {
	switch mode {
	case positionMode:
		return ic.mem[ic.mem[pos]]
	case immediateMode:
		return ic.mem[pos]
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}

func (ic *Mem) store(pos int, val int) {
	ic.mem[pos] = val
}

func (ic *Mem) Clone() *Mem {
	c := make([]int, len(ic.mem))
	copy(c, ic.mem)
	return NewFromOps(c)
}
