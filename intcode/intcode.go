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
	state         State
	seed          []int
	seedIdx       int
}

type State int

const (
	NeedInput = iota
	HaveOutput
	Halted
)

const (
	isUnstarted = iota
	isStarted
	isHalted
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
	AddOp        = 1
	MultOp       = 2
	InputOp      = 3
	OutputOp     = 4
	JmpTrueOp    = 5
	JmpFalseOp   = 6
	LtOp         = 7
	EqOp         = 8
	AdjRelBaseOp = 9
	HaltOp       = 99
)

const (
	FalseV = 0
	TrueV  = 1
)

const (
	PositionMode  = 0
	ImmediateMode = 1
	RelativeMode  = 2
)

func (ic *Mem) RunWithFixedInput(inputs []int) []int {
	ic.SeedInput(inputs)
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

func (ic *Mem) SeedInput(seed []int) {
	if ic.state != isUnstarted {
		panic("intcode is already started or already halted")
	}
	ic.seed = seed
	ic.seedIdx = 0
}

func (ic *Mem) Set(addr, val int) {
	if ic.state != isUnstarted {
		panic("intcode is already started or already halted")
	}
	ic.mem[addr] = val
}

func (ic *Mem) Run() {
	ic.state = isStarted
	defer close(ic.Output)
	defer close(ic.State)
	ip := 0
	for {
		op := ic.mem[ip]
		switch op % 100 {
		case HaltOp:
			ic.state = isHalted
			ic.State <- Halted
			return

		case AddOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			ic.store(out, p1+p2)
			ip += 4

		case MultOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			ic.store(out, p1*p2)
			ip += 4

		case InputOp:
			out := ic.loadWriteAddr(ip+1, mode(op, 0))
			if ic.seedIdx < len(ic.seed) {
				ic.store(out, ic.seed[ic.seedIdx])
				ic.seedIdx++
			} else {
				ic.State <- NeedInput
				select {
				case input := <-ic.Input:
					ic.store(out, input)
				case <-time.After(1 * time.Second):
					log.Fatalf("%s failed to get input in 1 second", ic.ID)
				}
			}
			ip += 2

		case OutputOp:
			out := ic.load(ip+1, mode(op, 0))
			ic.State <- HaveOutput
			ic.Output <- out
			ip += 2

		case JmpTrueOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			if p1 > 0 {
				ip = p2
			} else {
				ip += 3
			}

		case JmpFalseOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			if p1 == 0 {
				ip = p2
			} else {
				ip += 3
			}

		case LtOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			if p1 < p2 {
				ic.store(out, TrueV)
			} else {
				ic.store(out, FalseV)
			}
			ip += 4

		case EqOp:
			p1 := ic.load(ip+1, mode(op, 0))
			p2 := ic.load(ip+2, mode(op, 1))
			out := ic.loadWriteAddr(ip+3, mode(op, 2))
			if p1 == p2 {
				ic.store(out, TrueV)
			} else {
				ic.store(out, FalseV)
			}
			ip += 4

		case AdjRelBaseOp:
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
	case PositionMode:
		return ic.get(ic.get(pos))
	case ImmediateMode:
		return ic.get(pos)
	case RelativeMode:
		return ic.get(ic.relBase + ic.get(pos))
	default:
		panic(fmt.Sprintf("unknown mode %d", mode))
	}
}
func (ic *Mem) loadWriteAddr(pos int, mode int) int {
	switch mode {
	case PositionMode:
		return ic.get(pos)
	case ImmediateMode:
		panic("write addr should not be immediate")
	case RelativeMode:
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

func Mode(op int, modes ...int) int {
	v := 0
	for i := len(modes) - 1; i >= 0; i-- {
		m := modes[i]
		v *= 10
		v += m
	}
	return (v * 100) + op
}
