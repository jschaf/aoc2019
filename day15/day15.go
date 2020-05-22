package main

import (
	"aoc2019/files"
	"fmt"
	"log"
)

func main() {
	part1()
	part2()
}

func part1() {
	ints, err := files.ReadAllLinesAsInts("day15/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	for _, n := range ints {
		total += doPart1(n)
	}
	fmt.Printf("Part 1: %d\n", total)
}

func doPart1(n int) int {
	return 1
}

func part2() {
	ints, err := files.ReadAllLinesAsInts("day15/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	total := 0
	for _, n := range ints {
		total += doPart2(n)
	}
	fmt.Printf("Part 2: %d\n", total)
}

func doPart2(n int) int {
	return 1
}
