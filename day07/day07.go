package main

import (
	"aoc2019/files"
	"aoc2019/intcode"
	"aoc2019/lines"
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

func main() {
	opsRaw := files.ReadFirstLine("day07/input.txt")
	ops1 := lines.ParseCommaSeparatedInts(opsRaw)
	origCode := intcode.NewFromOps(ops1)

	// p1 := part1(origCode)
	// fmt.Printf("Part 1 best signal: %d", p1)

	p2 := part2(origCode)
	fmt.Printf("Part 2 best signal: %d", p2)

}

func part1(origCode *intcode.Mem) int {
	phases := []int{0, 1, 2, 3, 4}
	bestSignal := math.MinInt64
	for p := NewPermuter(phases); p.HasNext(); p.Next() {
		ampPhases := p.Get()
		prevOutput := 0
		for _, phase := range ampPhases {
			code := origCode.Clone()
			outputs := code.RunWithFixedInput([]int{phase, prevOutput})
			prevOutput = outputs[len(outputs)-1]
		}
		if prevOutput > bestSignal {
			bestSignal = prevOutput
		}
	}
	return bestSignal
}

func part2(origCode *intcode.Mem) int {
	possiblePhases := []int{5, 6, 7, 8, 9}
	numAmps := len(possiblePhases)
	bestSignal := math.MinInt64
	numPermutes := 0
	for p := NewPermuter(possiblePhases); p.HasNext(); p.Next() {
		amps := make([]*intcode.Mem, numAmps)
		for i := range amps {
			amps[i] = origCode.Clone()
			amps[i].ID = fmt.Sprintf("amp_%d", i)
			go amps[i].Run()
		}

		// Initial inputs
		phases := p.Get()
		fmt.Printf("\n\n\nPart 2: permutation %d: %v\n", numPermutes, phases)
		pendingInputs := make([][]int, len(amps))
		for i, phase := range phases {
			pendingInputs[i] = append(pendingInputs[i], phase)
		}
		pendingInputs[0] = append(pendingInputs[0], 0)

		// If last run, get Amp E's output instead of forwarding to pending input.
		// lastRun := false
		numRemaining := int64(numAmps)
		isQuit := make([]bool, numAmps)

		// Connect amps 1 step at a time
		step := 0
		for numRemaining > 0 {
			time.Sleep(1 * time.Microsecond)
			// fmt.Printf("  Step %d, num remaining: %d\n", step, numRemaining)
			step++
			for i := 0; i < numAmps; i++ {
				amp := amps[i]
				if isQuit[i] {
					// fmt.Printf("    %s - skipping because quit\n", amp.ID)
					continue
				}

				select {

				case o, ok := <-amp.Output:
					if !ok {
						// fmt.Printf("    %s - skipping output because closed\n", amp.ID)
						break
					}
					j := (i + 1) % len(amps)
					pendingInputs[j] = append(pendingInputs[j], o)
					// fmt.Printf("    %s - had output %d\n", amp.ID, o)

				case <-amp.Quit:
					// fmt.Printf("    %s - quit\n", amp.ID)
					isQuit[i] = true
					atomic.AddInt64(&numRemaining, -1)

				default:
					if len(pendingInputs[i]) > 0 {
						x := pendingInputs[i][0]
						// fmt.Printf("    %s - default have input %d\n", amp.ID, x)
						pendingInputs[i] = pendingInputs[i][1:]
						amp.Input <- x
					} else {
						// fmt.Printf("    %s - default nothing to do\n", amp.ID)
					}
				}
			}
		}

		final := pendingInputs[0][0]
		fmt.Printf("    %s - final %d, best: %d\n", amps[4].ID, final, bestSignal)
		if final > bestSignal {
			bestSignal = final
		}
	}
	return bestSignal
}

type perm struct {
	orig  []int
	state []int
}

func NewPermuter(orig []int) *perm {
	state := make([]int, len(orig))
	return &perm{orig, state}
}

func (p *perm) HasNext() bool {
	return p.state[0] < len(p.state)
}

func (p *perm) Next() {
	for i := len(p.state) - 1; i >= 0; i-- {
		if i == 0 || p.state[i] < len(p.state)-i-1 {
			p.state[i]++
			return
		}
		p.state[i] = 0
	}
}

func (p *perm) Get() []int {
	result := make([]int, len(p.orig))
	copy(result, p.orig)
	for i, v := range p.state {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}
