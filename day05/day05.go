package main

import (
	"aoc2019/files"
	"aoc2019/intcode"
	"aoc2019/lines"
	"fmt"
)

func main() {
	opsRaw := files.ReadFirstLine("day05/input.txt")

	ops1 := lines.ParseCommaSeparatedInts(opsRaw)
	code1 := intcode.NewFromOps(ops1)
	out1 := code1.RunWithFixedInput([]int{1})
	fmt.Printf("part 1 output: %v\n", out1[len(out1)-1])

	code2 := intcode.NewFromOps(ops1)
	out2 := code2.RunWithFixedInput([]int{5})
	fmt.Printf("part 2 output: %v\n", out2[len(out2)-1])
}
