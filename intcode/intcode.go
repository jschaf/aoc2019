package intcode

import (
	"fmt"
	"log"
	"time"
)

type Mem struct {
	mem           []int
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
		mem:    m,
		ID:     "Mem",
		Input:  make(chan int),
		Output: make(chan int),
		State:  make(chan State),
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

		case <-time.After(1 * time.Second):
			log.Fatalf("%s Timed out running with fixed input", ic.ID)
		}
	}
}

func (ic *Mem) Run() {
	defer close(ic.Output)
	defer close(ic.State)
	ip := 0
	curInputIdx := 0
	mem := ic.mem
	for {
		op := mem[ip]
		switch op % 100 {
		case haltOp:
			ic.State <- Halted
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
			ic.State <- NeedInput
			select {
			case input := <-ic.Input:
				ic.store(out, input)
				curInputIdx += 1
				ip += 2
			case <-time.After(1 * time.Second):
				log.Fatalf("%s failed to get input in 1 second", ic.ID)
			}

		case outputOp:
			out := ic.load(ip+1, immediateMode)
			ic.State <- HaveOutput
			ic.Output <- mem[out]
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
