package main

import (
	"aoc2019/files"
	"aoc2019/lines"
	"errors"
	"fmt"
	"log"
)

func main() {
	part1()
	part2()
}

// Opcode 1 adds together numbers read from two positions and stores the result
// in a third position.
//
// Opcode 2 works exactly like opcode 1, except it multiplies the two inputs
// instead of adding them.

func readIntCode() []int {
	line := files.ReadFirstLine("day02/input.txt")
	return lines.ParseCommaSeparatedInts(line)
}

func part1() {
	code := readIntCode()
	// replace position 1 with the value 12 and replace position 2 with the value 2.
	code[1] = 12
	code[2] = 2
	// What value is left at position 0 after the program halts?
	evalIntCode(code)
	fmt.Printf("Part 1: %d\n", code[0])
}

const (
	addOp  = 1
	multOp = 2
	haltOp = 99
)

func evalIntCode(mem []int) {
	ip := 0
	for mem[ip] != haltOp {
		op := mem[ip]
		switch op {
		case haltOp:
			return
		case addOp:
			arg1 := mem[mem[ip+1]]
			arg2 := mem[mem[ip+2]]
			mem[mem[ip+3]] = arg1 + arg2
		case multOp:
			arg1 := mem[mem[ip+1]]
			arg2 := mem[mem[ip+2]]
			mem[mem[ip+3]] = arg1 * arg2
		}
		ip += 4
	}
}

func part2() {
	// What pair of inputs produces 19690720
	goal := 19690720
	origMem := readIntCode()
	noun, verb, err := findInputsForGoal(goal, origMem)
	if err != nil {
		log.Fatal("goal not found")
	}
	fmt.Printf("Found goal: nount=%d, verb=%d, 100*noun+verb=%d", noun, verb, 100*noun+verb)
}

func findInputsForGoal(goal int, origMem []int) (int, int, error) {
	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			mem := make([]int, len(origMem))
			copy(mem, origMem)
			mem[1] = x
			mem[2] = y
			evalIntCode(mem)
			if mem[0] == goal {
				return x, y, nil
			}
		}
	}
	return 0, 0, errors.New("not found")
}
