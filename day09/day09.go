package main

import (
	"aoc2019/files"
	"aoc2019/intcode"
	"fmt"
	"log"
)

func main() {
	origCode := parseInputToIntCode("day09/input.txt")
	p1 := part1(origCode)
	fmt.Printf("Part 1 code: %d\n", p1)

	p2 := part2(origCode)
	fmt.Printf("Part 2 coords: %d\n", p2)
}

func parseInputToIntCode(file string) *intcode.Mem {
	line := files.ReadFirstLine(file)
	return intcode.ParseLine(line)
}

func part1(mem *intcode.Mem) int {
	output := mem.Clone().RunWithFixedInput([]int{1})
	if len(output) > 1 {
		log.Fatalf("Had malfuctioning opcodes: %v", output)
	}
	return output[0]
}

func part2(mem *intcode.Mem) int {
	output := mem.Clone().RunWithFixedInput([]int{2})
	if len(output) > 1 {
		log.Fatalf("Had malfuctioning opcodes: %v", output)
	}
	return output[0]
}
