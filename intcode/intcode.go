package intcode

import (
	"aoc2019/lines"
	"fmt"
	"log"
	"time"
)

type Mem struct {
	mem           []int
	relBase       int
	ID            string
	Input, Output chan int
	State         chan State
}

type State int

const (
	NeedInput = iota
	HaveOutput
	Halted
)

func NewFromOps(mem []int) *Mem {
	m := make([]int, len(mem))
	copy(m, mem)
	return &Mem{
		mem:     m,
		relBase: 0,
		ID:      "Mem",
		Input:   make(chan int),
		Output:  make(chan int),
		State:   make(chan State),
	}
}

func ParseLine(line string) *Mem {
	return NewFromOps(lines.ParseCommaSeparatedInts(line))
}

const (
	addOp        = 1
	multOp       = 2
	inputOp      = 3
	outputOp     = 4
	jmpTrueOp    = 5
	jmpFalseOp   = 6
	ltOp         = 7
	eqOp         = 8
	adjRelBaseOp = 9
	haltOp       = 99
)

const (
	falseV = 0
	trueV  = 1
)

const (
	positionMode  = 0
	immediateMode = 1
	relativeMode  = 2
)

func (ic *Mem) RunWithFixedInput(inputs []int) []int {
	go func() {
		for _, x := range inputs {
			ic.Input <- x
		}
	}()

	go ic.Run()

	outputs := make([]int, 0)
	for {
		select {

		case v := <-ic.Output:
			outputs = append(outputs, v)

		case s := <-ic.State:
			if s == Halted {
				return outputs
			}

		case <-time.After(10000 * time.Second):
			log.Fatalf("%s Timed out running with fixed input", ic.ID)
		}
	}
}

func (ic *Mem) Run() {
	defer close(ic.Output)
	defer close(ic.State)
	ip := 0
	curInputIdx := 0
	for {
		op := ic.mem[ip]
		switch op % 100 {
		case haltOp:
			ic.State <- Halted
			return

		case addOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			ic.store(out, p1+p2)
			ip += 4

		case multOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			ic.store(out, p1*p2)
			ip += 4

		case inputOp:
			out := ic.loadWriteAddr(ip+1, mode(op, 0))
			ic.State <- NeedInput
			select {
			case input := <-ic.Input:
				ic.store(out, input)
				curInputIdx += 1
				ip += 2
			case <-time.After(10000 * time.Second):
				log.Fatalf("%s failed to get input in 1 second", ic.ID)
			}

		case outputOp:
			out := ic.load(ip+1, mode(op, 0))
			ic.State <- HaveOutput
			ic.Output <- out
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
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			if p1 < p2 {
				ic.store(out, trueV)
			} else {
				ic.store(out, falseV)
			}
			ip += 4

		case eqOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			if p1 == p2 {
				ic.store(out, trueV)
			} else {
				ic.store(out, falseV)
			}
			ip += 4

		case adjRelBaseOp:
			p1 := ic.load(ip+1, mode(op, 0))
			ic.relBase += p1
			ip += 2

		default:
			log.Fatalf("unknown op code %d", op)
		}
	}
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
		return ic.get(ic.get(pos))
	case immediateMode:
		return ic.get(pos)
	case relativeMode:
		return ic.get(ic.relBase + ic.get(pos))
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}
func (ic *Mem) loadWriteAddr(pos int, mode int) int {
	switch mode {
	case positionMode:
		return ic.get(pos)
	case immediateMode:
		panic("write addr should not be immediate")
	case relativeMode:
		return ic.relBase + ic.get(pos)
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}

func (ic *Mem) get(i int) int {
	if i >= len(ic.mem) {
		return 0
	}
	return ic.mem[i]
}

func (ic *Mem) store(pos int, val int) {
	ic.maybeGrow(pos)
	ic.mem[pos] = val
}

func (ic *Mem) Clone() *Mem {
	c := make([]int, len(ic.mem))
	copy(c, ic.mem)
	return NewFromOps(c)
}

func (ic *Mem) maybeGrow(i int) {
	if i < len(ic.mem) {
		return
	}
	if i < cap(ic.mem) {
		ic.mem = ic.mem[:i+1]
		return
	}
	m := make([]int, i+1, (i+1)*2)
	copy(m, ic.mem)
	ic.mem = m
}
