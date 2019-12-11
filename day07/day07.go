package main

import (
	"aoc2019/combinations"
	"aoc2019/files"
	"aoc2019/intcode"
	"aoc2019/maths"
	"fmt"
	"log"
	"math"
	"sync"
)

func main() {
	origCode := parseInputToIntCode("day07/input.txt")
	p1 := part1(origCode)
	fmt.Printf("Part 1 best signal: %d\n", p1)

	p2 := part2(origCode)
	fmt.Printf("Part 2 best signal: %d\n", p2)
}

func parseInputToIntCode(file string) *intcode.Mem {
	line := files.ReadFirstLine(file)
	return intcode.ParseLine(line)
}

func part1(origCode *intcode.Mem) int {
	phases := []int{0, 1, 2, 3, 4}
	bestSignal := math.MinInt64
	for p := combinations.NewPermuter(phases); p.HasNext(); p.Next() {
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
	bestSignal := math.MinInt64
	for p := combinations.NewPermuter(possiblePhases); p.HasNext(); p.Next() {
		phases := p.Get()
		o := getOutputForPhases(origCode, phases)
		if o > bestSignal {
			bestSignal = o
		}
	}
	return bestSignal
}

func getOutputForPhases(origCode *intcode.Mem, phases []int) int {
	numAmps := len(phases)
	amps := make([]*intcode.Mem, numAmps)
	for i := range amps {
		amps[i] = origCode.Clone()
		amps[i].ID = fmt.Sprintf("amp_%d", i)
		go amps[i].Run()
	}

	outputs := make([]chan int, numAmps)
	for i, phase := range phases {
		outputs[i] = make(chan int, 10)
		outputs[i] <- phase
	}
	outputs[0] <- 0

	wg := sync.WaitGroup{}
	wg.Add(numAmps)

	for i, amp := range amps {
		go func(i int, amp *intcode.Mem) {
		runLoop:
			for {
				select {
				case s := <-amp.State:
					switch s {
					case intcode.NeedInput:
						prevI := maths.ModWithSameSign(i-1, numAmps)
						amp.Input <- <-outputs[prevI]

					case intcode.HaveOutput:
						o := <-amp.Output
						outputs[i] <- o

					case intcode.Halted:
						wg.Done()
						close(outputs[i])
						break runLoop
					}
				}
			}
		}(i, amp)
	}

	wg.Wait()

	lastOuts := make([]int, 0)
	for x := range outputs[0] {
		lastOuts = append(lastOuts, x)
	}
	if len(lastOuts) != 1 {
		log.Fatalf("expected exactly 1 element in last amp output but had %v", lastOuts)
	}
	return lastOuts[0]
}
